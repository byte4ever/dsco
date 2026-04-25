package dsco_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/svalue"
)

// TestLayerInventoryZeroValueIsUsable verifies that a zero-valued
// LayerInventory is meaningful (empty Provides, empty Note/Name) so
// callers can build it incrementally.
func TestLayerInventoryZeroValueIsUsable(t *testing.T) {
	t.Parallel()

	var inv dsco.LayerInventory
	assert.Empty(t, inv.Name)
	assert.Empty(t, inv.Note)
	assert.Empty(t, inv.Provides)
}

// TestFieldProvisionFields verifies the public field set of FieldProvision.
func TestFieldProvisionFields(t *testing.T) {
	t.Parallel()

	p := dsco.FieldProvision{
		FieldUID: "Database.Host",
		Key:      "MYAPP-DATABASE-HOST",
		Value:    nil,
	}
	assert.Equal(t, "Database.Host", p.FieldUID)
	assert.Equal(t, "MYAPP-DATABASE-HOST", p.Key)
	assert.Nil(t, p.Value)
}

// stubProvider is a minimal NamedStringValuesProvider for tests.
type stubProvider struct {
	name string
	vals svalue.Values
}

func (p *stubProvider) GetName() string                { return p.name }
func (p *stubProvider) GetStringValues() svalue.Values { return p.vals }

// TestStringBasedBuilderReportInventoryEnvKind builds a StringBasedBuilder
// with an env-style KeyFormatter (via the test seam in sbased.go) and
// verifies it reports the right canonical keys without performing any I/O.
//
// The test stays at unit granularity — it does not call Compute. End-to-end
// wiring is covered by Task 11.
func TestStringBasedBuilderReportInventoryEnvKind(t *testing.T) {
	t.Parallel()

	type sub struct {
		Host *string `yaml:"host"`
	}
	type cfg struct {
		Database *sub `yaml:"database"`
		Port     *int `yaml:"port"`
	}

	mdl, err := dsco.BuildModel(&cfg{})
	require.NoError(t, err)

	b, err := dsco.NewStringBasedBuilderForTest(
		&stubProvider{name: "stub", vals: svalue.Values{}},
		"env", "MYAPP",
	)
	require.NoError(t, err)

	inv, err := b.ReportInventory(mdl)
	require.NoError(t, err)
	assert.Equal(t, "env:MYAPP", inv.Name)

	keys := make(map[string]string)
	for _, p := range inv.Provides {
		keys[p.FieldUID] = p.Key
	}
	assert.Equal(t, "MYAPP-DATABASE-HOST", keys["Database.Host"])
	assert.Equal(t, "MYAPP-PORT", keys["Port"])
}

// TestStructBuilderReportInventoryUnit calls StructBuilder.ReportInventory
// directly with a small model and verifies only set fields appear.
// End-to-end Compute coverage lives in Task 11.
func TestStructBuilderReportInventoryUnit(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Host *string `yaml:"host"`
		Port *int    `yaml:"port"`
	}

	defaults := &cfg{Port: dsco.R(5432)}

	sb, err := dsco.NewStructBuilder(defaults, "defaults")
	require.NoError(t, err)

	mdl, err := dsco.BuildModel(&cfg{})
	require.NoError(t, err)

	inv, err := sb.ReportInventory(mdl)
	require.NoError(t, err)

	assert.Equal(t, "struct:defaults", inv.Name)

	uids := make(map[string]any)
	for _, p := range inv.Provides {
		uids[p.FieldUID] = p.Value
	}

	assert.Contains(t, uids, "Port")
	assert.Equal(t, 5432, uids["Port"])
	assert.NotContains(t, uids, "Host",
		"nil-pointer fields must not appear")
}

// TestCollectAliasesIncludesNestedFields verifies that collectAliases
// returns one entry per leaf, with dot→dash path conversion applied.
func TestCollectAliasesIncludesNestedFields(t *testing.T) {
	t.Parallel()

	type sub struct {
		Host *string `yaml:"host"`
	}
	type cfg struct {
		Database *sub `yaml:"database"`
		Port     *int `yaml:"port"`
	}

	mdl, err := dsco.BuildModel(&cfg{})
	require.NoError(t, err)

	aliases, err := dsco.CollectAliasesForTest(mdl)
	require.NoError(t, err)

	values := make(map[string]bool)
	for _, alias := range aliases {
		values[alias] = true
	}

	assert.True(t, values["database-host"], "expected database-host alias")
	assert.True(t, values["port"], "expected port alias")
}
