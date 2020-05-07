package logql

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"

	"github.com/grafana/loki/pkg/iter"
	"github.com/grafana/loki/pkg/logproto"
)

const (
	ValueName = "value"
)

type Extractor interface {
	Fields() []string
	Values(string) []string
}

func newPreviewExtractor(e Extractor, iter iter.EntryIterator) iter.EntryIterator {
	return &extractorPreview{
		iter:    iter,
		extract: e,
		fields:  e.Fields(),
	}
}

type extractorPreview struct {
	iter    iter.EntryIterator
	extract Extractor

	fields []string
	cur    logproto.Entry
	sb     strings.Builder
}

func (e *extractorPreview) Next() bool {
	for e.iter.Next() {
		entry := e.iter.Entry()
		values := e.extract.Values(entry.Line)
		if len(values) == 0 {
			continue
		}
		e.cur.Timestamp = entry.Timestamp
		e.sb.Reset()
		for i, f := range e.fields {
			e.sb.WriteString(f)
			e.sb.WriteString("=")
			e.sb.WriteString(values[i])
			if i != len(e.fields) {
				e.sb.WriteString(" ")
			}
		}
		e.cur.Line = e.sb.String()
		return true
	}
	return false
}
func (e *extractorPreview) Error() error          { return e.iter.Error() }
func (e *extractorPreview) Labels() string        { return e.iter.Labels() }
func (e *extractorPreview) Entry() logproto.Entry { return e.cur }
func (e *extractorPreview) Close() error          { return e.iter.Close() }

func newExtractor(op, params string) (Extractor, error) {
	switch op {
	case OpExtractRegexp:
		return newRegexpExtractor(params)
	default:
		return nil, fmt.Errorf("unsupported extractor %s", op)
	}
}

func newRegexpExtractor(params string) (*regexpExtractor, error) {
	re, err := regexp.Compile(params)
	if err != nil {
		return nil, err
	}

	f := re.SubexpNames()[1:]
	if len(f) == 0 {
		return nil, errors.New("at least one named capture must be supplied")
	}
	return &regexpExtractor{
		re:     re,
		fields: f,
	}, nil
}

type regexpExtractor struct {
	re     *regexp.Regexp
	fields []string
}

func (e *regexpExtractor) Fields() []string {
	return e.fields
}
func (e *regexpExtractor) Values(line string) []string {
	m := e.re.FindStringSubmatch(line)
	if len(m) <= 1 {
		return nil
	}
	return m[1:]
}

type SeriesIterator interface {
	Close() error
	Next() bool
	Peek() (Sample, bool)
}

type Sample struct {
	Labels        string
	Value         float64
	TimestampNano int64
}

type seriesIterator struct {
	iter    iter.PeekingEntryIterator
	sampler SampleExtractor

	updated bool
	cur     Sample
}

func (e *seriesIterator) Close() error {
	return e.iter.Close()
}

func (e *seriesIterator) Next() bool {
	e.updated = false
	return e.iter.Next()
}

func (e *seriesIterator) Peek() (Sample, bool) {
	if e.updated {
		return e.cur, true
	}

	for {
		lbs, entry, ok := e.iter.Peek()
		if !ok {
			return Sample{}, false
		}

		// transform
		e.cur, ok = e.sampler.From(lbs, entry)
		if ok {
			break
		}
		if !e.iter.Next() {
			return Sample{}, false
		}
	}
	e.updated = true
	return e.cur, true
}

func newSeriesIterator(expr *rangeAggregationExpr, it iter.EntryIterator) (SeriesIterator, error) {

	ex, err := expr.left.left.Extractor()
	if err != nil {
		return nil, err
	}

	if ex == nil && (expr.operation == OpRangeTypeRate || expr.operation == OpRangeTypeCount) {
		return &seriesIterator{
			iter:    iter.NewPeekingIterator(it),
			sampler: defaultSampleExtractor,
		}, nil
	}

	if expr.operation == OpRangeTypeRate || expr.operation == OpRangeTypeCount {
		return &seriesIterator{
			iter:    iter.NewPeekingIterator(it),
			sampler: newLabelsExtractor(ex),
		}, nil
	}
	se, err := newSampleExtractor(ex)
	if err != nil {
		return nil, err
	}
	return &seriesIterator{
		iter:    iter.NewPeekingIterator(it),
		sampler: se,
	}, nil
}

type SampleExtractor interface {
	From(string, logproto.Entry) (Sample, bool)
}

type countSampleExtractor struct{}

var defaultSampleExtractor = countSampleExtractor{}

func (countSampleExtractor) From(lbs string, entry logproto.Entry) (Sample, bool) {
	return Sample{
		Labels:        lbs,
		TimestampNano: entry.Timestamp.UnixNano(),
		Value:         1.,
	}, true
}

type labelsExtractor struct {
	Extractor

	fields []string
}

func newLabelsExtractor(ex Extractor) *labelsExtractor {
	return &labelsExtractor{
		Extractor: ex,
		fields:    ex.Fields(),
	}
}

func (l *labelsExtractor) From(lbs string, entry logproto.Entry) (Sample, bool) {
	values := l.Values(entry.Line)
	if len(values) == 0 {
		return Sample{}, false
	}
	metric, err := promql.ParseMetric(lbs)
	if err != nil {
		//todo warn
		return Sample{}, false
	}
	// todo this is inefficient and we don't validate if those values can be labels
	m := metric.Map()
	for i, f := range l.fields {
		m[f] = values[i]
	}
	return Sample{
		Labels:        labels.FromMap(m).String(),
		TimestampNano: entry.Timestamp.UnixNano(),
		Value:         1,
	}, true
}

type sampleExtractor struct {
	Extractor

	fields []string
}

func newSampleExtractor(ex Extractor) (*sampleExtractor, error) {
	fields := ex.Fields()
	var hasValue bool
	for _, f := range fields {
		if f == ValueName {
			hasValue = true
			break
		}
	}
	if !hasValue {
		return nil, fmt.Errorf("one named captured must be named '%s' to extract sample data", ValueName)
	}
	return &sampleExtractor{
		Extractor: ex,
		fields:    fields,
	}, nil
}

func (l *sampleExtractor) From(lbs string, entry logproto.Entry) (Sample, bool) {
	var err error
	values := l.Values(entry.Line)
	if len(values) == 0 {
		return Sample{}, false
	}
	metric, err := promql.ParseMetric(lbs)
	if err != nil {
		//todo warn
		return Sample{}, false
	}
	// todo this is inefficient and we don't validate if those values can be labels
	m := metric.Map()
	var value float64

	for i, f := range l.fields {
		if f == ValueName {
			value, err = strconv.ParseFloat(values[i], 64)
			if err != nil {
				return Sample{}, false
			}
			continue
		}
		m[f] = values[i]
	}
	return Sample{
		Labels:        labels.FromMap(m).String(),
		TimestampNano: entry.Timestamp.UnixNano(),
		Value:         value,
	}, true
}
