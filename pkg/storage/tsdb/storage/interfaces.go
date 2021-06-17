package storage

import (
	"context"

	"github.com/grafana/loki/pkg/storage/tsdb/chunkenc"
	"github.com/grafana/loki/pkg/storage/tsdb/chunks"
	"github.com/prometheus/prometheus/pkg/labels"
)

type ChunkQuerier interface {
	LabelQuerier

	// Select returns a set of series that matches the given label matchers.
	// Caller can specify if it requires returned series to be sorted. Prefer not requiring sorting for better performance.
	// It allows passing hints that can help in optimising select, but it's up to implementation how this is used if used at all.
	Select(sortSeries bool, hints *SelectHints, matchers ...*labels.Matcher) ChunkSeriesSet
}

// ChunkSeriesSet contains a set of chunked series.
type ChunkSeriesSet interface {
	Next() bool
	// At returns full chunk series. Returned series should be iterable even after Next is called.
	At() ChunkSeries
	// The error that iteration has failed with.
	// When an error occurs, set cannot continue to iterate.
	Err() error
	// A collection of warnings for the whole set.
	// Warnings could be return even iteration has not failed with error.
	Warnings() Warnings
}

// ChunkSeries exposes a single time series and allows iterating over chunks.
type ChunkSeries interface {
	Labels
	ChunkIterable
}

// Labels represents an item that has labels e.g. time series.
type Labels interface {
	// Labels returns the complete set of labels. For series it means all labels identifying the series.
	Labels() labels.Labels
}

type ChunkIterable interface {
	// Iterator returns a new, independent iterator that iterates over potentially overlapping
	// chunks of the series, sorted by min time.
	Iterator() chunks.Iterator
}

type Warnings []error

// LabelQuerier provides querying access over labels.
type LabelQuerier interface {
	// LabelValues returns all potential values for a label name.
	// It is not safe to use the strings beyond the lifefime of the querier.
	// If matchers are specified the returned result set is reduced
	// to label values of metrics matching the matchers.
	LabelValues(name string, matchers ...*labels.Matcher) ([]string, Warnings, error)

	// LabelNames returns all the unique label names present in the block in sorted order.
	// TODO(yeya24): support matchers or hints.
	LabelNames() ([]string, Warnings, error)

	// Close releases the resources of the Querier.
	Close() error
}

// Querier provides querying access over time series data of a fixed time range.
type Querier interface {
	LabelQuerier

	// Select returns a set of series that matches the given label matchers.
	// Caller can specify if it requires returned series to be sorted. Prefer not requiring sorting for better performance.
	// It allows passing hints that can help in optimising select, but it's up to implementation how this is used if used at all.
	Select(sortSeries bool, hints *SelectHints, matchers ...*labels.Matcher) SeriesSet
}

// SelectHints specifies hints passed for data selections.
// This is used only as an option for implementation to use.
type SelectHints struct {
	Start int64 // Start time in milliseconds for this select.
	End   int64 // End time in milliseconds for this select.

	Step int64  // Query step size in milliseconds.
	Func string // String representation of surrounding function or aggregation.

	Grouping []string // List of label names used in aggregation.
	By       bool     // Indicate whether it is without or by.
	Range    int64    // Range vector selector range in milliseconds.
}

// SeriesSet contains a set of series.
type SeriesSet interface {
	Next() bool
	// At returns full series. Returned series should be iterable even after Next is called.
	At() Series
	// The error that iteration as failed with.
	// When an error occurs, set cannot continue to iterate.
	Err() error
	// A collection of warnings for the whole set.
	// Warnings could be return even iteration has not failed with error.
	Warnings() Warnings
}

// Series exposes a single time series and allows iterating over samples.
type Series interface {
	Labels
	SampleIterable
}

type SampleIterable interface {
	// Iterator returns a new, independent iterator of the data of the series.
	Iterator() chunkenc.Iterator
}

// Appender provides batched appends against a storage.
// It must be completed with a call to Commit or Rollback and must not be reused afterwards.
//
// Operations on the Appender interface are not goroutine-safe.
type Appender interface {
	// Append adds a sample pair for the given series.
	// An optional reference number can be provided to accelerate calls.
	// A reference number is returned which can be used to add further
	// samples in the same or later transactions.
	// Returned reference numbers are ephemeral and may be rejected in calls
	// to Append() at any point. Adding the sample via Append() returns a new
	// reference number.
	// If the reference is 0 it must not be used for caching.
	Append(ref uint64, l labels.Labels, t int64, v []byte) (uint64, error)

	// Commit submits the collected samples and purges the batch. If Commit
	// returns a non-nil error, it also rolls back all modifications made in
	// the appender so far, as Rollback would do. In any case, an Appender
	// must not be used anymore after Commit has been called.
	Commit() error

	// Rollback rolls back all modifications made in the appender so far.
	// Appender has to be discarded after rollback.
	Rollback() error
}

// Appendable allows creating appenders.
type Appendable interface {
	// Appender returns a new appender for the storage. The implementation
	// can choose whether or not to use the context, for deadlines or to check
	// for errors.
	Appender(ctx context.Context) Appender
}

// SampleAndChunkQueryable allows retrieving samples as well as encoded samples in form of chunks.
type SampleAndChunkQueryable interface {
	Queryable
	ChunkQueryable
}

// A Queryable handles queries against a storage.
// Use it when you need to have access to all samples without chunk encoding abstraction e.g promQL.
type Queryable interface {
	// Querier returns a new Querier on the storage.
	Querier(ctx context.Context, mint, maxt int64) (Querier, error)
}

// A ChunkQueryable handles queries against a storage.
// Use it when you need to have access to samples in encoded format.
type ChunkQueryable interface {
	// ChunkQuerier returns a new ChunkQuerier on the storage.
	ChunkQuerier(ctx context.Context, mint, maxt int64) (ChunkQuerier, error)
}

// GetRef is an extra interface on Appenders used by downstream projects
// (e.g. Cortex) to avoid maintaining a parallel set of references.
type GetRef interface {
	// Returns reference number that can be used to pass to Appender.Append(),
	// and a set of labels that will not cause another copy when passed to Appender.Append().
	// 0 means the appender does not have a reference to this series.
	GetRef(lset labels.Labels) (uint64, labels.Labels)
}

type errSeriesSet struct {
	err error
}

func (s errSeriesSet) Next() bool         { return false }
func (s errSeriesSet) At() Series         { return nil }
func (s errSeriesSet) Err() error         { return s.err }
func (s errSeriesSet) Warnings() Warnings { return nil }

// ErrSeriesSet returns a series set that wraps an error.
func ErrSeriesSet(err error) SeriesSet {
	return errSeriesSet{err: err}
}

type errChunkSeriesSet struct {
	err error
}

func (s errChunkSeriesSet) Next() bool         { return false }
func (s errChunkSeriesSet) At() ChunkSeries    { return nil }
func (s errChunkSeriesSet) Err() error         { return s.err }
func (s errChunkSeriesSet) Warnings() Warnings { return nil }

// ErrChunkSeriesSet returns a chunk series set that wraps an error.
func ErrChunkSeriesSet(err error) ChunkSeriesSet {
	return errChunkSeriesSet{err: err}
}
