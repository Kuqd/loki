package queryrange

import (
	"context"
	"net/http"
	"strings"

	"github.com/grafana/loki/pkg/loghttp"
	"github.com/grafana/loki/pkg/logql"

	"github.com/cortexproject/cortex/pkg/querier/queryrange"
	"github.com/go-kit/kit/log"

	"github.com/cortexproject/cortex/pkg/querier/frontend"
)

// NewTripperware returns a Tripperware configured with middlewares to align, split and cache requests.
func NewTripperware(cfg queryrange.Config, log log.Logger, limits queryrange.Limits) (frontend.Tripperware, error) {
	metricsTripperware, err := queryrange.NewTripperware(cfg, log, limits, &lokiCodec{}, queryrange.PrometheusResponseExtractor)
	if err != nil {
		return nil, err
	}
	return frontend.Tripperware(func(next http.RoundTripper) http.RoundTripper {
		metricRT := metricsTripperware(next)
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
				return metricRT.RoundTrip(req)
			}
			return next.RoundTrip(req)
		})
	}), nil
}

type lokiCodec struct {
	queryrange.Codec
}

type request struct {
	*loghttp.RangeQuery
}

func (l *lokiCodec) ParseRequest(_ context.Context, r *http.Request) (queryrange.Request, error) {
	req, err := loghttp.ParseRangeQuery(r)
	if err != nil {
		return nil, err
	}
	return &queryrange.PrometheusRequest{
		Query:,
		Path:  r.URL.Path,
	}, nil
}
