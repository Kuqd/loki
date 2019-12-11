package queryrange

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/cortexproject/cortex/pkg/ingester/client"
	"github.com/cortexproject/cortex/pkg/querier/queryrange"
	"github.com/grafana/loki/pkg/loghttp"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/grafana/loki/pkg/logql/marshal"
	json "github.com/json-iterator/go"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/weaveworks/common/httpgrpc"
)

var lokiCodec = &codec{}

type codec struct{}

func (r *LokiRequest) GetEnd() int64 {
	return r.EndTs.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func (r *LokiRequest) GetStart() int64 {
	return r.StartTs.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func (r *LokiRequest) WithStartEnd(s int64, e int64) queryrange.Request {
	new := *r
	new.StartTs = time.Unix(0, s*int64(time.Millisecond))
	new.EndTs = time.Unix(0, e*int64(time.Millisecond))
	return &new
}

func (codec) DecodeRequest(_ context.Context, r *http.Request) (queryrange.Request, error) {
	req, err := loghttp.ParseRangeQuery(r)
	if err != nil {
		return nil, err
	}
	return &LokiRequest{
		Query:     req.Query,
		Limit:     req.Limit,
		Direction: req.Direction,
		StartTs:   req.Start,
		EndTs:     req.End,
		// GetStep must return milliseconds
		Step: int64(req.Step) / 1e6,
		Path: r.URL.Path,
	}, nil
}

func (codec) EncodeRequest(ctx context.Context, r queryrange.Request) (*http.Request, error) {
	lokiReq, ok := r.(*LokiRequest)
	if !ok {
		return nil, httpgrpc.Errorf(http.StatusInternalServerError, "invalid request format")
	}
	params := url.Values{
		"start":     []string{fmt.Sprintf("%d", lokiReq.StartTs.UnixNano())},
		"end":       []string{fmt.Sprintf("%d", lokiReq.EndTs.UnixNano())},
		"query":     []string{lokiReq.Query},
		"direction": []string{lokiReq.Direction.String()},
		"limit":     []string{fmt.Sprintf("%d", lokiReq.Limit)},
	}
	if lokiReq.Step != 0 {
		params["step"] = []string{fmt.Sprintf("%d", lokiReq.Step/int64(1e3))}
	}
	u := &url.URL{
		Path:     lokiReq.Path,
		RawQuery: params.Encode(),
	}
	req := &http.Request{
		Method:     "GET",
		RequestURI: u.String(), // This is what the httpgrpc code looks at.
		URL:        u,
		Body:       http.NoBody,
		Header:     http.Header{},
	}

	return req.WithContext(ctx), nil
}

func (codec) DecodeResponse(ctx context.Context, r *http.Response, req queryrange.Request) (queryrange.Response, error) {
	if r.StatusCode/100 != 2 {
		body, _ := ioutil.ReadAll(r.Body)
		return nil, httpgrpc.Errorf(r.StatusCode, string(body))
	}

	sp, _ := opentracing.StartSpanFromContext(ctx, "DecodeResponse")
	defer sp.Finish()

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sp.LogFields(otlog.Error(err))
		return nil, httpgrpc.Errorf(http.StatusInternalServerError, "error decoding response: %v", err)
	}

	sp.LogFields(otlog.Int("bytes", len(buf)))

	var resp loghttp.QueryResponse
	if err := json.Unmarshal(buf, &resp); err != nil {
		return nil, httpgrpc.Errorf(http.StatusInternalServerError, "error decoding response: %v", err)
	}
	if resp.Status != loghttp.QueryStatusSuccess {
		return nil, httpgrpc.Errorf(http.StatusInternalServerError, "error executing request: %v", resp.Status)
	}
	switch string(resp.Data.ResultType) {
	case loghttp.ResultTypeMatrix:
		return &queryrange.PrometheusResponse{
			Status: loghttp.QueryStatusSuccess,
			Data: queryrange.PrometheusData{
				ResultType: loghttp.ResultTypeMatrix,
				Result:     toProto(resp.Data.Result.(loghttp.Matrix)),
			},
		}, nil
	case loghttp.ResultTypeStream:
		return &LokiResponse{
			Status:    loghttp.QueryStatusSuccess,
			Direction: req.(*LokiRequest).Direction,
			Data: LokiData{
				ResultType: loghttp.ResultTypeStream,
				Result:     resp.Data.Result.(loghttp.Streams).ToProto(),
			},
		}, nil
	default:
		return nil, httpgrpc.Errorf(http.StatusBadRequest, "unsupported response type")
	}
}

func (codec) EncodeResponse(ctx context.Context, res queryrange.Response) (*http.Response, error) {
	sp, _ := opentracing.StartSpanFromContext(ctx, "APIResponse.ToHTTPResponse")
	defer sp.Finish()

	if _, ok := res.(*queryrange.PrometheusResponse); ok {
		return queryrange.PrometheusCodec.EncodeResponse(ctx, res)
	}

	proto, ok := res.(*LokiResponse)
	if !ok {
		return nil, httpgrpc.Errorf(http.StatusInternalServerError, "invalid response format")
	}

	streams := make(loghttp.Streams, len(proto.Data.Result))

	for i, stream := range proto.Data.Result {
		s, err := marshal.NewStream(&stream)
		if err != nil {
			return nil, err
		}
		streams[i] = s
	}

	queryRes := loghttp.QueryResponse{
		Status: proto.Status,
		Data: loghttp.QueryResponseData{
			ResultType: loghttp.ResultType(proto.Data.ResultType),
			Result:     streams,
		},
	}

	b, err := json.Marshal(queryRes)
	if err != nil {
		return nil, err
	}

	sp.LogFields(otlog.Int("bytes", len(b)))

	resp := http.Response{
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
		StatusCode: http.StatusOK,
	}
	return &resp, nil
}

type byMinTime []*LokiResponse

func (a byMinTime) Len() int           { return len(a) }
func (a byMinTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byMinTime) Less(i, j int) bool { return a[i].fistTime() < a[j].fistTime() }

type byMaxTime []*LokiResponse

func (a byMaxTime) Len() int           { return len(a) }
func (a byMaxTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byMaxTime) Less(i, j int) bool { return a[i].fistTime() > a[j].fistTime() }

func (resp *LokiResponse) fistTime() int64 {
	result := resp.Data.Result
	if len(result) == 0 {
		return -1
	}
	if len(result[0].Entries) == 0 {
		return -1
	}
	return result[0].Entries[0].Timestamp.UnixNano()
}

func (codec) MergeResponse(responses ...queryrange.Response) (queryrange.Response, error) {
	if len(responses) == 0 {
		return nil, errors.New("merging responses requires at least one response")
	}
	if _, ok := responses[0].(*queryrange.PrometheusResponse); ok {
		return queryrange.PrometheusCodec.MergeResponse(responses...)
	}
	lokiRes, ok := responses[0].(*LokiResponse)
	if !ok {
		return nil, errors.New("unexpected response type while merging")
	}

	lokiResponses := make([]*LokiResponse, 0, len(responses))
	for _, res := range responses {
		lokiResponses = append(lokiResponses, res.(*LokiResponse))
	}
	// todo limits results
	if lokiRes.Direction == logproto.FORWARD {
		sort.Sort(byMinTime(lokiResponses))
	} else {
		sort.Sort(byMaxTime(lokiResponses))
	}

	return &LokiResponse{
		Status:    loghttp.QueryStatusSuccess,
		Direction: lokiRes.Direction,
		Data: LokiData{
			ResultType: loghttp.ResultTypeStream,
			Result:     mergeStreams(lokiResponses),
		},
	}, nil
}

func mergeStreams(resps []*LokiResponse) []logproto.Stream {
	output := map[string]*logproto.Stream{}
	for _, resp := range resps {
		for _, stream := range resp.Data.Result {
			lbs := stream.Labels
			existing, ok := output[lbs]
			if !ok {
				existing = &logproto.Stream{
					Labels: lbs,
				}
			}
			existing.Entries = append(existing.Entries, stream.Entries...)
			output[lbs] = existing
		}
	}

	keys := make([]string, 0, len(output))
	for key := range output {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make([]logproto.Stream, 0, len(output))
	for _, key := range keys {
		result = append(result, *output[key])
	}

	return result
}

func toProto(m loghttp.Matrix) []queryrange.SampleStream {
	if len(m) == 0 {
		return nil
	}
	res := make([]queryrange.SampleStream, 0, len(m))
	for _, stream := range m {
		samples := make([]client.Sample, 0, len(stream.Values))
		for _, s := range stream.Values {
			samples = append(samples, client.Sample{
				Value:       float64(s.Value),
				TimestampMs: int64(s.Timestamp),
			})
		}
		res = append(res, queryrange.SampleStream{
			Labels:  client.FromMetricsToLabelAdapters(stream.Metric),
			Samples: samples,
		})
	}
	return res
}
