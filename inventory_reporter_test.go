package dsco_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/byte4ever/dsco"
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
