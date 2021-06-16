package chunkenc

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/grafana/loki/pkg/iter"
	"github.com/grafana/loki/pkg/logql/log"
)

// Errors returned by the chunk interface.
var (
	ErrChunkFull       = errors.New("chunk full")
	ErrOutOfOrder      = errors.New("entry out of order")
	ErrInvalidSize     = errors.New("invalid size")
	ErrInvalidFlag     = errors.New("invalid flag")
	ErrInvalidChecksum = errors.New("invalid chunk checksum")
)

// Encoding is the identifier for a chunk encoding.
type Encoding byte

// The different available encodings.
// Make sure to preserve the order, as these numeric values are written to the chunks!
const (
	EncNone Encoding = iota
	EncGZIP
	EncDumb
	EncLZ4_64k
	EncSnappy
	EncLZ4_256k
	EncLZ4_1M
	EncLZ4_4M
	EncFlate
	EncZstd
)

var supportedEncoding = []Encoding{
	EncNone,
	EncGZIP,
	EncLZ4_64k,
	EncSnappy,
	EncLZ4_256k,
	EncLZ4_1M,
	EncLZ4_4M,
	EncFlate,
	EncZstd,
}

func (e Encoding) String() string {
	switch e {
	case EncGZIP:
		return "gzip"
	case EncNone:
		return "none"
	case EncDumb:
		return "dumb"
	case EncLZ4_64k:
		return "lz4-64k"
	case EncLZ4_256k:
		return "lz4-256k"
	case EncLZ4_1M:
		return "lz4-1M"
	case EncLZ4_4M:
		return "lz4"
	case EncSnappy:
		return "snappy"
	case EncFlate:
		return "flate"
	case EncZstd:
		return "zstd"
	default:
		return "unknown"
	}
}

// ParseEncoding parses an chunk encoding (compression algorithm) by its name.
func ParseEncoding(enc string) (Encoding, error) {
	for _, e := range supportedEncoding {
		if strings.EqualFold(e.String(), enc) {
			return e, nil
		}
	}
	return 0, fmt.Errorf("invalid encoding: %s, supported: %s", enc, SupportedEncoding())
}

// SupportedEncoding returns the list of supported Encoding.
func SupportedEncoding() string {
	var sb strings.Builder
	for i := range supportedEncoding {
		sb.WriteString(supportedEncoding[i].String())
		if i != len(supportedEncoding)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

// Chunk is the interface for the compressed logs chunk format.
type Chunk interface {
	// Bytes returns the underlying byte slice of the chunk.
	Bytes() []byte

	// Encoding returns the encoding type of the chunk.
	Encoding() Encoding

	// Appender returns an appender to append samples to the chunk.
	Appender() (Appender, error)

	// The iterator passed as argument is for re-use.
	// Depending on implementation, the iterator can
	// be re-used or a new iterator can be allocated.
	Iterator(iter.EntryIterator) iter.EntryIterator

	// NumSamples returns the number of samples in the chunk.
	NumSamples() int

	// Compact is called whenever a chunk is expected to be complete (no more
	// samples appended) and the underlying implementation can eventually
	// optimize the chunk.
	// There's no strong guarantee that no samples will be appended once
	// Compact() is called. Implementing this function is optional.
	Compact()
}

// Block is a chunk block.
type Block interface {
	// MinTime is the minimum time of entries in the block
	MinTime() int64
	// MaxTime is the maximum time of entries in the block
	MaxTime() int64
	// Offset is the offset/position of the block in the chunk. Offset is unique for a given block per chunk.
	Offset() int
	// Entries is the amount of entries in the block.
	Entries() int
	// Iterator returns an entry iterator for the block.
	Iterator(ctx context.Context, pipeline log.StreamPipeline) iter.EntryIterator
	// SampleIterator returns a sample iterator for the block.
	SampleIterator(ctx context.Context, extractor log.StreamSampleExtractor) iter.SampleIterator
}

type Appender interface {
	Append(t int64, v []byte)
}

type Iterator interface {
	// Next advances the iterator by one.
	Next() bool

	// Seek advances the iterator forward to the first sample with the timestamp equal or greater than t.
	// If current sample found by previous `Next` or `Seek` operation already has this property, Seek has no effect.
	// Seek returns true, if such sample exists, false otherwise.
	// Iterator is exhausted when the Seek returns false.
	Seek(t int64) bool

	// At returns the current timestamp/value pair.
	// Before the iterator has advanced At behaviour is unspecified.
	At() (int64, []byte)
	// Err returns the current error. It should be used only after iterator is
	// exhausted, that is `Next` or `Seek` returns false.
	Err() error
}

// NewNopIterator returns a new chunk iterator that does not hold any data.
func NewNopIterator() Iterator {
	return nopIterator{}
}

type nopIterator struct{}

func (nopIterator) Seek(int64) bool     { return false }
func (nopIterator) At() (int64, []byte) { return math.MinInt64, nil }
func (nopIterator) Next() bool          { return false }
func (nopIterator) Err() error          { return nil }

// Pool is used to create and reuse chunk references to avoid allocations.
type Pool interface {
	Put(Chunk) error
	Get(e Encoding, b []byte) (Chunk, error)
}

// pool is a memory pool of chunk objects.
type chunkPool struct {
	chunks sync.Pool
}

// NewPool returns a new pool.
func NewPool() Pool {
	return &chunkPool{
		chunks: sync.Pool{
			New: func() interface{} {
				return NewMemChunk(EncSnappy, DefaultBlockSize, DefaultTargetSize)
			},
		},
	}
}

func (p *chunkPool) Get(e Encoding, b []byte) (Chunk, error) {
	c := p.chunks.Get().(*MemChunk)
	return c, nil
}

func (p *chunkPool) Put(c Chunk) error {
	return nil
}
