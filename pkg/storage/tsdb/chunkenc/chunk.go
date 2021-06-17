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
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/golang/snappy"
)

type MemChunk struct {
	// This is compressed bytes.
	b          *bytes.Buffer
	numEntries int

	// mint, maxt int64
	writer  *snappy.Writer
	encBuff []byte
}

func NewMemChunk() *MemChunk {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	writer := snappy.NewBufferedWriter(buffer)
	return &MemChunk{
		writer:  writer,
		b:       buffer,
		encBuff: make([]byte, binary.MaxVarintLen64),
	}
}

func (c *MemChunk) Bytes() []byte {
	c.writer.Flush()
	return c.b.Bytes()
}
func (_ *MemChunk) Encoding() Encoding { return 0 }

func (c *MemChunk) Appender() (Appender, error) {
	return c, nil
}

func (c *MemChunk) Iterator(_ Iterator) Iterator {
	c.writer.Flush()
	reader := bufio.NewReader(snappy.NewReader(c.b))
	return &SnappyIterator{
		reader: reader,
	}
}

func (c *MemChunk) NumSamples() int {
	return c.numEntries
}

func (_ *MemChunk) Compact() {
	fmt.Fprintln(os.Stdout, "Compact")
}

func (c *MemChunk) Append(t int64, v []byte) {
	c.numEntries++
	buf := make([]byte, binary.MaxVarintLen64)
	c.writer.Write(buf[:binary.PutVarint(buf, t)])
	c.writer.Write(buf[:binary.PutVarint(buf, int64(len(v)))])
	c.writer.Write(v)
}

type Reader interface {
	io.Reader
	io.ByteReader
}

type SnappyIterator struct {
	reader Reader
	err    error

	ts   int64
	line []byte
}

// Next advances the iterator by one.
func (it *SnappyIterator) Next() bool {
	ts, err := binary.ReadVarint(it.reader)
	if err != nil && err != io.EOF {
		it.err = err
		return false
	}
	le, err := binary.ReadVarint(it.reader)
	if err != nil && err != io.EOF {
		it.err = err
		return false
	}
	if it.line == nil || int(le) > cap(it.line) {
		it.line = make([]byte, le)
	}
	if _, err := it.reader.Read(it.line[:le]); err != nil && err != io.EOF {
		it.err = err
		return false
	}
	if err == io.EOF {
		return false
	}
	it.ts = ts
	return true
}

func (it *SnappyIterator) Seek(t int64) bool {
	return true
}

func (it *SnappyIterator) At() (int64, []byte) {
	return it.ts, it.line
}

func (it *SnappyIterator) Err() error { return it.err }
