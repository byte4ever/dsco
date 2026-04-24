package dsco

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newStrictLayer(t *testing.T) {
	t.Parallel()

	k := newStrictLayer(nil)
	require.True(t, k.isStrict())
}

func Test_newNormalLayer(t *testing.T) {
	t.Parallel()

	k := newNormalLayer(nil)
	require.False(t, k.isStrict())
}

// TestGetFieldValuesGetterStrictLayer verifies that getFieldValuesGetter
// returns the embedded FieldValuesGetter for a strictLayer.
func TestGetFieldValuesGetterStrictLayer(t *testing.T) {
	t.Parallel()

	fvg := NewMockFieldValuesGetter(t)
	sl := newStrictLayer(fvg)
	assert.Equal(t, fvg, sl.getFieldValuesGetter())
}

// TestGetFieldValuesGetterNormalLayer verifies that getFieldValuesGetter
// returns the embedded FieldValuesGetter for a normalLayer.
func TestGetFieldValuesGetterNormalLayer(t *testing.T) {
	t.Parallel()

	fvg := NewMockFieldValuesGetter(t)
	nl := newNormalLayer(fvg)
	assert.Equal(t, fvg, nl.getFieldValuesGetter())
}
