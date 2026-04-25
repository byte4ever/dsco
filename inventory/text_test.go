package inventory_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/inventory"
)

// TestWriteTextMatchesGolden verifies the human-readable layout is stable.
func TestWriteTextMatchesGolden(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	require.NoError(t, fixtureReport().WriteText(&buf))

	checkGolden(t, "testdata/sample.txt", buf.Bytes())
}

// TestWriteTextEmDashForEmpty verifies missing key/default cells print "—".
func TestWriteTextEmDashForEmpty(t *testing.T) {
	t.Parallel()

	rep := &inventory.Report{
		Type: "Cfg",
		Fields: []inventory.Field{
			{Path: "X", GoType: "*string"},
		},
	}
	var buf bytes.Buffer
	require.NoError(t, rep.WriteText(&buf))

	out := buf.String()
	assert.Contains(t, out, "—", "empty cells must use em-dash")
}

// TestWriteTextTruncatesLongDefaults verifies values >40 chars get
// truncated with ellipsis.
func TestWriteTextTruncatesLongDefaults(t *testing.T) {
	t.Parallel()

	long := strings.Repeat("x", 80)
	rep := &inventory.Report{
		Type: "Cfg",
		Fields: []inventory.Field{
			{
				Path: "X", GoType: "*string",
				Satisfied: &inventory.Satisfaction{LayerID: "d", Value: long},
			},
		},
	}
	var buf bytes.Buffer
	require.NoError(t, rep.WriteText(&buf))

	assert.Contains(t, buf.String(), "…", "long values must end with ellipsis")
	assert.NotContains(t, buf.String(), long, "full value must not appear")
}

// TestWriteTextPropagatesWriterError covers the error branch.
func TestWriteTextPropagatesWriterError(t *testing.T) {
	t.Parallel()
	err := fixtureReport().WriteText(errWriter{err: assert.AnError})
	require.ErrorIs(t, err, assert.AnError)
}
