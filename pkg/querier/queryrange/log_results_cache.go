package queryrange

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	proto "github.com/gogo/protobuf/proto"

	"github.com/cortexproject/cortex/pkg/chunk/cache"
	"github.com/cortexproject/cortex/pkg/querier/queryrange"
)

type logResultsCache struct {
	logger   log.Logger
	cfg      queryrange.ResultsCacheConfig
	next     queryrange.Handler
	cache    cache.Cache
	limits   Limits
	splitter queryrange.CacheSplitter
}

func NewLogResultsCacheMiddleware(
	logger log.Logger,
	cfg queryrange.ResultsCacheConfig,
	splitter queryrange.CacheSplitter,
	limits Limits,
) (queryrange.Middleware, Stopper, error) {
	c, err := cache.New(cfg.CacheConfig)
	if err != nil {
		return nil, nil, err
	}

	return queryrange.MiddlewareFunc(func(next queryrange.Handler) queryrange.Handler {
		return &logResultsCache{
			logger:   logger,
			cfg:      cfg,
			next:     next,
			cache:    c,
			limits:   limits,
			splitter: splitter,
		}
	}), c, nil
}

func (s logResultsCache) Do(ctx context.Context, r queryrange.Request) (queryrange.Response, error) {
	return s.next.Do(ctx, r)
}

func (s logResultsCache) put(ctx context.Context, key string, extents []queryrange.Extent) {
	buf, err := proto.Marshal(&queryrange.CachedResponse{
		Key:     key,
		Extents: extents,
	})
	if err != nil {
		level.Error(s.logger).Log("msg", "error marshalling cached value", "err", err)
		return
	}

	s.cache.Store(ctx, []string{cache.HashKey(key)}, [][]byte{buf})
}
