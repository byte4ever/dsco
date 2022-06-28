package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_stackEmbed_push(t *testing.T) {
	t.Parallel()

	t.Run(
		"initial state",
		func(t *testing.T) {
			t.Parallel()

			v := &elemEmbedded{}

			var s stackEmbed

			s.push(v)
			require.True(t, s.more())
			require.Len(t, s, 1)
			require.Equal(t, v, s[0])
		},
	)

	t.Run(
		"initial state",
		func(t *testing.T) {
			t.Parallel()

			v0 := &elemEmbedded{}
			v := &elemEmbedded{}

			s := stackEmbed{v0}

			s.push(v)
			require.True(t, s.more())
			require.Len(t, s, 2)
			require.Equal(t, v, s[1])
		},
	)
}

func Test_stackEmbed_pop(t *testing.T) {
	t.Parallel()

	t.Run(
		"initial state",
		func(t *testing.T) {
			t.Parallel()

			v := &elemEmbedded{}

			s := stackEmbed{v}

			got := s.pop()
			require.False(t, s.more())
			require.Len(t, s, 0)
			require.Equal(t, v, got)
		},
	)

	t.Run(
		"initial state",
		func(t *testing.T) {
			t.Parallel()

			v0 := &elemEmbedded{}
			v := &elemEmbedded{}

			s := stackEmbed{v0, v}

			got := s.pop()
			require.True(t, s.more())
			require.Len(t, s, 1)
			require.Equal(t, v, got)
		},
	)
}
