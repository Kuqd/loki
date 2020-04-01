package queryrange

import (
	"context"
	"net/http"
	"sort"
	"time"

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
	"github.com/cortexproject/cortex/pkg/util/spanlogger"
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
	if ok {
		response, extents, err = s.handleHit(ctx, r, cached)
	} else {
		// response, extents, err = s.handleMiss(ctx, r)
	}

	if err == nil && len(extents) > 0 {
		// extents, err := s.filterRecentExtents(r, extents)
		// if err != nil {
		// 	return nil, err
		// }
		s.put(ctx, key, extents)
	}

	return response, err
}

type Extent interface {
	Start() time.Time
	End() time.Time
	TraceID() string
	Extract(start, end time.Time) (queryrange.Response, error)
}

func (s logResultsCache) handleHit(ctx context.Context, r queryrange.Request, extents []queryrange.Extent) (queryrange.Response, []queryrange.Extent, error) {
	var (
		reqResps []queryrange.RequestResponse
		err      error
	)
	log, ctx := spanlogger.New(ctx, "handleHit")
	defer log.Finish()

	requests, responses, err := partition(r, extents, s.extractor)
	if err != nil {
		return nil, nil, err
	}
	if len(requests) == 0 {
		response, err := s.merger.MergeResponse(responses...)
		// No downstream requests so no need to write back to the cache.
		return response, nil, err
	}

	reqResps, err = DoRequests(ctx, s.next, requests, s.limits)
	if err != nil {
		return nil, nil, err
	}

	for _, reqResp := range reqResps {
		responses = append(responses, reqResp.Response)
		extent, err := toExtent(ctx, reqResp.Request, reqResp.Response)
		if err != nil {
			return nil, nil, err
		}
		extents = append(extents, extent)
	}
	sort.Slice(extents, func(i, j int) bool {
		return extents[i].Start < extents[j].Start
	})

	// Merge any extents - they're guaranteed not to overlap.
	accumulator, err := newAccumulator(extents[0])
	if err != nil {
		return nil, nil, err
	}
	mergedExtents := make([]Extent, 0, len(extents))

	for i := 1; i < len(extents); i++ {
		if accumulator.End+r.GetStep() < extents[i].Start {
			mergedExtents, err = merge(mergedExtents, accumulator)
			if err != nil {
				return nil, nil, err
			}
			accumulator, err = newAccumulator(extents[i])
			if err != nil {
				return nil, nil, err
			}
			continue
		}

		accumulator.TraceId = jaegerTraceID(ctx)
		accumulator.End = extents[i].End
		currentRes, err := extents[i].toResponse()
		if err != nil {
			return nil, nil, err
		}
		merged, err := s.merger.MergeResponse(accumulator.Response, currentRes)
		if err != nil {
			return nil, nil, err
		}
		accumulator.Response = merged
	}

	mergedExtents, err = merge(mergedExtents, accumulator)
	if err != nil {
		return nil, nil, err
	}

	response, err := s.merger.MergeResponse(responses...)
	return response, mergedExtents, err
}

// partition calculates the required requests to satisfy req given the cached data.
func partition(req queryrange.Request, extents []Extent) ([]queryrange.Request, []queryrange.Response, error) {
	var requests []queryrange.Request
	var cachedResponses []queryrange.Response
	start := req.GetStart()

	for _, extent := range extents {
		// If there is no overlap, ignore this extent.
		if extent.End() < start || extent.Start() > req.GetEnd() {
			continue
		}

		// If there is a bit missing at the front, make a request for that.
		if start < extent.Start {
			r := req.WithStartEnd(start, extent.Start)
			requests = append(requests, r)
		}
		res, err := extent.toResponse()
		if err != nil {
			return nil, nil, err
		}
		// extract the overlap from the cached extent.
		cachedResponses = append(cachedResponses, extractor.Extract(start, req.GetEnd(), res))
		start = extent.End
	}

	if start < req.GetEnd() {
		r := req.WithStartEnd(start, req.GetEnd())
		requests = append(requests, r)
	}

	return requests, cachedResponses, nil
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
