package queryrange

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/types"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/weaveworks/common/user"

	"github.com/cortexproject/cortex/pkg/chunk/cache"
	"github.com/cortexproject/cortex/pkg/querier/queryrange"
	"github.com/cortexproject/cortex/pkg/util/spanlogger"
)

type ExtentCache interface {
	Get(ctx context.Context, key string) ([]queryrange.Extent, bool)
	Put(ctx context.Context, key string, extents []queryrange.Extent)
}

type resultsCache struct {
	logger log.Logger
	cfg    queryrange.ResultsCacheConfig
	next   queryrange.Handler
	cache  ExtentCache
	limits queryrange.Limits
}

// NewResultsCacheMiddleware creates results cache middleware from config.
// Caches only a limited amount of logs per interval, specified by the request log limit.
// If there is more entries for the interval, the cache will keep newest/oldest logs based on the request direction.
// The rest of the data will be fetched via downstream handler.
func NewResultsCacheMiddleware(logger log.Logger, cfg queryrange.ResultsCacheConfig, limits queryrange.Limits) (queryrange.Middleware, error) {
	c, err := cache.New(cfg.CacheConfig)
	if err != nil {
		return nil, err
	}
	return queryrange.MiddlewareFunc(func(next queryrange.Handler) queryrange.Handler {
		return &resultsCache{
			logger: logger,
			cfg:    cfg,
			next:   next,
			cache:  queryrange.NewExtentCache(c, logger),
			limits: limits,
		}
	}), nil
}

func (s resultsCache) Do(ctx context.Context, r queryrange.Request) (queryrange.Response, error) {
	req, ok := r.(*LokiRequest)
	if !ok {
		return nil, fmt.Errorf("unsupported request type: %T", r)
	}
	userID, err := user.ExtractOrgID(ctx)
	if err != nil {
		return nil, err
	}

	var (
		key      = generateKey(userID, req, s.cfg.SplitInterval)
		extents  []queryrange.Extent
		response queryrange.Response
	)

	maxCacheTime := time.Now().Add(-s.cfg.MaxCacheFreshness)
	if req.StartTs.After(maxCacheTime) {
		return s.next.Do(ctx, r)
	}

	cached, ok := s.cache.Get(ctx, key)
	if ok {
		response, extents, err = s.handleHit(ctx, req, cached)
	} else {
		response, extents, err = s.handleMiss(ctx, req)
	}

	if err == nil && len(extents) > 0 {
		extents, err := s.filterRecentExtents(req, extents)
		if err != nil {
			return nil, err
		}
		s.cache.Put(ctx, key, extents)
	}

	return response, err
}

func (s resultsCache) handleMiss(ctx context.Context, r *LokiRequest) (queryrange.Response, []queryrange.Extent, error) {
	response, err := s.next.Do(ctx, r)
	if err != nil {
		return nil, nil, err
	}

	//todo sets start or end if the result contains the same amount of data as the result.
	extent, err := queryrange.NewExtent(ctx, r.StartTs.UnixNano(), r.EndTs.UnixNano(), response)
	if err != nil {
		return nil, nil, err
	}

	extents := []queryrange.Extent{
		extent,
	}
	return response, extents, nil
}

