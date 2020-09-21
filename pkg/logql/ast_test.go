package logql

import (
	"strings"
	"testing"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_logSelectorExpr_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		selector     string
		expectFilter bool
	}{
		{`{foo!~"bar"}`, false},
		{`{foo="bar", bar!="baz"}`, false},
		{`{foo="bar", bar!="baz"} != "bip" !~ ".+bop"`, true},
		{`{foo="bar"} |= "baz" |~ "blip" != "flip" !~ "flap"`, true},
		{`{foo="bar", bar!="baz"} |= ""`, false},
		{`{foo="bar", bar!="baz"} |~ ""`, false},
		{`{foo="bar", bar!="baz"} |~ ".*"`, false},
		{`{foo="bar", bar!="baz"} |= "" |= ""`, false},
		{`{foo="bar", bar!="baz"} |~ "" |= "" |~ ".*"`, false},
		{`{foo="bar", bar!="baz"} != "bip" !~ ".+bop" | json`, true},
		{`{foo="bar"} |= "baz" |~ "blip" != "flip" !~ "flap" | logfmt`, true},
		{`{foo="bar"} |= "baz" |~ "blip" != "flip" !~ "flap" | regexp "(?P<foo>foo|bar)"`, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.selector, func(t *testing.T) {
			t.Parallel()
			expr, err := ParseLogSelector(tt.selector)
			if err != nil {
				t.Fatalf("failed to parse log selector: %s", err)
			}
			p, err := expr.Pipeline()
			if err != nil {
				t.Fatalf("failed to get filter: %s", err)
			}
			require.Equal(t, tt.expectFilter, p != NoopPipeline)
			if expr.String() != strings.Replace(tt.selector, " ", "", -1) {
				t.Fatalf("error expected: %s got: %s", tt.selector, expr.String())
			}
		})
	}
}

func Test_SampleExpr_String(t *testing.T) {
	t.Parallel()
	for _, tc := range []string{
		`rate( ( {job="mysql"} |="error" !="timeout" ) [10s] )`,
		`sum without(a) ( rate ( ( {job="mysql"} |="error" !="timeout" ) [10s] ) )`,
		`sum by(a) (rate( ( {job="mysql"} |="error" !="timeout" ) [10s] ) )`,
		`sum(count_over_time({job="mysql"}[5m]))`,
		`sum(count_over_time({job="mysql"} | json [5m]))`,
		`sum(count_over_time({job="mysql"} | logfmt [5m]))`,
		`sum(count_over_time({job="mysql"} | regexp "(?P<foo>foo|bar)" [5m]))`,
		`topk(10,sum(rate({region="us-east1"}[5m])) by (name))`,
		`avg( rate( ( {job="nginx"} |= "GET" ) [10s] ) ) by (region)`,
		`sum by (cluster) (count_over_time({job="mysql"}[5m]))`,
		`sum by (cluster) (count_over_time({job="mysql"}[5m])) / sum by (cluster) (count_over_time({job="postgres"}[5m])) `,
		`
		sum by (cluster) (count_over_time({job="postgres"}[5m])) /
		sum by (cluster) (count_over_time({job="postgres"}[5m])) /
		sum by (cluster) (count_over_time({job="postgres"}[5m]))
		`,
		`sum by (cluster) (count_over_time({job="mysql"}[5m])) / min(count_over_time({job="mysql"}[5m])) `,
		`sum by (job) (
			count_over_time({namespace="tns"} |= "level=error"[5m])
		/
			count_over_time({namespace="tns"}[5m])
		)`,
	} {
		t.Run(tc, func(t *testing.T) {
			expr, err := ParseExpr(tc)
			require.Nil(t, err)

			expr2, err := ParseExpr(expr.String())
			require.Nil(t, err)
			require.Equal(t, expr, expr2)
		})
	}
}

func Test_NilFilterDoesntPanic(t *testing.T) {
	t.Parallel()
	for _, tc := range []string{
		`{namespace="dev", container_name="cart"} |= "" |= "bloop"`,
		`{namespace="dev", container_name="cart"} |= "bleep" |= ""`,
		`{namespace="dev", container_name="cart"} |= "bleep" |= "" |= "bloop"`,
		`{namespace="dev", container_name="cart"} |= "bleep" |= "" |= "bloop"`,
		`{namespace="dev", container_name="cart"} |= "bleep" |= "bloop" |= ""`,
	} {
		t.Run(tc, func(t *testing.T) {
			expr, err := ParseLogSelector(tc)
			require.Nil(t, err)

			p, err := expr.Pipeline()
			require.Nil(t, err)
			_, _, ok := p.Process([]byte("bleepbloop"), labelBar)

			require.True(t, ok)
		})
	}

}

type linecheck struct {
	l string
	e bool
}

