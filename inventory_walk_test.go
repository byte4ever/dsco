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
