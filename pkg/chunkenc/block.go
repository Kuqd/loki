package chunkenc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/grafana/loki/pkg/iter"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/grafana/loki/pkg/logql/log"
	"github.com/grafana/loki/pkg/logql/stats"
	"github.com/pkg/errors"
)

var blockPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(nil)
	},
}

func getBlockBuffer() *bytes.Buffer {
	b := blockPool.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

func putBlockBuffer(b *block) {
	if b.reuse {
		blockPool.Put(b.buff)
	}
}

type block struct {
	reuse    bool
	isClose  bool
	released bool
	mtx      sync.Mutex
	inUsed   int

	buff       *bytes.Buffer
	numEntries int

	mint, maxt int64

	offset           int // The offset of the block in the chunk.
	uncompressedSize int // Total uncompressed size in bytes when the chunk is cut.
}

// This block holds the un-compressed entries. Once it has enough data, this is
// emptied into a block with only compressed entries.
type headBlock struct {
	// This is the list of raw entries.
	entries []entry
	size    int // size of uncompressed bytes.

	mint, maxt int64

	pool WriterPool
}

func (hb *headBlock) isEmpty() bool {
	return len(hb.entries) == 0
}

func (hb *headBlock) clear() {
	if hb.entries != nil {
		hb.entries = hb.entries[:0]
	}
	hb.size = 0
	hb.mint = 0
	hb.maxt = 0
}

func (hb *headBlock) append(ts int64, line string) error {
	if !hb.isEmpty() && hb.maxt > ts {
		return ErrOutOfOrder
	}

	hb.entries = append(hb.entries, entry{ts, line})
	if hb.mint == 0 || hb.mint > ts {
		hb.mint = ts
	}
	hb.maxt = ts
	hb.size += len(line)

	return nil
}

func (hb *headBlock) WriteTo(w io.Writer) (int64, error) {
	inBuf := serializeBytesBufferPool.Get().(*bytes.Buffer)
	defer func() {
		inBuf.Reset()
		serializeBytesBufferPool.Put(inBuf)
	}()

	encBuf := make([]byte, binary.MaxVarintLen64)
	compressedWriter := hb.pool.GetWriter(w)
	defer hb.pool.PutWriter(compressedWriter)
	for _, logEntry := range hb.entries {
		n := binary.PutVarint(encBuf, logEntry.t)
		inBuf.Write(encBuf[:n])

		n = binary.PutUvarint(encBuf, uint64(len(logEntry.s)))
		inBuf.Write(encBuf[:n])

		inBuf.WriteString(logEntry.s)
	}
	var n int
	var err error
	if n, err = compressedWriter.Write(inBuf.Bytes()); err != nil {
		return int64(n), errors.Wrap(err, "appending entry")
	}
	if err := compressedWriter.Close(); err != nil {
		return int64(n), errors.Wrap(err, "flushing pending compress buffer")
	}

	return int64(n), nil
}

// encBlock is an internal wrapper for a block, mainly to avoid binding an encoding in a block itself.
// This may seem roundabout, but the encoding is already a field on the parent MemChunk type. encBlock
// then allows us to bind a decoding context to a block when requested, but otherwise helps reduce the
// chances of chunk<>block encoding drift in the codebase as the latter is parameterized by the former.
type encBlock struct {
	enc Encoding
	*block
}

func (b encBlock) Iterator(ctx context.Context, pipeline log.StreamPipeline) iter.EntryIterator {
	if len(b.buff.Bytes()) == 0 {
		return iter.NoopIterator
	}
	if !b.block.use() {
		return iter.NoopIterator
	}
	return newEntryIterator(ctx, getReaderPool(b.enc), b.block, pipeline)
}

func (b encBlock) SampleIterator(ctx context.Context, extractor log.StreamSampleExtractor) iter.SampleIterator {
	if len(b.buff.Bytes()) == 0 {
		return iter.NoopIterator
	}
	if !b.block.use() {
		return iter.NoopIterator
	}
	return newSampleIterator(ctx, getReaderPool(b.enc), b.block, extractor)
}

func (b *block) Offset() int {
	return b.offset
}

//todo close on block should release only if close and last.
// close on chunk should release only if no block leaked.
func (b *block) close() {
	if !b.reuse {
		return
	}
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if b.isClose {
		return
	}
	b.isClose = true
	if b.inUsed == 0 && !b.released {
		b.released = true
		putBlockBuffer(b)
		return
	}

}

func (b *block) use() bool {
	if !b.reuse {
		return true
	}
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if b.isClose || b.released {
		return false
	}
	b.inUsed++
	return true
}

func (b *block) release() {
	if !b.reuse {
		return
	}
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.inUsed--
	if b.inUsed == 0 && b.isClose && !b.released {
		b.released = true
		putBlockBuffer(b)
	}
}

func (b *block) Entries() int {
	return b.numEntries
}
func (b *block) MinTime() int64 {
	return b.mint
}
func (b *block) MaxTime() int64 {
	return b.maxt
}

func (hb *headBlock) iterator(ctx context.Context, direction logproto.Direction, mint, maxt int64, pipeline log.StreamPipeline) iter.EntryIterator {
	if hb.isEmpty() || (maxt < hb.mint || hb.maxt < mint) {
		return iter.NoopIterator
	}

	chunkStats := stats.GetChunkData(ctx)

	// We are doing a copy everytime, this is because b.entries could change completely,
	// the alternate would be that we allocate a new b.entries everytime we cut a block,
	// but the tradeoff is that queries to near-realtime data would be much lower than
	// cutting of blocks.
	chunkStats.HeadChunkLines += int64(len(hb.entries))
	streams := map[uint64]*logproto.Stream{}
	for _, e := range hb.entries {
		chunkStats.HeadChunkBytes += int64(len(e.s))
		line := []byte(e.s)
		newLine, parsedLbs, ok := pipeline.Process(line)
		if !ok {
			continue
		}
		var stream *logproto.Stream
		lhash := parsedLbs.Hash()
		if stream, ok = streams[lhash]; !ok {
			stream = &logproto.Stream{
				Labels: parsedLbs.String(),
			}
			streams[lhash] = stream
		}
		stream.Entries = append(stream.Entries, logproto.Entry{
			Timestamp: time.Unix(0, e.t),
			Line:      string(newLine),
		})

	}

	if len(streams) == 0 {
		return iter.NoopIterator
	}
	streamsResult := make([]logproto.Stream, 0, len(streams))
	for _, stream := range streams {
		streamsResult = append(streamsResult, *stream)
	}
	return iter.NewStreamsIterator(ctx, streamsResult, direction)
}

func (hb *headBlock) sampleIterator(ctx context.Context, mint, maxt int64, extractor log.StreamSampleExtractor) iter.SampleIterator {
	if hb.isEmpty() || (maxt < hb.mint || hb.maxt < mint) {
		return iter.NoopIterator
	}
	chunkStats := stats.GetChunkData(ctx)
	chunkStats.HeadChunkLines += int64(len(hb.entries))
	series := map[uint64]*logproto.Series{}
	for _, e := range hb.entries {
		chunkStats.HeadChunkBytes += int64(len(e.s))
		line := []byte(e.s)
		value, parsedLabels, ok := extractor.Process(line)
		if !ok {
			continue
		}
		var found bool
		var s *logproto.Series
		lhash := parsedLabels.Hash()
		if s, found = series[lhash]; !found {
			s = &logproto.Series{
				Labels: parsedLabels.String(),
			}
			series[lhash] = s
		}
		s.Samples = append(s.Samples, logproto.Sample{
			Timestamp: e.t,
			Value:     value,
			Hash:      xxhash.Sum64([]byte(e.s)),
		})
	}

	if len(series) == 0 {
		return iter.NoopIterator
	}
	seriesRes := make([]logproto.Series, 0, len(series))
	for _, s := range series {
		sort.Sort(s)
		seriesRes = append(seriesRes, *s)
	}
	return iter.NewMultiSeriesIterator(ctx, seriesRes)
}

type bufferedIterator struct {
	block *block
	stats *stats.ChunkData

	bufReader *bufio.Reader
	reader    io.Reader
	pool      ReaderPool

	err error

	decBuf   []byte // The buffer for decoding the lengths.
	buf      []byte // The buffer for a single entry.
	currLine []byte // the current line, this is the same as the buffer but sliced the the line size.
	currTs   int64

	closed bool
}
