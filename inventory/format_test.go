package inventory_test

import (
	"bytes"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/inventory"
)

var update = flag.Bool("update", false, "update golden files")

// fixtureReport returns a deterministic Report covering the three
// interesting cases (key only, satisfied only, both).
func fixtureReport() *inventory.Report {
	return &inventory.Report{
		Type: "github.com/example/myapp.Config",
		Fields: []inventory.Field{
			{
				Path:   "Database.Host",
				GoType: "*string",
				Key: &inventory.KeySpec{
					Layer: "env", Key: "MYAPP-DATABASE-HOST",
				},
			},
			{
				Path:   "Database.Port",
				GoType: "*int",
				Satisfied: &inventory.Satisfaction{
					LayerID: "defaults", Value: 5432,
				},
				Key: &inventory.KeySpec{
					Layer: "cmdline", Key: "--database-port=",
				},
			},
			{
				Path:   "Server.Timeout",
				GoType: "*time.Duration",
				Satisfied: &inventory.Satisfaction{
					LayerID: "defaults", Value: 30 * time.Second,
				},
			},
		},
	}
}

// TestWriteJSONMatchesGolden verifies JSON output is byte-stable.
func TestWriteJSONMatchesGolden(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	require.NoError(t, fixtureReport().WriteJSON(&buf))

	checkGolden(t, "testdata/sample.json", buf.Bytes())
}

// TestWriteYAMLMatchesGolden verifies YAML output is byte-stable.
func TestWriteYAMLMatchesGolden(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	require.NoError(t, fixtureReport().WriteYAML(&buf))

	checkGolden(t, "testdata/sample.yaml", buf.Bytes())
}

// checkGolden compares got to the contents of path; with -update, writes
// got to path instead.
func checkGolden(t *testing.T, path string, got []byte) {
	t.Helper()
	if *update {
		require.NoError(t, os.WriteFile(path, got, 0o644))
		return
	}
	want, err := os.ReadFile(path)
	require.NoError(t, err, "missing golden — run with -update to generate")
	assert.Equal(t, string(want), string(got))
}

type errWriter struct{ err error }

func (w errWriter) Write([]byte) (int, error) { return 0, w.err }

// nthErrWriter succeeds for the first n-1 Write calls and fails on the
// nth call with the configured error.
type nthErrWriter struct {
	err   error
	n     int
	calls int
}

func (w *nthErrWriter) Write(p []byte) (int, error) {
	w.calls++
	if w.calls >= w.n {
		return 0, w.err
	}

	return len(p), nil
}

// TestWriteJSONPropagatesWriterError covers the error branch when the
// underlying io.Writer fails.
func TestWriteJSONPropagatesWriterError(t *testing.T) {
	t.Parallel()
	err := fixtureReport().WriteJSON(errWriter{err: assert.AnError})
	require.ErrorIs(t, err, assert.AnError)
}

// TestWriteJSONPropagatesSecondWriteError covers the error branch when
// the second writer.Write call (the trailing newline) fails.
func TestWriteJSONPropagatesSecondWriteError(t *testing.T) {
	t.Parallel()

	w := &nthErrWriter{err: assert.AnError, n: 2}
	err := fixtureReport().WriteJSON(w)
	require.ErrorIs(t, err, assert.AnError)
}

// TestWriteYAMLPropagatesWriterError covers the same for WriteYAML.
func TestWriteYAMLPropagatesWriterError(t *testing.T) {
	t.Parallel()
	err := fixtureReport().WriteYAML(errWriter{err: assert.AnError})
	require.ErrorIs(t, err, assert.AnError)
}
