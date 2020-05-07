package logql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newRegexpExtractor(t *testing.T) {

	e, err := newRegexpExtractor("bar")
	require.Error(t, err)

	e, err = newRegexpExtractor("(?P<foo>.*) bar (?P<foo>.*)")
	require.Nil(t, err)

	require.Equal(t, []string{"foo", "foo"}, e.Fields())

	require.Equal(t, []string{"foo", "foo"}, e.Values("foo bar foo"))
	require.Equal(t, 0, len(e.Values("foo bar")))

}
