package regexp

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestFindEqualMatches(t *testing.T) {
	tests := []struct {
		re   string
		want []string
	}{
		{"foo|bar", []string{"foo", "bar"}},
		{"foo|foo|bar", []string{"foo", "bar"}},
		{"foo|foobar", []string{"foo", "foobar"}},
		{"(foo|(ba|ar))", []string{"foo", "ba", "ar"}},
		{"^(?:fo\\.o|bar\\?|\\^baz)$", []string{"fo.o", "bar?", "^baz"}},
		{"(?i)foo|bar", nil},
	}
	for _, tt := range tests {
		t.Run(tt.re, func(t *testing.T) {
			if got := FindEqualMatches(tt.re); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindEqualMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}

var result []string

func Benchmark_FindMatches(b *testing.B) {
	for _, test := range []struct {
		re string
	}{
		{"foo|bar"},
		{"foo|foo|bar"},
		{"foo|foobar"},
		{"(foo|(ba|ar))"},
		{"^(?:fo\\.o|bar\\?|\\^baz)$"},
		{"(?i)foo|bar"},
	} {
		b.Run(fmt.Sprintf("prom_%s", test.re), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result = append(result, FindSetMatches(test.re)...)
			}
		})
		b.Run(fmt.Sprintf("loki_%s", test.re), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result = append(result, FindEqualMatches(test.re)...)
			}
		})
	}
}

// Bitmap used by func isRegexMetaCharacter to check whether a character needs to be escaped.
var regexMetaCharacterBytes [16]byte

// isRegexMetaCharacter reports whether byte b needs to be escaped.
func isRegexMetaCharacter(b byte) bool {
	return b < utf8.RuneSelf && regexMetaCharacterBytes[b%16]&(1<<(b/16)) != 0
}
func init() {
	for _, b := range []byte(`.+*?()|[]{}^$`) {
		regexMetaCharacterBytes[b%16] |= 1 << (b / 16)
	}
}

// FindSetMatches returns list of values that can be equality matched on.
// copied from Prometheus querier.go, removed check for Prometheus wrapper.
// Returns list of values that can regex matches.
func FindSetMatches(pattern string) []string {
	escaped := false
	sets := []*strings.Builder{{}}
	for i := 0; i < len(pattern); i++ {
		if escaped {
			switch {
			case isRegexMetaCharacter(pattern[i]):
				sets[len(sets)-1].WriteByte(pattern[i])
			case pattern[i] == '\\':
				sets[len(sets)-1].WriteByte('\\')
			default:
				return nil
			}
			escaped = false
		} else {
			switch {
			case isRegexMetaCharacter(pattern[i]):
				if pattern[i] == '|' {
					sets = append(sets, &strings.Builder{})
				} else {
					return nil
				}
			case pattern[i] == '\\':
				escaped = true
			default:
				sets[len(sets)-1].WriteByte(pattern[i])
			}
		}
	}
	matches := make([]string, 0, len(sets))
	for _, s := range sets {
		if s.Len() > 0 {
			matches = append(matches, s.String())
		}
	}
	return matches
}