func Test_FilterMatcher(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		q string

		expectedMatchers []*labels.Matcher
		// test line against the resulting filter, if empty filter should also be nil
		lines []linecheck
	}{
		{
			`{app="foo",cluster=~".+bar"}`,
			[]*labels.Matcher{
				mustNewMatcher(labels.MatchEqual, "app", "foo"),
				mustNewMatcher(labels.MatchRegexp, "cluster", ".+bar"),
			},
			nil,
		},
		{
			`{app!="foo",cluster=~".+bar",bar!~".?boo"}`,
			[]*labels.Matcher{
				mustNewMatcher(labels.MatchNotEqual, "app", "foo"),
				mustNewMatcher(labels.MatchRegexp, "cluster", ".+bar"),
				mustNewMatcher(labels.MatchNotRegexp, "bar", ".?boo"),
			},
			nil,
		},
		{
			`{app="foo"} |= "foo"`,
			[]*labels.Matcher{
				mustNewMatcher(labels.MatchEqual, "app", "foo"),
			},
			[]linecheck{{"foobar", true}, {"bar", false}},
		},
		{
			`{app="foo"} |= "foo" != "bar"`,
			[]*labels.Matcher{
				mustNewMatcher(labels.MatchEqual, "app", "foo"),
			},
			[]linecheck{{"foobuzz", true}, {"bar", false}},
		},
		{
			`{app="foo"} |= "foo" !~ "f.*b"`,
			[]*labels.Matcher{
				mustNewMatcher(labels.MatchEqual, "app", "foo"),
			},
			[]linecheck{{"foo", true}, {"bar", false}, {"foobar", false}},
		},
		{
			`{app="foo"} |= "foo" |~ "f.*b"`,
			[]*labels.Matcher{
				mustNewMatcher(labels.MatchEqual, "app", "foo"),
			},
			[]linecheck{{"foo", false}, {"bar", false}, {"foobar", true}},
		},
		{
			`{app="foo"} |~ "foo"`,
			[]*labels.Matcher{
				mustNewMatcher(labels.MatchEqual, "app", "foo"),
			},
			[]linecheck{{"foo", true}, {"bar", false}, {"foobar", true}},
		},
	} {
		tt := tt
		t.Run(tt.q, func(t *testing.T) {
			t.Parallel()
			expr, err := ParseLogSelector(tt.q)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedMatchers, expr.Matchers())
			p, err := expr.Pipeline()
			assert.Nil(t, err)
			if tt.lines == nil {
				assert.Equal(t, p, NoopPipeline)
			} else {
				for _, lc := range tt.lines {
					_, _, ok := p.Process([]byte(lc.l), labelBar)
					assert.Equal(t, lc.e, ok)
				}
			}
		})
	}
}

func TestStringer(t *testing.T) {
	for _, tc := range []struct {
		in  string
		out string
	}{
		{
			in:  `1 > 1 > 1`,
			out: `0`,
		},
		{
			in:  `1.6`,
			out: `1.6`,
		},
		{
			in:  `1 > 1 > bool 1`,
			out: `0`,
		},
		{
			in:  `1 > bool 1 > count_over_time({foo="bar"}[1m])`,
			out: `0 > count_over_time({foo="bar"}[1m])`,
		},
		{
			in:  `1 > bool 1 > bool count_over_time({foo="bar"}[1m])`,
			out: `0 > bool count_over_time({foo="bar"}[1m])`,
		},
		{

			in:  `0 > count_over_time({foo="bar"}[1m])`,
			out: `0 > count_over_time({foo="bar"}[1m])`,
		},
	} {
		t.Run(tc.in, func(t *testing.T) {
			expr, err := ParseExpr(tc.in)
			require.Nil(t, err)
			require.Equal(t, tc.out, expr.String())
		})
	}
}

func BenchmarkContainsFilter(b *testing.B) {
	expr, err := ParseLogSelector(`{app="foo"} |= "foo"`)
	if err != nil {
		b.Fatal(err)
	}

	p, err := expr.Pipeline()
	if err != nil {
		b.Fatal(err)
	}

	line := []byte("hello world foo bar")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, _, ok := p.Process(line, labelBar); !ok {
			b.Fatal("doesn't match")
		}
	}
}

func Test_parserExpr_Parser(t *testing.T) {
	tests := []struct {
		name    string
		op      string
		param   string
		want    LabelParser
		wantErr bool
	}{
		{"json", OpParserTypeJSON, "", NewJSONParser(), false},
		{"logfmt", OpParserTypeLogfmt, "", NewLogfmtParser(), false},
		{"regexp", OpParserTypeRegexp, "(?P<foo>foo)", mustNewRegexParser("(?P<foo>foo)"), false},
		{"regexp err ", OpParserTypeRegexp, "foo", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &labelParserExpr{
				op:    tt.op,
				param: tt.param,
			}
			got, err := e.parser()
			if (err != nil) != tt.wantErr {
				t.Errorf("parserExpr.Parser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				require.Nil(t, got)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
