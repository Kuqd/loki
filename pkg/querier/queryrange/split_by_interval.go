package queryrange

import (
	"context"
	"time"

	"github.com/cortexproject/cortex/pkg/querier/queryrange"
)

// SplitByIntervalMiddleware creates a new Middleware that splits log requests by a given interval.
func SplitByIntervalMiddleware(interval time.Duration, limits queryrange.Limits, merger queryrange.Merger) queryrange.Middleware {
	return queryrange.MiddlewareFunc(func(next queryrange.Handler) queryrange.Handler {
		return splitByInterval{
			next:     next,
			limits:   limits,
			merger:   merger,
			interval: interval,
		}
	})
}

type splitByInterval struct {
	next     queryrange.Handler
	limits   queryrange.Limits
	merger   queryrange.Merger
	interval time.Duration
}

func (s splitByInterval) Do(ctx context.Context, r queryrange.Request) (queryrange.Response, error) {
	reqs := splitQuery(r, s.interval)

	reqResps, err := queryrange.DoRequests(ctx, s.next, reqs, s.limits)
	if err != nil {
		return nil, err
	}

	resps := make([]queryrange.Response, 0, len(reqResps))
	for _, reqResp := range reqResps {
		resps = append(resps, reqResp.Response)
	}

	response, err := s.merger.MergeResponse(resps...)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func splitQuery(r queryrange.Request, interval time.Duration) []queryrange.Request {
	lokiRequest := r.(*LokiRequest)
	var reqs []queryrange.Request
	for start := lokiRequest.StartTs; start.Before(lokiRequest.EndTs); start = start.Add(interval) {
		end := start.Add(interval)
		if end.After(lokiRequest.EndTs) {
			end = lokiRequest.EndTs
		}
		reqs = append(reqs, &LokiRequest{
			Query:     lokiRequest.Query,
			Limit:     lokiRequest.Limit,
			Step:      lokiRequest.Step,
			Direction: lokiRequest.Direction,
			Path:      lokiRequest.Path,
			StartTs:   start,
			EndTs:     end,
		})
	}
	return reqs
}
