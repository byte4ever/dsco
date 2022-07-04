package dsco

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newStrictLayer(t *testing.T) {
	k := newStrictLayer(nil)
	require.True(t, k.isStrict())
}

func Test_newNormalLayer(t *testing.T) {
	k := newNormalLayer(nil)
	require.False(t, k.isStrict())
}
