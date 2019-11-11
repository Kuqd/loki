package queryrange

import (
	"context"
	"net/http"
	"time"

	"github.com/cortexproject/cortex/pkg/querier/queryrange"
	"github.com/grafana/loki/pkg/loghttp"
	"github.com/grafana/loki/pkg/logproto"
)

type codec struct {
	queryrange.Codec
}

type request struct {
	*logproto.QueryRequest
}

func (r *request) GetEnd() int64 {
	return 0
}

func (r *request) GetStart() int64 {
	return 0
}

func (r *request) GetStep() int64 {
	return int64(0 / time.Millisecond)
}

func (r *request) GetQuery() string {
	return ""
}

func (r *request) WithStartEnd(s int64, e int64) queryrange.Request {
	return r
}

func (c *codec) ParseRequest(_ context.Context, r *http.Request) (queryrange.Request, error) {
	req, err := loghttp.ParseRangeQuery(r)
	if err != nil {
		return nil, err
	}
	return &request{
		QueryRequest: &logproto.QueryRequest{
			Selector:  req.Query,
			Limit:     req.Limit,
			Direction: req.Direction,
			Start:     req.Start,
			End:       req.End,
		},
	}, nil
}
