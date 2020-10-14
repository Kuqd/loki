package log

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/prometheus/prometheus/pkg/labels"
)

var (
	_ LabelFilterer = &BinaryLabelFilter{}
	_ LabelFilterer = &BytesLabelFilter{}
	_ LabelFilterer = &DurationLabelFilter{}
	_ LabelFilterer = &NumericLabelFilter{}
	_ LabelFilterer = &StringLabelFilter{}

	NoopLabelFilter = noopLabelFilter{}
)

// LabelFilterType is an enum for label filtering types.
type LabelFilterType int

// Possible LabelFilterType.
const (
	LabelFilterEqual LabelFilterType = iota
	LabelFilterNotEqual
	LabelFilterGreaterThan
	LabelFilterGreaterThanOrEqual
	LabelFilterLesserThan
	LabelFilterLesserThanOrEqual
)

func (f LabelFilterType) String() string {
	switch f {
	case LabelFilterEqual:
		return "=="
	case LabelFilterNotEqual:
		return "!="
	case LabelFilterGreaterThan:
		return ">"
	case LabelFilterGreaterThanOrEqual:
		return ">="
	case LabelFilterLesserThan:
		return "<"
	case LabelFilterLesserThanOrEqual:
		return "<="
	default:
		return ""
	}
}

type LabelFilterer interface {
	Stage
	fmt.Stringer
}

type BinaryLabelFilter struct {
	Left  LabelFilterer
	Right LabelFilterer
	and   bool
}

func NewAndLabelFilter(left LabelFilterer, right LabelFilterer) *BinaryLabelFilter {
	return &BinaryLabelFilter{
		Left:  left,
		Right: right,
		and:   true,
	}
}

func NewOrLabelFilter(left LabelFilterer, right LabelFilterer) *BinaryLabelFilter {
	return &BinaryLabelFilter{
		Left:  left,
		Right: right,
	}
}

func (b *BinaryLabelFilter) Process(line []byte, lbs Labels) ([]byte, bool) {
	line, lok := b.Left.Process(line, lbs)
	if !b.and && lok {
		return line, true
	}
	line, rok := b.Right.Process(line, lbs)
	if !b.and {
		return line, lok || rok
	}
	return line, lok && rok
}

func (b *BinaryLabelFilter) String() string {
	var sb strings.Builder
	sb.WriteString("( ")
	sb.WriteString(b.Left.String())
	if b.and {
		sb.WriteString(" , ")
	} else {
		sb.WriteString(" or ")
	}
	sb.WriteString(b.Right.String())
	sb.WriteString(" )")
	return sb.String()
}

type noopLabelFilter struct{}

func (noopLabelFilter) String() string                                 { return "" }
func (noopLabelFilter) Process(line []byte, lbs Labels) ([]byte, bool) { return line, true }

func ReduceAndLabelFilter(filters []LabelFilterer) LabelFilterer {
	if len(filters) == 0 {
		return NoopLabelFilter
	}
	if len(filters) == 1 {
		return filters[0]
	}
	result := filters[0]
	for _, f := range filters[1:] {
		result = NewAndLabelFilter(result, f)
	}
	return result
}

type BytesLabelFilter struct {
	Name  string
	Value uint64
	Type  LabelFilterType
}

func NewBytesLabelFilter(t LabelFilterType, name string, b uint64) *BytesLabelFilter {
	return &BytesLabelFilter{
		Name:  name,
		Type:  t,
		Value: b,
	}
}

