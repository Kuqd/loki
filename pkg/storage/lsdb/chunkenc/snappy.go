package chunkenc

import (
	"bytes"
	"encoding/binary"

	"github.com/golang/snappy"
)

type SnappyChunk struct {
	buff *bytes.Buffer
}

func NewSnappyChunk() *SnappyChunk {
	buff := bytes.NewBuffer(make([]byte, 4, 1024)) // 4 first bytes contains the num of lines
	return &SnappyChunk{
		buff:   buff,
		writer: snappy.NewWriter(buff.),
	}
}

// Encoding returns the encoding type.
func (c *SnappyChunk) Encoding() Encoding {
	return EncSnappy
}

// Bytes returns the underlying byte slice of the chunk.
func (c *SnappyChunk) Bytes() []byte {
	return c.buff.Bytes()
}

// NumSamples returns the number of samples in the chunk.
func (c *SnappyChunk) NumSamples() uint32 {
	return binary.BigEndian.Uint32(c.Bytes())
}

type snappyAppender struct {
	writer *snappy.Writer
}

func (a *snappyAppender) Append(int64, string) {

}

// Appender implements the Chunk interface.
func (c *SnappyChunk) Appender() (Appender, error) {
	// it := c.iterator(nil)

	// // To get an appender we must know the state it would have if we had
	// // appended all existing data from scratch.
	// // We iterate through the end and populate via the iterator's state.
	// for it.Next() {
	// }
	// if err := it.Err(); err != nil {
	// 	return nil, err
	// }

	// a := &xorAppender{
	// 	b:        &c.b,
	// 	t:        it.t,
	// 	v:        it.val,
	// 	tDelta:   it.tDelta,
	// 	leading:  it.leading,
	// 	trailing: it.trailing,
	// }
	// if binary.BigEndian.Uint16(a.b.bytes()) == 0 {
	// 	a.leading = 0xff
	// }
	return nil, nil
}
