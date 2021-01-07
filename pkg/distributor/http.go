package distributor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/pkg/pool"
	"github.com/weaveworks/common/httpgrpc"

	"github.com/grafana/loki/pkg/loghttp"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/grafana/loki/pkg/logql/unmarshal"
	unmarshal_legacy "github.com/grafana/loki/pkg/logql/unmarshal/legacy"
)

var contentType = http.CanonicalHeaderKey("Content-Type")

const applicationJSON = "application/json"

// PushHandler reads a snappy-compressed proto from the HTTP body.
func (d *Distributor) PushHandler(w http.ResponseWriter, r *http.Request) {

	req, err := ParseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = d.Push(r.Context(), req)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp, ok := httpgrpc.HTTPResponseFromError(err)
	if ok {
		http.Error(w, string(resp.Body), int(resp.Code))
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ParseRequest(r *http.Request) (*logproto.PushRequest, error) {
	var req logproto.PushRequest

	switch r.Header.Get(contentType) {
	case applicationJSON:
		var err error

		if loghttp.GetVersion(r.RequestURI) == loghttp.VersionV1 {
			err = unmarshal.DecodePushRequest(r.Body, &req)
		} else {
			err = unmarshal_legacy.DecodePushRequest(r.Body, &req)
		}

		if err != nil {
			return nil, err
		}

	default:
		if err := ParseProtoRequest(r.Context(), r.Body, int(r.ContentLength), math.MaxInt32, &req); err != nil {
			return nil, err
		}
	}
	return &req, nil
}

var BytesBufferPool = pool.New(1*1024, 2*1024*1024, 2, func(size int) interface{} { return make([]byte, 0, size) })

// ParseProtoRequest parses a compressed proto pushRequest from an io.Reader.
func ParseProtoRequest(ctx context.Context, reader io.Reader, expectedSize, maxSize int, req *logproto.PushRequest) error {
	var err error

	var buf bytes.Buffer

	if expectedSize > 0 {
		if expectedSize > maxSize {
			return fmt.Errorf("message expected size larger than max (%d vs %d)", expectedSize, maxSize)
		}
		bufData := BytesBufferPool.Get(expectedSize + bytes.MinRead).([]byte)
		buf = *bytes.NewBuffer(bufData)
		buf.Reset()
		defer BytesBufferPool.Put(bufData)
	}

	_, err = buf.ReadFrom(reader)

	var body []byte
	var l int
	if err == nil && buf.Len() <= maxSize {
		l, err = snappy.DecodedLen(buf.Bytes())
		if err != nil {
			return err
		}
		body = BytesBufferPool.Get(l).([]byte)
		body = body[:l]
		defer BytesBufferPool.Put(body)
		body, err = snappy.Decode(body, buf.Bytes())
	}

	if err != nil {
		return err
	}
	if len(body) > maxSize {
		return fmt.Errorf("received message larger than max (%d vs %d)", len(body), maxSize)
	}

	req.Reset()
	err = req.Unmarshal(body)
	if err != nil {
		return err
	}

	return nil
}
