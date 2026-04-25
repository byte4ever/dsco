package dsco_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
)

// TestPrepareInventoryWalkBuildsModelAndReporters verifies that
// PrepareInventoryWalk yields a non-nil model and one InventoryReporter
// per layer without performing any I/O.
func TestPrepareInventoryWalkBuildsModelAndReporters(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Host *string `yaml:"host"`
	}

	walk, err := dsco.PrepareInventoryWalk(
		&cfg{},
		dsco.WithEnvLayer("MYAPP"),
		dsco.WithStructLayer(&cfg{Host: dsco.R("localhost")}, "defaults"),
	)
	require.NoError(t, err)

	require.NotNil(t, walk)
	require.NotNil(t, walk.Model)
	assert.Len(t, walk.Reporters, 2)
}

// TestBuildModelRejectsNonPointerCfg verifies error path for non-pointer
// cfg values (mirrors Fill's contract).
func TestBuildModelRejectsNonPointerCfg(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Host *string `yaml:"host"`
	}

	_, err := dsco.BuildModel(cfg{})
	require.Error(t, err)
}

// TestPrepareInventoryWalkRejectsNonPointerCfg verifies that
// PrepareInventoryWalk propagates the buildModel error when cfg is not a
// pointer.
func TestPrepareInventoryWalkRejectsNonPointerCfg(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Host *string `yaml:"host"`
	}

	_, err := dsco.PrepareInventoryWalk(cfg{})
	require.Error(t, err)
}

// TestPrepareInventoryWalkPropagatesGetPoliciesError verifies that
// PrepareInventoryWalk propagates errors returned by GetPolicies.  A
// duplicate cmdline layer triggers a dedup error in builders.go.
func TestPrepareInventoryWalkPropagatesGetPoliciesError(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Host *string `yaml:"host"`
	}

	// Two cmdline layers trigger the dedup error in builders.go.
	_, err := dsco.PrepareInventoryWalk(
		&cfg{},
		dsco.WithCmdlineLayer(),
		dsco.WithCmdlineLayer(),
	)
	require.Error(t, err)
}

// TestPrepareInventoryWalkRejectsNonReporterLayer verifies that
// PrepareInventoryWalk returns ErrLayerNotInventoryReporter when a layer's
// FieldValuesGetter does not implement InventoryReporter.
func TestPrepareInventoryWalkRejectsNonReporterLayer(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Host *string `yaml:"host"`
	}

	_, err := dsco.PrepareInventoryWalk(
		&cfg{},
		dsco.WithNonReporterLayerForTest(),
	)
	require.Error(t, err)
	require.ErrorIs(t, err, dsco.ErrLayerNotInventoryReporter)
}

// noopInventoryReporter is a minimal InventoryReporter for use in tests.
// ReportInventory returns an empty LayerInventory without error.
type noopInventoryReporter struct{}

func (noopInventoryReporter) ReportInventory(
	_ dsco.ModelInterface,
) (dsco.LayerInventory, error) {
	return dsco.LayerInventory{}, nil
}

// TestWithInventoryReporterLayerForTestIntegration verifies that
// WithInventoryReporterLayerForTest produces a Layer whose reporter is
// surfaced by PrepareInventoryWalk.
func TestWithInventoryReporterLayerForTestIntegration(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Host *string `yaml:"host"`
	}

	walk, err := dsco.PrepareInventoryWalk(
		&cfg{},
		dsco.WithInventoryReporterLayerForTest(noopInventoryReporter{}),
	)
	require.NoError(t, err)
	assert.Len(t, walk.Reporters, 1)
}