func (s resultsCache) handleHit(ctx context.Context, r *LokiRequest, extents []queryrange.Extent) (queryrange.Response, []queryrange.Extent, error) {
	var (
		reqResps []queryrange.RequestResponse
		err      error
	)
	log, ctx := spanlogger.New(ctx, "handleHit")
	defer log.Finish()

	requests, responses, err := partition(r, extents)
	if err != nil {
		return nil, nil, err
	}
	if len(requests) == 0 {
		response, err := s.merger.MergeResponse(responses...)
		// No downstream requests so no need to write back to the cache.
		return response, nil, err
	}

	reqResps, err = queryrange.DoRequests(ctx, s.next, requests, s.limits)
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
	mergedExtents := make([]queryrange.Extent, 0, len(extents))

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
		currentRes, err := toResponse(extents[i])
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

type accumulator struct {
	queryrange.Response
	queryrange.Extent
}

func merge(extents []queryrange.Extent, acc *accumulator) ([]queryrange.Extent, error) {
	any, err := types.MarshalAny(acc.Response)
	if err != nil {
		return nil, err
	}
	return append(extents, queryrange.Extent{
		Start:    acc.Extent.Start,
		End:      acc.Extent.End,
		Response: any,
		TraceId:  acc.Extent.TraceId,
	}), nil
}

func newAccumulator(base queryrange.Extent) (*accumulator, error) {
	res, err := toResponse(base)
	if err != nil {
		return nil, err
	}
	return &accumulator{
		Response: res,
		Extent:   base,
	}, nil
}

// partition calculates the required requests to satisfy req given the cached data.
func partition(req *LokiRequest, extents []queryrange.Extent) ([]queryrange.Request, []queryrange.Response, error) {
	var requests []queryrange.Request
	var cachedResponses []queryrange.Response
	start := req.StartTs.UnixNano()

	for _, extent := range extents {
		// If there is no overlap, ignore this extent.
		if extent.GetEnd() < start || extent.GetStart() > req.EndTs.UnixNano() {
			continue
		}

		// If there is a bit missing at the front, make a request for that.
		if start < extent.GetStart() {
			r := *req
			r.StartTs = time.Unix(0, start)
			r.EndTs = time.Unix(0, extent.GetEnd())
			requests = append(requests, &r)
		}
		res, err := extent.ToResponse()
		if err != nil {
			return nil, nil, err
		}
		// extract the overlap from the cached extent.
		cachedResponses = append(cachedResponses, extractLogs(start, req.EndTs.UnixNano(), res))
		start = extent.End
	}

	if start < req.EndTs.UnixNano() {
		r := *req
		r.StartTs = time.Unix(0, start)
		r.EndTs = req.EndTs
		requests = append(requests, &r)
	}

	return requests, cachedResponses, nil
}

func (s resultsCache) filterRecentExtents(req *LokiRequest, extents []queryrange.Extent) ([]queryrange.Extent, error) {
	maxCacheTime := time.Now().Add(-s.cfg.MaxCacheFreshness).UnixNano()
	for i := range extents {
		// Never cache data for the latest freshness period.
		if extents[i].End > maxCacheTime {
			extents[i].End = maxCacheTime
			res, err := extents[i].ToResponse()
			if err != nil {
				return nil, err
			}
			extracted := extractLogs(extents[i].Start, maxCacheTime, res)
			any, err := types.MarshalAny(extracted)
			if err != nil {
				return nil, err
			}
			extents[i].Response = any
		}
	}
	return extents, nil
}

// generateKey generates a cache key based on the userID, Request and interval.
func generateKey(userID string, r *LokiRequest, interval time.Duration) string {

}

func jaegerTraceID(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	spanContext, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return ""
	}

	return spanContext.TraceID().String()
}

func extractLogs(start, end int64, from queryrange.Response) queryrange.Response {
	res := from.(*LokiResponse)
	return &LokiResponse{
		Status: statusSuccess,
		Data: LokiData{
			ResultType: res.Data.ResultType,
			// start and end are passed as milliseconds
			Result: extractStreams(start*1e6, end*1e6, res.Data.Result),
		},
	}
}

func countEntries(streams []logproto.Stream) int64 {
	if len(streams) == 0 {
		return 0
	}
	res := int64(0)
	for _, s := range streams {
		res += int64(len(s.Entries))
	}
	return res
}

func extractStreams(start, end int64, streams []logproto.Stream) []logproto.Stream {
	result := make([]logproto.Stream, 0, len(streams))
	for _, stream := range streams {
		extracted, ok := extractStream(start, end, stream)
		if ok {
			result = append(result, extracted)
		}
	}
	return result
}

func extractStream(start, end int64, stream logproto.Stream) (logproto.Stream, bool) {
	result := logproto.Stream{
		Labels:  stream.Labels,
		Entries: make([]logproto.Entry, 0, len(stream.Entries)),
	}
	for _, entry := range stream.Entries {
		ts := entry.Timestamp.UnixNano()
		if start <= ts && ts < end {
			result.Entries = append(result.Entries, entry)
		}
	}
	if len(result.Entries) == 0 {
		return logproto.Stream{}, false
	}
	return result, true
}
