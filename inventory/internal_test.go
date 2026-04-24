package inventory

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
)

// TestNormalizeValueStringer verifies that any fmt.Stringer is converted
// to its String() form for serialization.
func TestNormalizeValueStringer(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   any
		want any
	}{
		{"duration", 30 * time.Second, "30s"},
		{
			"time",
			time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC),
			"2026-04-24 12:00:00 +0000 UTC",
		},
		{
			"url",
			mustParseURL("https://example.com/path"),
			"https://example.com/path",
		},
	}

	for _, cc := range cases {
		cc := cc
		t.Run(cc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, cc.want, normalizeValue(cc.in))
		})
	}
}

// TestNormalizeValuePrimitivesPassThrough verifies primitives are
// returned unchanged.
func TestNormalizeValuePrimitivesPassThrough(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 42, normalizeValue(42))
	assert.Equal(t, "hello", normalizeValue("hello"))
	assert.Equal(t, true, normalizeValue(true))
	assert.InEpsilon(t, 3.14, normalizeValue(3.14), 0.0001)
	assert.Nil(t, normalizeValue(nil))
}

// mustParseURL parses a URL string and panics if it fails. Used in
// test setup to keep table entries concise.
func mustParseURL(str string) *url.URL {
	parsed, err := url.Parse(str)
	if err != nil {
		panic(err)
	}

	return parsed
}

// stubFailingReporter is a minimal InventoryReporter whose ReportInventory
// always returns an error. Used in-package to exercise the invErr branch
// of computeFromWalk without going through the public layer API.
type stubFailingReporter struct{ err error }

func (r *stubFailingReporter) ReportInventory(
	_ dsco.ModelInterface,
) (dsco.LayerInventory, error) {
	return dsco.LayerInventory{}, r.err
}

// TestComputeFromWalkPropagatesReporterError verifies that computeFromWalk
// propagates errors returned by a reporter's ReportInventory and wraps them
// with dsco.ErrFiller.
func TestComputeFromWalkPropagatesReporterError(t *testing.T) {
	t.Parallel()

	type cfg struct{ Host *string `yaml:"host"` }

	mdl, err := dsco.BuildModel(&cfg{})
	require.NoError(t, err)

	walk := &dsco.InventoryWalk{
		Model: mdl,
		Reporters: []dsco.InventoryReporter{
			&stubFailingReporter{err: assert.AnError},
		},
	}

	_, err = computeFromWalk(walk)
	require.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
	assert.ErrorIs(t, err, dsco.ErrFiller)
}

// TestTrimStructPrefix verifies that trimStructPrefix removes the
// "struct:" prefix when present, and returns the input unchanged when
// the prefix is absent or the string is too short.
func TestTrimStructPrefix(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   string
		want string
	}{
		{"with prefix", "struct:defaults", "defaults"},
		{"no prefix — cmdline", "cmdline", "cmdline"},
		{
			"no prefix — env",
			"env:MYAPP",
			"env:MYAPP",
		},
		{"exact prefix length — no suffix", "struct:", "struct:"},
	}

	for _, cc := range cases {
		cc := cc
		t.Run(cc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, cc.want, trimStructPrefix(cc.in))
		})
	}
}
