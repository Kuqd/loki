package chunkenc

import (
	"context"
	"errors"
	"time"

	"github.com/cortexproject/cortex/pkg/chunk"

	"github.com/grafana/loki/pkg/iter"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/grafana/loki/pkg/logql"
)

// LazyChunk loads the chunk when it is accessed.
type LazyChunk struct {
	Chunk   chunk.Chunk
	Fetcher *chunk.Fetcher

	MinT, MaxT time.Time

	cache *iter.CachedIterator
}

// Iterator returns an entry iterator.
func (c *LazyChunk) Iterator(ctx context.Context, from, through time.Time, direction logproto.Direction, filter logql.LineFilter) (iter.EntryIterator, error) {
	// If the chunk is already loaded, then use that.
	if c.Chunk.Data == nil {
		return nil, errors.New("chunk is not loaded")

	}

	lokiChunk := c.Chunk.Data.(*Facade).LokiChunk()

	// for filter queries we expect to run through many chunks if it doesn't match a lot.
	// So we cache filtered chunk entries to avoid to decompress multiple time the same chunk.
	if filter != nil {
		if c.cache == nil {
			// the underlaying iterator boundaries is clipped with the min and max of the request
			// since on the next iteration from and through can be different and we want to make
			// sure we cache everything we will need.
			it, err := lokiChunk.Iterator(ctx, c.MinT, c.MaxT, direction, filter)
			if err != nil {
				return nil, err
			}
			c.cache = iter.NewCachedIterator(it)
		}
		return iter.NewTimeRangedIterator(c.cache, from, through, direction), nil
	}
	return lokiChunk.Iterator(ctx, from, through, direction, filter)
}