func (d *BytesLabelFilter) Process(line []byte, lbs Labels) ([]byte, bool) {
	if lbs.HasError() {
		// if there's an error only the string matchers can filter it out.
		return line, true
	}
	v, ok := lbs[d.Name]
	if !ok {
		// we have not found this label.
		return line, false
	}
	value, err := humanize.ParseBytes(v)
	if err != nil {
		lbs.SetError(errLabelFilter)
		return line, true
	}
	switch d.Type {
	case LabelFilterEqual:
		return line, value == d.Value
	case LabelFilterNotEqual:
		return line, value != d.Value
	case LabelFilterGreaterThan:
		return line, value > d.Value
	case LabelFilterGreaterThanOrEqual:
		return line, value >= d.Value
	case LabelFilterLesserThan:
		return line, value < d.Value
	case LabelFilterLesserThanOrEqual:
		return line, value <= d.Value
	default:
		lbs.SetError(errLabelFilter)
		return line, true
	}
}

func (d *BytesLabelFilter) String() string {
	return fmt.Sprintf("%s%s%d", d.Name, d.Type, d.Value)
}

type DurationLabelFilter struct {
	Name  string
	Value time.Duration
	Type  LabelFilterType
}

func NewDurationLabelFilter(t LabelFilterType, name string, d time.Duration) *DurationLabelFilter {
	return &DurationLabelFilter{
		Name:  name,
		Type:  t,
		Value: d,
	}
}

func (d *DurationLabelFilter) Process(line []byte, lbs Labels) ([]byte, bool) {
	if lbs.HasError() {
		// if there's an error only the string matchers can filter out.
		return line, true
	}
	v, ok := lbs[d.Name]
	if !ok {
		// we have not found this label.
		return line, false
	}
	value, err := time.ParseDuration(v)
	if err != nil {
		lbs.SetError(errLabelFilter)
		return line, true
	}
	switch d.Type {
	case LabelFilterEqual:
		return line, value == d.Value
	case LabelFilterNotEqual:
		return line, value != d.Value
	case LabelFilterGreaterThan:
		return line, value > d.Value
	case LabelFilterGreaterThanOrEqual:
		return line, value >= d.Value
	case LabelFilterLesserThan:
		return line, value < d.Value
	case LabelFilterLesserThanOrEqual:
		return line, value <= d.Value
	default:
		lbs.SetError(errLabelFilter)
		return line, true
	}
}

func (d *DurationLabelFilter) String() string {
	return fmt.Sprintf("%s%s%s", d.Name, d.Type, d.Value)
}

type NumericLabelFilter struct {
	Name  string
	Value float64
	Type  LabelFilterType
}

func NewNumericLabelFilter(t LabelFilterType, name string, v float64) *NumericLabelFilter {
	return &NumericLabelFilter{
		Name:  name,
		Type:  t,
		Value: v,
	}
}

func (n *NumericLabelFilter) Process(line []byte, lbs Labels) ([]byte, bool) {
	if lbs.HasError() {
		// if there's an error only the string matchers can filter out.
		return line, true
	}
	if lbs.HasError() {
		// if there's an error only the string matchers can filter out.
		return line, true
	}
	v, ok := lbs[n.Name]
	if !ok {
		// we have not found this label.
		return line, false
	}
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		lbs.SetError(errLabelFilter)
		return line, true
	}
	switch n.Type {
	case LabelFilterEqual:
		return line, value == n.Value
	case LabelFilterNotEqual:
		return line, value != n.Value
	case LabelFilterGreaterThan:
		return line, value > n.Value
	case LabelFilterGreaterThanOrEqual:
		return line, value >= n.Value
	case LabelFilterLesserThan:
		return line, value < n.Value
	case LabelFilterLesserThanOrEqual:
		return line, value <= n.Value
	default:
		lbs.SetError(errLabelFilter)
		return line, true
	}

}

func (n *NumericLabelFilter) String() string {
	return fmt.Sprintf("%s%s%s", n.Name, n.Type, strconv.FormatFloat(n.Value, 'f', -1, 64))
}

type StringLabelFilter struct {
	*labels.Matcher
}

func NewStringLabelFilter(m *labels.Matcher) *StringLabelFilter {
	return &StringLabelFilter{
		Matcher: m,
	}
}

func (s *StringLabelFilter) Process(line []byte, lbs Labels) ([]byte, bool) {
	return line, s.Matches(lbs[s.Name])
}
