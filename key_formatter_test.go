package dsco

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnvKeyFormatter verifies env-layer key formatting:
// uppercase, dashes between segments, prefix-prepended.
func TestEnvKeyFormatter(t *testing.T) {
	t.Parallel()
	f := newEnvKeyFormatter("MYAPP")

	assert.Equal(t, "env", f.LayerKind())
	assert.Equal(t, "env:MYAPP", f.LayerName())
	assert.Equal(t, "MYAPP-DATABASE-HOST", f.FormatKey("database-host"))
	assert.Equal(t, "MYAPP-MAX_RETRY", f.FormatKey("max_retry"))
}

// TestCmdlineKeyFormatter verifies cmdline-layer key formatting:
// dashes between segments, --name= prefix.
func TestCmdlineKeyFormatter(t *testing.T) {
	t.Parallel()
	f := newCmdlineKeyFormatter()

	assert.Equal(t, "cmdline", f.LayerKind())
	assert.Equal(t, "cmdline", f.LayerName())
	assert.Equal(t, "--database-host=", f.FormatKey("database-host"))
}

// TestFileKeyFormatter verifies file-layer key formatting:
// raw alias path, dot-separated for human readability.
func TestFileKeyFormatter(t *testing.T) {
	t.Parallel()
	f := newFileKeyFormatter("config.yaml")

	assert.Equal(t, "file", f.LayerKind())
	assert.Equal(t, "file:config.yaml", f.LayerName())
	assert.Equal(t, "database.host", f.FormatKey("database-host"))
}

// TestNilKeyFormatter verifies the no-op formatter returned when a layer
// cannot enumerate keys (custom string providers).
func TestNilKeyFormatter(t *testing.T) {
	t.Parallel()
	f := newNilKeyFormatter("my-provider")

	assert.Equal(t, "", f.LayerKind())
	assert.Equal(t, "my-provider", f.LayerName())
	assert.Equal(t, "", f.FormatKey("anything"))
}
