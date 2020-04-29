package iter

import (
	"context"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/grafana/loki/pkg/logproto"
)

const bufSize = 1024 * 1024

var res = []*logproto.QueryResponse{}

var result = &logproto.QueryResponse{
	Streams: []*logproto.Stream{
		&logproto.Stream{
			Labels: `{foo="bar"}`,
			Entries: []logproto.Entry{
				logproto.Entry{Timestamp: time.Now(), Line: "foo"},
			},
		},
		&logproto.Stream{
			Labels: `{foo="buzz"}`,
			Entries: []logproto.Entry{
				logproto.Entry{Timestamp: time.Now(), Line: "buzz"},
			},
		},
	},
}

func BenchmarkQueryClientIterator(b *testing.B) {

	lis := bufconn.Listen(bufSize)
	server := grpc.NewServer()
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		b.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	ing := ingesterFn(func(req *logproto.QueryRequest, s logproto.Querier_QueryServer) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			_ = s.Send(result)
		}
	})
	logproto.RegisterQuerierServer(server, ing)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	ingClient := logproto.NewQuerierClient(conn)
	// iters := make([]EntryIterator, 10)
	// for i := range iters {
	// 	c, err := ingClient.Query(ctx, &logproto.QueryRequest{})
	// 	if err != nil {
	// 		b.Fatal(err)
	// 	}
	// 	iters[i] = NewQueryClientIterator(c, logproto.FORWARD)
	// }
	// iter := NewHeapIterator(ctx, iters, logproto.FORWARD)
	// defer iter.Close()
	c, err := ingClient.Query(ctx, &logproto.QueryRequest{})
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		// if iter.Next() {
		// 	entries = append(entries, iter.Entry())
		// }

		r, err := c.Recv()
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
		res = append(res, r)
	}
}

type ingesterFn func(*logproto.QueryRequest, logproto.Querier_QueryServer) error

func (i ingesterFn) Query(req *logproto.QueryRequest, s logproto.Querier_QueryServer) error {
	return i(req, s)
}
func (ingesterFn) Label(context.Context, *logproto.LabelRequest) (*logproto.LabelResponse, error) {
	return nil, nil
}
func (ingesterFn) Tail(*logproto.TailRequest, logproto.Querier_TailServer) error { return nil }
func (ingesterFn) Series(context.Context, *logproto.SeriesRequest) (*logproto.SeriesResponse, error) {
	return nil, nil
}
func (ingesterFn) TailersCount(context.Context, *logproto.TailersCountRequest) (*logproto.TailersCountResponse, error) {
	return nil, nil
}
