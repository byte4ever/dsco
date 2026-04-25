package dsco

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/internal/fvalue"
)

// bareFieldValuesGetter is an unexported FieldValuesGetter that does NOT
// implement InventoryReporter. Used to exercise the
// ErrLayerNotInventoryReporter branch in prepareInventoryWalkFromPolicies.
type bareFieldValuesGetter struct{}

func (bareFieldValuesGetter) GetFieldValuesFrom(
	_ ModelInterface,
) (fvalue.Values, error) {
	return nil, nil //nolint:nilnil // test stub
}

// TestPrepareInventoryWalkFromPoliciesRejectsNonReporter verifies that
// prepareInventoryWalkFromPolicies returns ErrLayerNotInventoryReporter when
// a policy's FieldValuesGetter does not implement InventoryReporter.
func TestPrepareInventoryWalkFromPoliciesRejectsNonReporter(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Host *string `yaml:"host"`
	}

	mdl, err := buildModel(&cfg{})
	require.NoError(t, err)

	policies := constraintLayerPolicies{
		newNormalLayer(bareFieldValuesGetter{}),
	}

	_, err = prepareInventoryWalkFromPolicies(policies, mdl)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrLayerNotInventoryReporter)
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
