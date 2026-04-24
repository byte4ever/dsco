package dsco

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMinimalFieldValuesGetterReturnsNil verifies that the test-only stub
// returns (nil, nil) as documented. This covers the GetFieldValuesFrom
// body of minimalFieldValuesGetter.
func TestMinimalFieldValuesGetterReturnsNil(t *testing.T) {
	t.Parallel()

	var g minimalFieldValuesGetter

	vals, err := g.GetFieldValuesFrom(nil)
	assert.Nil(t, vals)
	require.NoError(t, err)
}

// TestInventoryReporterFVGReturnsNil verifies that the test-only adapter
// returns (nil, nil) as documented, covering the GetFieldValuesFrom body
// of inventoryReporterFVG.
func TestInventoryReporterFVGReturnsNil(t *testing.T) {
	t.Parallel()

	fvg := &inventoryReporterFVG{
		InventoryReporter: NewMockInventoryReporter(t),
	}

	vals, err := fvg.GetFieldValuesFrom(nil)
	assert.Nil(t, vals)
	require.NoError(t, err)
}
