package inventory_test

import (
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/inventory"
)

// TestComputeCanonicalKeyFirstLayerWins verifies that when env and
// cmdline both can supply the same field, the env key (first in the
// layer list) wins — matching dsco.Fill's first-layer-wins semantics.
//
//nolint:paralleltest // modifies os.Args
func TestComputeCanonicalKeyFirstLayerWins(t *testing.T) {
	// cmdline layer reads os.Args at registration time; supply clean args.
	os.Args = []string{"testapp"}

	type cfg struct {
		Host *string `yaml:"host"`
	}
	var c *cfg

	report, err := inventory.Compute(
		&c,
		dsco.WithEnvLayer("MYAPP"),
		dsco.WithCmdlineLayer(),
	)
	require.NoError(t, err)
	require.Len(t, report.Fields, 1)

	require.NotNil(t, report.Fields[0].Key)
	assert.Equal(t, "env", report.Fields[0].Key.Layer)
	assert.Equal(t, "MYAPP-HOST", report.Fields[0].Key.Key)
}

// TestComputeSatisfiedByDefaults verifies that struct-layer values
// appear in Field.Satisfied while string-layer keys remain in Field.Key.
func TestComputeSatisfiedByDefaults(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Port *int `yaml:"port"`
	}
	defaults := &cfg{Port: dsco.R(5432)}
	var c *cfg

	report, err := inventory.Compute(
		&c,
		dsco.WithStructLayer(defaults, "defaults"),
		dsco.WithEnvLayer("MYAPP"),
	)
	require.NoError(t, err)
	require.Len(t, report.Fields, 1)

	require.NotNil(t, report.Fields[0].Satisfied)
	assert.Equal(t, "defaults", report.Fields[0].Satisfied.LayerID)
	assert.Equal(t, 5432, report.Fields[0].Satisfied.Value)

	require.NotNil(t, report.Fields[0].Key)
	assert.Equal(t, "env", report.Fields[0].Key.Layer)
	assert.Equal(t, "MYAPP-PORT", report.Fields[0].Key.Key)
}

// TestComputeSortsFieldsByPath verifies that the report's Fields slice
// is sorted lexicographically by Path.
func TestComputeSortsFieldsByPath(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Zeta  *string `yaml:"zeta"`
		Alpha *string `yaml:"alpha"`
	}
	var c *cfg

	report, err := inventory.Compute(&c, dsco.WithEnvLayer("MYAPP"))
	require.NoError(t, err)

	paths := make([]string, len(report.Fields))
	for i, f := range report.Fields {
		paths[i] = f.Path
	}

	assert.True(t, sort.StringsAreSorted(paths), "fields must be sorted by path")
}

// TestComputeRejectsNonPointerCfg verifies the error path when cfg is not a
// pointer.
func TestComputeRejectsNonPointerCfg(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Host *string `yaml:"host"`
	}
	_, err := inventory.Compute(cfg{}, dsco.WithEnvLayer("MYAPP"))
	require.Error(t, err)
}

// TestComputePropagatesPrepareWalkError verifies that Compute propagates an
// error returned by PrepareInventoryWalk (e.g. duplicate cmdline layers).
func TestComputePropagatesPrepareWalkError(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Host *string `yaml:"host"`
	}
	var c *cfg

	// Two cmdline layers trigger a dedup error inside PrepareInventoryWalk.
	_, err := inventory.Compute(
		&c,
		dsco.WithCmdlineLayer(),
		dsco.WithCmdlineLayer(),
	)
	require.Error(t, err)
}

// TestComputeErrorsSatisfyErrFiller verifies that every error returned by
// Compute satisfies errors.Is(err, dsco.ErrFiller), per spec.
func TestComputeErrorsSatisfyErrFiller(t *testing.T) {
	t.Parallel()

	type cfg struct{ Host *string `yaml:"host"` }

	_, err := inventory.Compute(cfg{}, dsco.WithEnvLayer("MYAPP"))
	require.Error(t, err)
	assert.ErrorIs(t, err, dsco.ErrFiller)
}
