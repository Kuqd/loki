package tsdb

import (
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
	Append(ref uint64, l labels.Labels, t int64, line []byte) (uint64, error)

	// Commit submits the collected samples and purges the batch. If Commit
	// returns a non-nil error, it also rolls back all modifications made in
	// the appender so far, as Rollback would do. In any case, an Appender
	// must not be used anymore after Commit has been called.
	Commit() error

	// Rollback rolls back all modifications made in the appender so far.
	// Appender has to be discarded after rollback.
	Rollback() error

	storage.ExemplarAppender
}
