// Copyright 2017 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chunkenc

import (
	"math"
)

// Encoding is the identifier for a chunk encoding.
type Encoding uint8

func (e Encoding) String() string {
	switch e {
	case EncNone:
		return "none"
	case EncXOR:
		return "XOR"
	}
	return "<unknown>"
}

// The different available chunk encodings.
const (
	EncNone Encoding = iota
	EncXOR
)

// Chunk holds a sequence of sample pairs that can be iterated over and appended to.
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
	Iterator(Iterator) Iterator

	// NumSamples returns the number of samples in the chunk.
	NumSamples() int

	// Compact is called whenever a chunk is expected to be complete (no more
	// samples appended) and the underlying implementation can eventually
	// optimize the chunk.
	// There's no strong guarantee that no samples will be appended once
	// Compact() is called. Implementing this function is optional.
	Compact()
}

// Appender adds sample pairs to a chunk.
type Appender interface {
	Append(int64, []byte)
}

// Iterator is a simple iterator that can only get the next value.
// Iterator iterates over the samples of a time series, in timestamp-increasing order.
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
