package distributor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/prometheus/pkg/pool"
	"github.com/weaveworks/common/httpgrpc"
	"github.com/weaveworks/common/httpgrpc/server"

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

var EntryPool = pool.New(256, 2048, 2, func(size int) interface{} { return make([]logproto.Entry, 0, size) })

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
		if err := ParseProtoReader(r.Context(), r.Body, int(r.ContentLength), math.MaxInt32, &req); err != nil {
			return nil, err
		}
	}
	return &req, nil
}

// ParseProtoReader parses a compressed proto from an io.Reader.
func ParseProtoReader(ctx context.Context, reader io.ReadCloser, expectedSize, maxSize int, req proto.Message) error {
	var err error
	sp := opentracing.SpanFromContext(ctx)
	if sp != nil {
		sp.LogFields(otlog.String("event", "util.ParseProtoRequest[start reading]"))
	}

	var buf bytes.Buffer
	if bodyBuff, ok := server.BufferFromReader(reader); ok {
		buf = *bodyBuff
	} else {
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
	}

	if sp != nil {
		sp.LogFields(otlog.String("event", "util.ParseProtoRequest[decompress]"),
			otlog.Int("size", buf.Len()))
	}
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

	if sp != nil {
		sp.LogFields(otlog.String("event", "util.ParseProtoRequest[unmarshal]"),
			otlog.Int("size", len(body)))
	}

	// We re-implement proto.Unmarshal here as it calls XXX_Unmarshal first,
	// which we can't override without upsetting golint.
	req.Reset()
	if u, ok := req.(proto.Unmarshaler); ok {
		err = u.Unmarshal(body)
	} else {
		err = proto.NewBuffer(body).Unmarshal(req)
	}
	if err != nil {
		return err
	}

	return nil
}
