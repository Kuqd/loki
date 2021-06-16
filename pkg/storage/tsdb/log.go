package tsdb

import (
	"github.com/grafana/loki/pkg/storage/tsdb/chunkenc"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
)

type LogAppender interface {
	// Append adds a sample pair for the given series.
	// An optional reference number can be provided to accelerate calls.
	// A reference number is returned which can be used to add further
	// samples in the same or later transactions.
	// Returned reference numbers are ephemeral and may be rejected in calls
	// to Append() at any point. Adding the sample via Append() returns a new
	// reference number.
	// If the reference is 0 it must not be used for caching.
	AppendLog(ref uint64, l labels.Labels, t int64, line []byte) (uint64, error)

	storage.Appender
}

type LogChunkAppender interface {
	chunkenc.Appender
	AppendLog(t int64, v []byte)
}
