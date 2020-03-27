package queryrange

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	proto "github.com/gogo/protobuf/proto"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/common/model"
	"github.com/weaveworks/common/httpgrpc"
	"github.com/weaveworks/common/user"

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
	userID, err := user.ExtractOrgID(ctx)
	if err != nil {
		return nil, httpgrpc.Errorf(http.StatusBadRequest, err.Error())
	}

	var (
		key      = s.splitter.GenerateCacheKey(userID, r)
		extents  []queryrange.Extent
		response queryrange.Response
	)
	maxCacheTime := int64(model.Now().Add(-s.cfg.MaxCacheFreshness))
	if r.GetStart() > maxCacheTime {
		return s.next.Do(ctx, r)
	}

	cached, ok := s.get(ctx, key)

	return s.next.Do(ctx, r)
}

func (s logResultsCache) get(ctx context.Context, key string) ([]queryrange.Extent, bool) {
	found, bufs, _ := s.cache.Fetch(ctx, []string{cache.HashKey(key)})
	if len(found) != 1 {
		return nil, false
	}

	// todo uncompress

	var resp queryrange.CachedResponse
	sp, _ := opentracing.StartSpanFromContext(ctx, "unmarshal-extent")
	defer sp.Finish()

	sp.LogFields(otlog.Int("bytes", len(bufs[0])))

	if err := proto.Unmarshal(bufs[0], &resp); err != nil {
		level.Error(s.logger).Log("msg", "error unmarshalling cached value", "err", err)
		sp.LogFields(otlog.Error(err))
		return nil, false
	}

	if resp.Key != key {
		return nil, false
	}

	// Refreshes the cache if it contains an old proto schema.
	for _, e := range resp.Extents {
		if e.Response == nil {
			return nil, false
		}
	}

	return resp.Extents, true
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
	// todo compress

	s.cache.Store(ctx, []string{cache.HashKey(key)}, [][]byte{buf})
}
