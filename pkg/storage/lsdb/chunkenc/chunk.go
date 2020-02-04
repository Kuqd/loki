package chunkenc

// Encoding is the identifier for a chunk encoding.
type Encoding uint8

func (e Encoding) String() string {
	switch e {
	case EncNone:
		return "none"
	case EncSnappy:
		return "snappy"
	}
	return "<unknown>"
}

// The different available chunk encodings.
const (
	EncNone Encoding = iota
	EncSnappy
)

type Chunk interface {
	Bytes() []byte
	Encoding() Encoding
	Appender() (Appender, error)
	// The iterator passed as argument is for re-use.
	// Depending on implementation, the iterator can
	// be re-used or a new iterator can be allocated.
	Iterator(Iterator) Iterator
	NumSamples() uint32
}

// Appender adds sample pairs to a chunk.
type Appender interface {
	Append(int64, string)
	Close() error
}

// Iterator is a simple iterator that can only get the next value.
type Iterator interface {
	At() (int64, string)
	Err() error
	Next() bool
}
