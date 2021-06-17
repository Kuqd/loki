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
	"bytes"
	"sync"
)

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
				return NewMemChunk()
			},
		},
	}
}

func (p *chunkPool) Get(e Encoding, b []byte) (Chunk, error) {
	c := p.chunks.Get().(*MemChunk)
	c.b = bytes.NewBuffer(b)
	return c, nil
}

func (p *chunkPool) Put(c Chunk) error {
	return nil
}
