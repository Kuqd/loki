package queryrange

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/cortexproject/cortex/pkg/querier/queryrange"
	"github.com/cortexproject/cortex/pkg/util"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/grafana/loki/pkg/logql/marshal"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/stretchr/testify/require"
	"github.com/weaveworks/common/middleware"
	"github.com/weaveworks/common/user"
)

var now = time.Now()

func TestMetricsTripperware(t *testing.T) {

	cfg := queryrange.Config{
		SplitQueriesByInterval: 4 * time.Hour,
		AlignQueriesWithStep:   true,
		MaxRetries:             3,
	}

	tpw, err := NewTripperware(cfg, util.Logger, fakeLimits{})
	require.NoError(t, err)

	lreq := &LokiRequest{
		Query:     `rate({app="foo"} |= "foo"[1m])`,
		Limit:     1000,
		Step:      30000, //30sec
		StartTs:   now.Add(-6 * time.Hour),
		EndTs:     now,
		Direction: logproto.FORWARD,
		Path:      "/query_range",
	}

	ctx := user.InjectOrgID(context.Background(), "1")
	req, err := lokiCodec.EncodeRequest(ctx, lreq)
	require.NoError(t, err)

	req = req.WithContext(ctx)
	err = user.InjectOrgIDIntoHTTPRequest(ctx, req)
	require.NoError(t, err)

	count, h := counter()
	_, err = tpw(newfakeRoundTripper(t, h)).RoundTrip(req)
	// 2 split so 6 retries
	require.Equal(t, 6, *count)
	require.Error(t, err)

	count, h = promqlResult(promql.Matrix{
		{
			Points: []promql.Point{
				{
					T: 1568404331324,
					V: 0.013333333333333334,
				},
			},
			Metric: []labels.Label{
				{
					Name:  "filename",
					Value: `/var/hostlog/apport.log`,
				},
				{
					Name:  "job",
					Value: "varlogs",
				},
			},
		},
		{
			Points: []promql.Point{
				{
					T: 1568404331324,
					V: 3.45,
				},
				{
					T: 1568404331339,
					V: 4.45,
				},
			},
			Metric: []labels.Label{
				{
					Name:  "filename",
					Value: `/var/hostlog/syslog`,
				},
				{
					Name:  "job",
					Value: "varlogs",
				},
			},
		},
	})
	_, err = tpw(newfakeRoundTripper(t, h)).RoundTrip(req)
	// 2 queries todo test result merge.
	require.Equal(t, 2, *count)
	require.NoError(t, err)
}

type fakeLimits struct{}

func (fakeLimits) MaxQueryLength(string) time.Duration {
	return 0 // Disable.
}

func (fakeLimits) MaxQueryParallelism(string) int {
	return 14 // Flag default.
}

func counter() (*int, http.Handler) {
	count := 0
	return &count, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
	})
}

func promqlResult(v promql.Value) (*int, http.Handler) {
	count := 0
	return &count, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := marshal.WriteQueryResponseJSON(v, w); err != nil {
			panic(err)
		}
		count++
	})

}

type fakeRoundTripper struct {
	*httptest.Server
	host string
}

func newfakeRoundTripper(t *testing.T, handler http.Handler) *fakeRoundTripper {
	s := httptest.NewServer(
		middleware.AuthenticateUser.Wrap(
			handler,
		),
	)

	u, err := url.Parse(s.URL)
	require.NoError(t, err)
	return &fakeRoundTripper{
		Server: s,
		host:   u.Host,
	}
}

func (s fakeRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = "http"
	r.URL.Host = s.host
	return http.DefaultTransport.RoundTrip(r)
}
