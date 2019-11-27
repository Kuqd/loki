package queryrange

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/cortexproject/cortex/pkg/querier/queryrange"
	"github.com/cortexproject/cortex/pkg/util"
	"github.com/stretchr/testify/require"
	"github.com/weaveworks/common/middleware"
)

func TestNewTripperware(t *testing.T) {

	cfg := queryrange.Config{}

	tpw, err := NewTripperware(cfg, util.Logger, fakeLimits{})
	require.NoError(t, err)

	req, err := http.NewRequest("GET", "/query_range", http.NoBody)
	require.NoError(t, err)

	tpw(newfakeRoundTripper(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`responseBody`))
	}))).RoundTrip(req)

	// type args struct {
	// 	cfg    queryrange.Config
	// 	log    log.Logger
	// 	limits queryrange.Limits
	// }
	// tests := []struct {
	// 	name    string
	// 	args    args
	// 	want    frontend.Tripperware
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		got, err := NewTripperware(tt.args.cfg, tt.args.log, tt.args.limits)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("NewTripperware() error = %v, wantErr %v", err, tt.wantErr)
	// 			return
	// 		}
	// 		got.
	// 		if !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("NewTripperware() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}

type fakeLimits struct{}

func (fakeLimits) MaxQueryLength(string) time.Duration {
	return 0 // Disable.
}

func (fakeLimits) MaxQueryParallelism(string) int {
	return 14 // Flag default.
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
