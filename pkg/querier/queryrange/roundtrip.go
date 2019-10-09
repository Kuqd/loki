package queryrange

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/loki/pkg/logql"
	"github.com/weaveworks/common/httpgrpc"

	"github.com/cortexproject/cortex/pkg/querier/queryrange"
	"github.com/go-kit/kit/log"

	"github.com/cortexproject/cortex/pkg/querier/frontend"
)

// NewTripperware returns a Tripperware configured with middlewares to align, split and cache requests.
func NewTripperware(cfg queryrange.Config, log log.Logger, limits queryrange.Limits) (frontend.Tripperware, error) {
	metricsTripperware, err := queryrange.NewTripperware(cfg, log, limits)
	if err != nil {
		return nil, err
	}
	return frontend.Tripperware(func(next http.RoundTripper) http.RoundTripper {
		return frontend.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if !strings.HasSuffix(req.URL.Path, "/query_range") {
				return next.RoundTrip(req)
			}
			params := req.URL.Query()
			query := params.Get("query")
			expr, err := logql.ParseExpr(query)
			if err != nil {
				return nil, err
			}
			if _, ok := expr.(logql.SampleExpr); ok {
				mRequest, err := toMetricRequest(req)
				if err != nil {
					return nil, err
				}
				httpRequest, err := mRequest.ToHTTPRequest(req.Context())
				if err != nil {
					return nil, err
				}
				return metricsTripperware(next).RoundTrip(httpRequest)
			}
			return next.RoundTrip(req)
		})
	}), nil
	// var queryRangeMiddleware []Middleware
	// if cfg.AlignQueriesWithStep {
	// 	queryRangeMiddleware = append(queryRangeMiddleware, InstrumentMiddleware("step_align"), StepAlignMiddleware)
	// }
	// if cfg.SplitQueriesByDay {
	// 	queryRangeMiddleware = append(queryRangeMiddleware, InstrumentMiddleware("split_by_day"), SplitByDayMiddleware(limits))
	// }
	// if cfg.CacheResults {
	// 	queryCacheMiddleware, err := NewResultsCacheMiddleware(log, cfg.ResultsCacheConfig, limits)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	queryRangeMiddleware = append(queryRangeMiddleware, InstrumentMiddleware("results_cache"), queryCacheMiddleware)
	// }

	// return frontend.Tripperware(func(next http.RoundTripper) http.RoundTripper {
	// 	// Finally, if the user selected any query range middleware, stitch it in.
	// 	if len(queryRangeMiddleware) > 0 {
	// 		return NewRoundTripper(next, MergeMiddlewares(queryRangeMiddleware...).Wrap(&ToRoundTripperMiddleware{Next: next}), limits)
	// 	}
	// 	return next
	// }), nil
}

const defaultSince = 1 * time.Hour

func toMetricRequest(r *http.Request) (*queryrange.Request, error) {
	params := r.URL.Query()
	now := time.Now()

	query := params.Get("query")

	start, err := unixNanoTimeParam(params, "start", now.Add(-defaultSince))
	if err != nil {
		return nil, httpgrpc.Errorf(http.StatusBadRequest, err.Error())
	}

	end, err := unixNanoTimeParam(params, "end", now)
	if err != nil {
		return nil, httpgrpc.Errorf(http.StatusBadRequest, err.Error())
	}

	return &queryrange.Request{
		Path:  "/queryrange",
		Query: query,
	}, nil
}

func unixNanoTimeParam(values url.Values, name string, def time.Time) (time.Time, error) {
	value := values.Get(name)
	if value == "" {
		return def, nil
	}

	if strings.Contains(value, ".") {
		if t, err := strconv.ParseFloat(value, 64); err == nil {
			s, ns := math.Modf(t)
			ns = math.Round(ns*1000) / 1000
			return time.Unix(int64(s), int64(ns*float64(time.Second))), nil
		}
	}
	nanos, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		if ts, err := time.Parse(time.RFC3339Nano, value); err == nil {
			return ts, nil
		}
		return time.Time{}, err
	}
	if len(value) <= 10 {
		return time.Unix(nanos, 0), nil
	}
	return time.Unix(0, nanos), nil
}
