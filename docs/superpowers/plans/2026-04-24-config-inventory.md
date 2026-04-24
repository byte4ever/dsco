# Config Inventory Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a static, no-I/O inventory of configuration keys a `dsco.Fill(&cfg, layers...)` call would expect, exposed via a new `inventory/` sub-package with JSON / YAML / plain-text serialization.

**Architecture:** New `InventoryReporter` interface implemented by every layer builder; a `KeyFormatter` injected into `StringBasedBuilder` so env / cmdline / file each format their own canonical key. `inventory.Compute` calls a small dsco-root helper (`PrepareInventoryWalk`) that runs only the model-build + builder-construction phases (no value-loading) and reduces the per-layer reports into one canonical `*Report` per leaf field.

**Tech Stack:** Go 1.21+, `github.com/goccy/go-json` (per project import policy), `github.com/goccy/go-yaml` (per project import policy), `github.com/stretchr/testify`, `github.com/vektra/mockery/v2` (in-package, expecter pattern).

**Spec:** `docs/superpowers/specs/2026-04-24-config-inventory-design.md`

---

## File Structure

**New files:**

- `inventory/doc.go` — package overview
- `inventory/inventory.go` — `Compute` + `Report` / `Field` / `Satisfaction` / `KeySpec` types
- `inventory/normalize.go` — `normalizeValue` helper + custom `Satisfaction` marshalers
- `inventory/format.go` — `WriteJSON`, `WriteYAML`
- `inventory/text.go` — `WriteText` (column layout)
- `inventory/inventory_test.go`, `normalize_test.go`, `format_test.go`, `text_test.go`
- `inventory/testdata/sample.json`, `sample.yaml`, `sample.txt` — golden files
- `inventory_reporter.go` (dsco root) — `InventoryReporter` interface, `LayerInventory`, `FieldProvision`
- `inventory_reporter_test.go` (dsco root) — tests for `ReportInventory` on each builder
- `key_formatter.go` (dsco root) — `KeyFormatter` interface + env / cmdline / file / nil implementations
- `key_formatter_test.go` (dsco root)
- `inventory_walk.go` (dsco root) — `PrepareInventoryWalk` helper consumed by `inventory/`
- `inventory_walk_test.go` (dsco root)
- `mock_InventoryReporter_test.go` (dsco root, mockery-generated)

**Modified files:**

- `sbased.go` — accept `KeyFormatter` in constructor; add `aliases()` accessor; add `ReportInventory`
- `builders.go` — pass per-layer `KeyFormatter` when constructing `StringBasedBuilder`
- `structs.go` — add `ReportInventory`
- `filler.go` — extract `generateModel` and `generateBuilders` into helpers shared with `PrepareInventoryWalk`
- `.mockery.yaml` — register `MockInventoryReporter`
- `README.md` — add Inventory section
- `QUICKSTART.md` — add a small example
- `doc.go` — extend package overview with the inventory feature

---

## Task 1: Branch verification

**Files:** none

- [ ] **Step 1: Confirm working branch**

Run: `git status`
Expected: `On branch feature/config-inventory`. If detached or on master, run `git switch feature/config-inventory`.

- [ ] **Step 2: Confirm spec is committed**

Run: `git log --oneline -5`
Expected: see commit `docs: add design spec for config key inventory feature`.

---

## Task 2: Create inventory package skeleton

**Files:**
- Create: `inventory/doc.go`
- Create: `inventory/inventory.go`

- [ ] **Step 1: Create `inventory/doc.go`**

```go
// Package inventory computes the static list of configuration keys a
// dsco-managed struct expects, given a set of layers, without performing
// any I/O.
//
// See https://github.com/byte4ever/dsco for the parent project.
package inventory
```

- [ ] **Step 2: Create `inventory/inventory.go` with type definitions only (no Compute yet)**

```go
package inventory

// Report is the static inventory of a config struct against a layer set.
type Report struct {
	Type   string  `json:"type"   yaml:"type"`
	Fields []Field `json:"fields" yaml:"fields"`
}

// Field describes the canonical key (and any baked-in default) for one
// leaf field of the config struct.
type Field struct {
	Path      string        `json:"path"                yaml:"path"`
	GoType    string        `json:"go_type"             yaml:"go_type"`
	Satisfied *Satisfaction `json:"satisfied,omitempty" yaml:"satisfied,omitempty"`
	Key       *KeySpec      `json:"key,omitempty"       yaml:"key,omitempty"`
}

// Satisfaction records that a struct layer already provides a value
// for this field.
type Satisfaction struct {
	LayerID string `json:"layer_id" yaml:"layer_id"`
	Value   any    `json:"value"    yaml:"value"`
}

// KeySpec is the canonical (highest-precedence) key form a string-based
// layer would accept to supply this field.
type KeySpec struct {
	Layer string `json:"layer" yaml:"layer"`
	Key   string `json:"key"   yaml:"key"`
}
```

- [ ] **Step 3: Run `go build ./inventory/...` to confirm package compiles**

Run: `go build ./inventory/...`
Expected: no output (success).

- [ ] **Step 4: Commit**

```bash
git add inventory/doc.go inventory/inventory.go
git commit -m "inventory: scaffold package and core data types"
```

---

## Task 3: Define InventoryReporter interface and LayerInventory types

**Files:**
- Create: `inventory_reporter.go`
- Test: `inventory_reporter_test.go`

- [ ] **Step 1: Write the failing test**

Create `inventory_reporter_test.go`:

```go
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
```

- [ ] **Step 2: Run test to verify it fails to compile**

Run: `go test ./...`
Expected: build error — `LayerInventory` / `FieldProvision` undefined.

- [ ] **Step 3: Create `inventory_reporter.go` in dsco root**

```go
package dsco

// InventoryReporter is implemented by layer builders that can describe
// what they would contribute to an inventory without performing I/O.
//
// Pattern: Strategy — each layer reports its own contribution shape (key
// form for string-based layers, baked-in values for struct layers).
type InventoryReporter interface {
	// ReportInventory returns the layer's contribution to an inventory walk
	// for the given model. Implementations must not perform any I/O.
	ReportInventory(model ModelInterface) (LayerInventory, error)
}

// LayerInventory is one layer's contribution to a Report.
type LayerInventory struct {
	// Name uniquely identifies the layer instance, e.g. "env:MYAPP",
	// "cmdline", "file:<id>", "struct:<id>", or a custom provider name.
	Name string

	// Note carries optional information for callers that cannot enumerate
	// keys (typically custom string providers).
	Note string

	// Provides lists every (field, key|value) pair this layer can supply
	// to the model.
	Provides []FieldProvision
}

// FieldProvision is one (field, layer) pair from a layer's perspective.
type FieldProvision struct {
	// FieldUID matches the model's field UID.
	FieldUID string

	// Key is the canonical key for string-based layers; empty for
	// struct layers.
	Key string

	// Value is the baked-in value for struct layers; nil for string-based
	// layers.
	Value any
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./...`
Expected: all tests pass.

- [ ] **Step 5: Commit**

```bash
git add inventory_reporter.go inventory_reporter_test.go
git commit -m "dsco: add InventoryReporter interface and layer inventory types"
```

---

## Task 4: Define KeyFormatter abstraction

**Files:**
- Create: `key_formatter.go`
- Test: `key_formatter_test.go`

- [ ] **Step 1: Write the failing tests**

Create `key_formatter_test.go`:

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./...`
Expected: build error — formatters undefined.

- [ ] **Step 3: Create `key_formatter.go`**

```go
package dsco

import "strings"

// KeyFormatter renders a layer-internal alias path into the canonical
// user-facing key form for a specific layer kind (env, cmdline, file).
//
// Pattern: Strategy — each layer kind formats keys differently; the
// formatter is injected into StringBasedBuilder at construction time.
type KeyFormatter interface {
	// LayerKind returns the layer category (e.g. "env", "cmdline", "file").
	// Empty for layers that cannot enumerate keys.
	LayerKind() string

	// LayerName returns the layer instance identifier
	// (e.g. "env:MYAPP", "cmdline", "file:config.yaml").
	LayerName() string

	// FormatKey converts an internal alias path (dash-separated, lowercase)
	// to the canonical user-facing key for this layer kind. Returns the
	// empty string when the layer cannot enumerate keys.
	FormatKey(aliasPath string) string
}

// envKeyFormatter formats keys for environment-variable layers:
// PREFIX-UPPER-CASE-DASHED.
type envKeyFormatter struct {
	prefix string
}

func newEnvKeyFormatter(prefix string) *envKeyFormatter {
	return &envKeyFormatter{prefix: prefix}
}

func (f *envKeyFormatter) LayerKind() string { return "env" }
func (f *envKeyFormatter) LayerName() string { return "env:" + f.prefix }

func (f *envKeyFormatter) FormatKey(aliasPath string) string {
	return f.prefix + "-" + strings.ToUpper(aliasPath)
}

// cmdlineKeyFormatter formats keys for command-line layers: --name=.
type cmdlineKeyFormatter struct{}

func newCmdlineKeyFormatter() *cmdlineKeyFormatter {
	return &cmdlineKeyFormatter{}
}

func (f *cmdlineKeyFormatter) LayerKind() string { return "cmdline" }
func (f *cmdlineKeyFormatter) LayerName() string { return "cmdline" }

func (f *cmdlineKeyFormatter) FormatKey(aliasPath string) string {
	return "--" + aliasPath + "="
}

// fileKeyFormatter formats keys for file layers: dot-separated YAML path.
type fileKeyFormatter struct {
	id string
}

func newFileKeyFormatter(id string) *fileKeyFormatter {
	return &fileKeyFormatter{id: id}
}

func (f *fileKeyFormatter) LayerKind() string { return "file" }
func (f *fileKeyFormatter) LayerName() string { return "file:" + f.id }

func (f *fileKeyFormatter) FormatKey(aliasPath string) string {
	return strings.ReplaceAll(aliasPath, "-", ".")
}

// nilKeyFormatter is a no-op formatter for layers (custom string
// providers) that cannot enumerate keys statically. LayerKind is empty so
// reduce-pass logic skips them when picking a canonical key.
type nilKeyFormatter struct {
	name string
}

func newNilKeyFormatter(name string) *nilKeyFormatter {
	return &nilKeyFormatter{name: name}
}

func (f *nilKeyFormatter) LayerKind() string             { return "" }
func (f *nilKeyFormatter) LayerName() string             { return f.name }
func (f *nilKeyFormatter) FormatKey(_ string) string     { return "" }
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -run KeyFormatter ./...`
Expected: all four formatter tests pass.

- [ ] **Step 5: Commit**

```bash
git add key_formatter.go key_formatter_test.go
git commit -m "dsco: add KeyFormatter strategy for per-layer key rendering"
```

---

## Task 5: Inject KeyFormatter into StringBasedBuilder

**Files:**
- Modify: `sbased.go`
- Modify: `builders.go`
- Modify: `sbased_test.go`, `builders_test.go` (if any tests break)

- [ ] **Step 1: Add a `keyFormatter` field to `StringBasedBuilder` and a constructor variant**

Modify `sbased.go`. Find the `StringBasedBuilder` struct definition (around line 117–122):

```go
// StringBasedBuilder is a value bases builder depending on text values.
type StringBasedBuilder struct {
	internalOpts
	values         svalue.Values
	expandedValues map[string]*fvalue.Value
}
```

Replace with:

```go
// StringBasedBuilder is a value bases builder depending on text values.
type StringBasedBuilder struct {
	internalOpts
	values         svalue.Values
	expandedValues map[string]*fvalue.Value
	keyFormatter   KeyFormatter
}
```

- [ ] **Step 2: Add a key-formatter constructor variant (keep the old one as a thin wrapper for backward compatibility within the package)**

Add immediately after `NewStringBasedBuilder` in `sbased.go`:

```go
// newStringBasedBuilderWithFormatter is an internal constructor that
// behaves like NewStringBasedBuilder but records the KeyFormatter used to
// render the layer's keys in inventory reports.
func newStringBasedBuilderWithFormatter(
	provider StringValuesProvider,
	formatter KeyFormatter,
	options ...Option,
) (*StringBasedBuilder, error) {
	b, err := NewStringBasedBuilder(provider, options...)
	if err != nil {
		return nil, err
	}

	b.keyFormatter = formatter

	return b, nil
}
```

- [ ] **Step 3: Update `wrapEnvBuild`, `wrapCmdlineBuild`, `wrapStringProviderBuild` in `builders.go` to use the formatter constructor**

In `builders.go`, replace the body of `wrapCmdlineBuild` (lines around 145–164):

```go
func wrapCmdlineBuild(
	to *layerBuilder,
	wrap func(FieldValuesGetter) constraintLayerPolicy,
	options []Option,
) error {
	if idx := to.dedupId("cmdLine"); idx != nil {
		return CmdlineAlreadyUsedError{
			Index: *idx,
		}
	}

	cmdLine, err := cmdline.NewEntriesProvider(os.Args[1:])
	if err != nil {
		return fmt.Errorf("cmdline builder: %w", err)
	}

	builder, err := newStringBasedBuilderWithFormatter(
		cmdLine,
		newCmdlineKeyFormatter(),
		options...,
	)
	if err != nil {
		return err
	}

	to.builders = append(to.builders, wrap(builder))

	return nil
}
```

Replace `wrapEnvBuild` body (lines around 199–225):

```go
func wrapEnvBuild(
	to *layerBuilder,
	wrap func(FieldValuesGetter) constraintLayerPolicy,
	prefix string,
	options []Option,
) error {
	if idx := to.dedupId(fmt.Sprintf("env(%s)", prefix)); idx != nil {
		return DuplicateEnvPrefixError{
			Index:  *idx,
			Prefix: prefix,
		}
	}

	envProvider, err := env.NewEntriesProvider(prefix)
	if err != nil {
		return fmt.Errorf("env builder: %w", err)
	}

	builder, err := newStringBasedBuilderWithFormatter(
		envProvider,
		newEnvKeyFormatter(prefix),
		options...,
	)
	if err != nil {
		return err
	}

	to.addBuilder(wrap(builder))

	return nil
}
```

Replace `wrapStringProviderBuild` body (lines around 362–390):

```go
func wrapStringProviderBuild(
	to *layerBuilder,
	wrap func(bg FieldValuesGetter) constraintLayerPolicy,
	provider NamedStringValuesProvider,
	options []Option,
) error {
	providerName := provider.GetName()

	if idx := to.dedupId(
		fmt.Sprintf("stringProvider(%s)", providerName),
	); idx != nil {
		return DuplicateStringProviderError{
			Index: *idx,
			ID:    providerName,
		}
	}

	builder, err := newStringBasedBuilderWithFormatter(
		provider,
		newNilKeyFormatter(providerName),
		options...,
	)
	if err != nil {
		return err
	}

	to.addBuilder(wrap(builder))

	return nil
}
```

- [ ] **Step 4: Run all tests to verify nothing regressed**

Run: `go test -race ./...`
Expected: all tests pass.

- [ ] **Step 5: Commit**

```bash
git add sbased.go builders.go
git commit -m "dsco: inject KeyFormatter into StringBasedBuilder for env/cmdline/provider layers"
```

---

## Task 6: Implement ReportInventory on StringBasedBuilder

**Files:**
- Modify: `sbased.go`
- Test: `inventory_reporter_test.go` (extend)

- [ ] **Step 1: Write the failing unit test (no Compute dependency)**

Append to `inventory_reporter_test.go`:

```go
import (
	"github.com/byte4ever/dsco/svalue"
)

// stubProvider is a minimal NamedStringValuesProvider for tests.
type stubProvider struct {
	name string
	vals svalue.Values
}

func (p *stubProvider) GetName() string                   { return p.name }
func (p *stubProvider) GetStringValues() svalue.Values    { return p.vals }

// TestStringBasedBuilderReportInventoryEnvKind builds a StringBasedBuilder
// with an env-style KeyFormatter (via the exported NewStringBasedBuilder
// + an internal helper) and verifies it reports the right canonical keys
// without performing any I/O.
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

	// Construct a StringBasedBuilder with an env-style formatter via the
	// internal test seam exposed in inventory_walk.go (Task 8).
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
```

`NewStringBasedBuilderForTest` is a tiny test seam added in step 3 below.

- [ ] **Step 2: Implement `ReportInventory` on `StringBasedBuilder`**

Append to `sbased.go`:

```go
// ReportInventory implements InventoryReporter by walking the model's
// alias map and rendering each entry through the layer's KeyFormatter.
// No I/O is performed.
func (s *StringBasedBuilder) ReportInventory(
	m ModelInterface,
) (LayerInventory, error) {
	const errCtx = "reporting inventory"

	if s.keyFormatter == nil {
		// Defensive: any builder constructed via the layer wrappers in
		// builders.go has a non-nil formatter (see Task 5).
		return LayerInventory{}, fmt.Errorf(
			"%s: nil key formatter", errCtx,
		)
	}

	aliases, err := collectAliases(m)
	if err != nil {
		return LayerInventory{}, fmt.Errorf("%s: %w", errCtx, err)
	}

	provides := make([]FieldProvision, 0, len(aliases))
	for fieldUID, aliasPath := range aliases {
		provides = append(provides, FieldProvision{
			FieldUID: fieldUID,
			Key:      s.keyFormatter.FormatKey(aliasPath),
		})
	}

	inv := LayerInventory{
		Name:     s.keyFormatter.LayerName(),
		Provides: provides,
	}

	if s.keyFormatter.LayerKind() == "" {
		inv.Note = "custom provider — keys not enumerable"
		// Drop key strings — they cannot be rendered for custom providers.
		for i := range inv.Provides {
			inv.Provides[i].Key = ""
		}
	}

	return inv, nil
}

// collectAliases walks the model and returns a map of fieldUID →
// canonical alias path (the dash-separated form StringBasedBuilder
// recognises). Implemented in inventory_walk.go.
```

(The `collectAliases` helper is implemented in Task 9 — until then leave it as a stub to keep this task self-contained.)

- [ ] **Step 3: Add the test seam `NewStringBasedBuilderForTest`**

Append to `sbased.go`:

```go
// NewStringBasedBuilderForTest constructs a StringBasedBuilder with a
// synthetic KeyFormatter. Intended solely for tests that need to
// exercise ReportInventory without going through the layer wrappers in
// builders.go.
//
// kind must be one of "env", "cmdline", "file", or "" (nil formatter
// for custom-provider behaviour). For "env" / "file", metaOrPrefix is
// the prefix or file id.
func NewStringBasedBuilderForTest(
	provider StringValuesProvider,
	kind, metaOrPrefix string,
) (*StringBasedBuilder, error) {
	var kf KeyFormatter
	switch kind {
	case "env":
		kf = newEnvKeyFormatter(metaOrPrefix)
	case "cmdline":
		kf = newCmdlineKeyFormatter()
	case "file":
		kf = newFileKeyFormatter(metaOrPrefix)
	case "":
		kf = newNilKeyFormatter(metaOrPrefix)
	default:
		return nil, fmt.Errorf("unknown key-formatter kind %q", kind)
	}

	return newStringBasedBuilderWithFormatter(provider, kf)
}
```

- [ ] **Step 4: Add a temporary stub for `collectAliases`**

Add to `sbased.go` (will be replaced in Task 9):

```go
// collectAliases returns the field-uid → alias-path map for the given
// model. Stubbed; full implementation lives in inventory_walk.go.
func collectAliases(m ModelInterface) (map[string]string, error) {
	return nil, fmt.Errorf("collectAliases: not yet implemented")
}
```

- [ ] **Step 5: Verify build**

Run: `go build ./...`
Expected: success. The unit test from Step 1 will fail at runtime because `collectAliases` returns the stub error — that's expected; it starts passing in Task 9.

- [ ] **Step 6: Commit**

```bash
git add sbased.go inventory_reporter_test.go
git commit -m "dsco: add ReportInventory to StringBasedBuilder (with stub alias collector)"
```

---

## Task 7: Implement ReportInventory on StructBuilder

**Files:**
- Modify: `structs.go`
- Test: `inventory_reporter_test.go` (extend)

- [ ] **Step 1: Write the failing unit-level test (no Compute dependency)**

Append to `inventory_reporter_test.go`:

```go
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
	assert.NotContains(t, uids, "Host", "nil-pointer fields must not appear")
}
```

(`dsco.BuildModel` is added in Task 8; this test will fail to compile until then.)

- [ ] **Step 2: Implement `ReportInventory` on `StructBuilder`**

Append to `structs.go`:

```go
// ReportInventory implements InventoryReporter by enumerating every
// non-nil field of the source struct and recording its value as a
// FieldProvision. No I/O is performed.
func (s *StructBuilder) ReportInventory(
	m ModelInterface,
) (LayerInventory, error) {
	const errCtx = "reporting struct inventory"

	values := m.GetFieldValuesFor(s.id, s.value)

	provides := make([]FieldProvision, 0, len(values))
	for fieldUID, v := range values {
		// Dereference pointer values so the report shows the user-visible
		// scalar, not a *T pointer.
		val := v.Value
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}

		provides = append(provides, FieldProvision{
			FieldUID: fieldUID,
			Value:    val.Interface(),
		})
	}

	return LayerInventory{
		Name:     "struct:" + s.id,
		Provides: provides,
	}, nil
}
```

- [ ] **Step 3: Verify build (test will fail at link step until Task 8)**

Run: `go build ./...`
Expected: success.

Run: `go test -run TestStructBuilderReportInventoryUnit ./...`
Expected: fails on missing `dsco.BuildModel`. That's expected — `BuildModel` lands in Task 8 and the test starts passing then.

- [ ] **Step 4: Commit**

```bash
git add structs.go inventory_reporter_test.go
git commit -m "dsco: add ReportInventory to StructBuilder"
```

---

## Task 8: Add PrepareInventoryWalk helper and refactor filler.go phases

**Files:**
- Create: `inventory_walk.go`
- Test: `inventory_walk_test.go`
- Modify: `filler.go`

- [ ] **Step 1: Write the failing test**

Create `inventory_walk_test.go`:

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run PrepareInventoryWalk ./...`
Expected: build error — undefined references.

- [ ] **Step 3: Refactor the first two phases of filler.go into shared helpers**

Modify `filler.go`. Replace `generateModel` and `generateBuilders` (lines 60–81) with calls into shared package-level helpers:

```go
func (c *dscoContext) generateModel() {
	if c.err.None() {
		m, err := buildModel(c.inputModelRef)
		if err != nil {
			c.err.Add(err)
			return
		}
		c.model = m
	}
}

func (c *dscoContext) generateBuilders() {
	if c.err.None() {
		builders, err := c.layers.GetPolicies()
		if err != nil {
			c.err.Add(err)
			return
		}
		c.builders = builders
	}
}

// buildModel constructs the model from a pointer-to-struct configuration.
func buildModel(inputModelRef any) (ModelInterface, error) {
	const errCtx = "building model"

	t := reflect.TypeOf(inputModelRef)
	if t == nil || t.Kind() != reflect.Pointer {
		return nil, fmt.Errorf(
			"%s: cfg must be a pointer to a struct", errCtx,
		)
	}

	m, err := model2.NewModel(t.Elem())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCtx, err)
	}

	return m, nil
}
```

- [ ] **Step 4: Create `inventory_walk.go`**

```go
package dsco

import (
	"fmt"
)

// InventoryWalk holds the prepared model and per-layer reporters used by
// the inventory sub-package to compute a Report. It contains no live
// configuration values — only structural metadata.
type InventoryWalk struct {
	Model     ModelInterface
	Reporters []InventoryReporter
}

// BuildModel constructs the configuration model from a pointer-to-struct
// value, mirroring the model-build phase of Fill. Exposed for the
// inventory sub-package; callers should generally use Fill or
// PrepareInventoryWalk instead.
func BuildModel(cfg any) (ModelInterface, error) {
	return buildModel(cfg)
}

// PrepareInventoryWalk constructs the model and the per-layer
// InventoryReporter list without performing any I/O. Used by the
// inventory sub-package to compute a Report.
//
// Pattern: Factory — assembles the structural inputs needed for an
// inventory walk.
func PrepareInventoryWalk(
	cfg any,
	layers ...Layer,
) (*InventoryWalk, error) {
	const errCtx = "preparing inventory walk"

	m, err := buildModel(cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCtx, err)
	}

	policies, err := Layers(layers).GetPolicies()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCtx, err)
	}

	reporters := make([]InventoryReporter, 0, len(policies))
	for i, p := range policies {
		fvg, ok := p.getFieldValuesGetter().(InventoryReporter)
		if !ok {
			return nil, fmt.Errorf(
				"%s: layer #%d does not implement InventoryReporter", errCtx, i,
			)
		}
		reporters = append(reporters, fvg)
	}

	return &InventoryWalk{
		Model:     m,
		Reporters: reporters,
	}, nil
}
```

- [ ] **Step 5: Add `getFieldValuesGetter` accessor to the policy interface**

Modify `policy.go`. Add to the interface (or as a method on each wrapper) — locate the `constraintLayerPolicy` interface and append:

```go
// getFieldValuesGetter exposes the wrapped FieldValuesGetter for inventory
// walks; not part of the Fill code path.
getFieldValuesGetter() FieldValuesGetter
```

Add the implementation to both `strictLayer` and `normalLayer` wrappers in the same file:

```go
func (l *strictLayer) getFieldValuesGetter() FieldValuesGetter { return l.FieldValuesGetter }
func (l *normalLayer) getFieldValuesGetter() FieldValuesGetter { return l.FieldValuesGetter }
```

(If the wrapper field name differs, adjust accordingly. Read `policy.go` first to confirm.)

- [ ] **Step 6: Run tests to verify they pass**

Run: `go test -run "PrepareInventoryWalk|BuildModel" ./...`
Expected: both tests pass.

Run: `go test -race ./...`
Expected: full suite still passes.

- [ ] **Step 7: Commit**

```bash
git add inventory_walk.go inventory_walk_test.go filler.go policy.go
git commit -m "dsco: extract buildModel helper and add PrepareInventoryWalk"
```

---

## Task 9: Implement collectAliases (replace stub from Task 6)

**Files:**
- Modify: `sbased.go` (replace stub)
- Test: `inventory_reporter_test.go` (extend)

- [ ] **Step 1: Read the model package to find the alias enumeration entry point**

Run: `grep -nR "FieldUID\|UID" internal/model/`
Expected: locate the field-iteration helper in `internal/model/`.

If no public iterator exists, add one:

```go
// In internal/model/model.go (or appropriate file):
//
// EachLeaf walks all scalar leaf fields, calling fn with each field's
// UID, dot-separated Go path, alias-path (dash-separated), and reflect
// type.
func (m *Model) EachLeaf(fn func(uid, path, alias string, t reflect.Type))
```

(Coordinate the exact signature with the existing field-tree implementation. The plan assumes a UID + alias-path are extractable; if the existing model already exposes both, reuse that accessor.)

- [ ] **Step 2: Write the failing test for `collectAliases`**

Append to `inventory_reporter_test.go`:

```go
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

	assert.True(t, values["database-host"], "expected database-host")
	assert.True(t, values["port"], "expected port")
}
```

(`CollectAliasesForTest` is a tiny exported test seam added in step 4.)

- [ ] **Step 3: Replace the `collectAliases` stub in `sbased.go`**

Replace the stub:

```go
// collectAliases returns the field-uid → alias-path map for the given
// model. The alias path is dash-separated, lowercase, matching the form
// StringBasedBuilder.values keys use.
func collectAliases(m ModelInterface) (map[string]string, error) {
	const errCtx = "collecting aliases"

	out := make(map[string]string)

	if err := m.EachLeaf(func(uid, _, alias string, _ reflect.Type) {
		out[uid] = alias
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", errCtx, err)
	}

	return out, nil
}
```

If `EachLeaf` does not return an error in the existing model API, drop the inner error wrap. Adapt to the actual signature.

- [ ] **Step 4: Add the test seam in `inventory_walk.go`**

Append to `inventory_walk.go`:

```go
// CollectAliasesForTest exposes collectAliases for testing only. It is
// intentionally not part of the public API.
func CollectAliasesForTest(m ModelInterface) (map[string]string, error) {
	return collectAliases(m)
}
```

- [ ] **Step 5: Run tests**

Run: `go test -race ./...`
Expected: the unit-level reporter tests from Tasks 6 and 7 now pass; the new alias-collection test passes.

- [ ] **Step 6: Commit**

```bash
git add sbased.go inventory_walk.go inventory_reporter_test.go internal/model/
git commit -m "dsco: implement collectAliases via model EachLeaf walk"
```

---

## Task 10: Add Mockery configuration for InventoryReporter

**Files:**
- Modify: `.mockery.yaml`
- Generate: `mock_InventoryReporter_test.go`

- [ ] **Step 1: Read current `.mockery.yaml`**

Run: `cat .mockery.yaml`

- [ ] **Step 2: Add an entry for InventoryReporter**

Append to `.mockery.yaml` under `packages.<dsco-root>`:

```yaml
packages:
  github.com/byte4ever/dsco:
    interfaces:
      # ... existing entries
      InventoryReporter:
```

(Match the indentation and config style of existing entries.)

- [ ] **Step 3: Generate the mock**

Run: `make mocks` (per CLAUDE.md note that mockery uses make, not go generate)
Expected: `mock_InventoryReporter_test.go` created in the dsco root.

- [ ] **Step 4: Verify it builds**

Run: `go build ./...`
Expected: success.

- [ ] **Step 5: Commit**

```bash
git add .mockery.yaml mock_InventoryReporter_test.go
git commit -m "dsco: generate MockInventoryReporter via mockery"
```

---

## Task 11: Implement inventory.Compute (collect + reduce)

**Files:**
- Modify: `inventory/inventory.go`
- Test: `inventory/inventory_test.go`

- [ ] **Step 1: Write the failing test**

Create `inventory/inventory_test.go`:

```go
package inventory_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/inventory"
)

// TestComputeCanonicalKeyHighestPrecedenceWins verifies that when env
// and cmdline both can supply the same field, the cmdline key (last in
// the layer list) wins.
func TestComputeCanonicalKeyHighestPrecedenceWins(t *testing.T) {
	t.Parallel()

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
	assert.Equal(t, "cmdline", report.Fields[0].Key.Layer)
	assert.Equal(t, "--host=", report.Fields[0].Key.Key)
}

// TestComputeSatisfiedByDefaults verifies that struct-layer values appear
// in Field.Satisfied while string-layer keys remain in Field.Key.
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

// TestComputeRejectsNonPointerCfg verifies the error path.
func TestComputeRejectsNonPointerCfg(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Host *string `yaml:"host"`
	}
	_, err := inventory.Compute(cfg{}, dsco.WithEnvLayer("MYAPP"))
	require.Error(t, err)
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./inventory/...`
Expected: build error — `Compute` undefined.

- [ ] **Step 3: Implement Compute**

Append to `inventory/inventory.go`:

```go
import (
	"fmt"
	"reflect"
	"sort"

	"github.com/byte4ever/dsco"
)

// Compute walks the model and layer builders exactly like dsco.Fill, but
// instead of loading values it returns the canonical key each
// string-based layer would accept for every required leaf field of cfg.
// No environment variables, command-line arguments, or files are read.
//
// Pattern: Factory — assembles a Report from model + layers without I/O.
func Compute(cfg any, layers ...dsco.Layer) (*Report, error) {
	const errCtx = "computing inventory"

	walk, err := dsco.PrepareInventoryWalk(cfg, layers...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCtx, err)
	}

	perLayer := make([]dsco.LayerInventory, 0, len(walk.Reporters))
	for i, r := range walk.Reporters {
		inv, err := r.ReportInventory(walk.Model)
		if err != nil {
			return nil, fmt.Errorf("%s: layer #%d: %w", errCtx, i, err)
		}
		perLayer = append(perLayer, inv)
	}

	return reduce(walk.Model, perLayer), nil
}

// reduce collapses per-layer reports into one Field per leaf, applying
// precedence rules: the last string-based layer that can supply a field
// wins for Key; any struct layer that bakes in a value populates
// Satisfied.
func reduce(m dsco.ModelInterface, perLayer []dsco.LayerInventory) *Report {
	type leaf struct {
		uid    string
		path   string
		goType string
	}

	var leaves []leaf
	_ = m.EachLeaf(func(uid, path string, _ string, t reflect.Type) {
		leaves = append(leaves, leaf{uid: uid, path: path, goType: t.String()})
	})

	fields := make([]Field, 0, len(leaves))
	for _, lf := range leaves {
		f := Field{Path: lf.path, GoType: lf.goType}

		// Walk layers in declaration order; later ones override earlier
		// for Key (per dsco precedence).
		for _, inv := range perLayer {
			for _, p := range inv.Provides {
				if p.FieldUID != lf.uid {
					continue
				}
				if p.Value != nil {
					f.Satisfied = &Satisfaction{
						LayerID: trimStructPrefix(inv.Name),
						Value:   p.Value,
					}
				}
				if p.Key != "" {
					f.Key = &KeySpec{
						Layer: layerKindFromName(inv.Name),
						Key:   p.Key,
					}
				}
			}
		}

		fields = append(fields, f)
	}

	sort.Slice(fields, func(i, j int) bool { return fields[i].Path < fields[j].Path })

	return &Report{
		Type:   m.TypeName(),
		Fields: fields,
	}
}

// trimStructPrefix strips the "struct:" prefix from a layer Name to
// expose just the user-supplied id.
func trimStructPrefix(name string) string {
	const prefix = "struct:"
	if len(name) > len(prefix) && name[:len(prefix)] == prefix {
		return name[len(prefix):]
	}
	return name
}

// layerKindFromName extracts the kind (e.g. "env") from a layer Name
// like "env:MYAPP". Returns the whole name if no colon is present
// (e.g. "cmdline").
func layerKindFromName(name string) string {
	for i := 0; i < len(name); i++ {
		if name[i] == ':' {
			return name[:i]
		}
	}
	return name
}
```

- [ ] **Step 4: Run tests**

Run: `go test -race ./inventory/...`
Expected: all four Compute tests pass.

- [ ] **Step 5: Run the full suite**

Run: `go test -race ./...`
Expected: all tests pass — the per-builder unit tests, the reporter walks, and Compute's reduction logic all green.

- [ ] **Step 6: Commit**

```bash
git add inventory/inventory.go inventory/inventory_test.go
git commit -m "inventory: implement Compute with precedence and satisfaction reduction"
```

---

## Task 12: Implement normalizeValue and custom marshalers

**Files:**
- Create: `inventory/normalize.go`
- Test: `inventory/normalize_test.go`

- [ ] **Step 1: Write the failing test**

Create `inventory/normalize_test.go`:

```go
package inventory

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, c.want, normalizeValue(c.in))
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

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./inventory/...`
Expected: build error — `normalizeValue` undefined.

- [ ] **Step 3: Implement normalize.go**

Create `inventory/normalize.go`:

```go
package inventory

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/goccy/go-yaml"
)

// normalizeValue converts fmt.Stringer values (time.Duration, time.Time,
// *url.URL, …) to their String() form so JSON / YAML / text output stays
// human-readable. Primitives and plain structs pass through.
func normalizeValue(v any) any {
	if v == nil {
		return nil
	}

	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}

	return v
}

// MarshalJSON implements json.Marshaler so Satisfaction.Value is
// normalized before serialization.
func (s Satisfaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LayerID string `json:"layer_id"`
		Value   any    `json:"value"`
	}{
		LayerID: s.LayerID,
		Value:   normalizeValue(s.Value),
	})
}

// MarshalYAML implements yaml.InterfaceMarshaler so Satisfaction.Value is
// normalized before serialization.
func (s Satisfaction) MarshalYAML() (any, error) {
	return struct {
		LayerID string `yaml:"layer_id"`
		Value   any    `yaml:"value"`
	}{
		LayerID: s.LayerID,
		Value:   normalizeValue(s.Value),
	}, nil
}

// Compile-time interface assertions.
var (
	_ json.Marshaler = Satisfaction{}
	_ yaml.InterfaceMarshaler = Satisfaction{}
)
```

- [ ] **Step 4: Run tests**

Run: `go test -race ./inventory/...`
Expected: all normalize tests pass.

- [ ] **Step 5: Commit**

```bash
git add inventory/normalize.go inventory/normalize_test.go
git commit -m "inventory: normalize Stringer values for JSON/YAML output"
```

---

## Task 13: Implement WriteJSON and WriteYAML

**Files:**
- Create: `inventory/format.go`
- Test: `inventory/format_test.go`
- Create: `inventory/testdata/sample.json`, `sample.yaml`

- [ ] **Step 1: Write the failing tests**

Create `inventory/format_test.go`:

```go
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
```

- [ ] **Step 2: Run tests to verify failure mode**

Run: `go test ./inventory/...`
Expected: build error — `WriteJSON`, `WriteYAML` undefined.

- [ ] **Step 3: Implement format.go**

Create `inventory/format.go`:

```go
package inventory

import (
	"fmt"
	"io"

	"github.com/goccy/go-json"
	"github.com/goccy/go-yaml"
)

// WriteJSON writes the inventory as JSON via github.com/goccy/go-json.
// Indentation is two spaces; output ends with a trailing newline.
func (r *Report) WriteJSON(w io.Writer) error {
	const errCtx = "writing JSON inventory"

	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	if _, err := w.Write([]byte{'\n'}); err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	return nil
}

// WriteYAML writes the inventory as YAML via github.com/goccy/go-yaml.
// Output ends with a trailing newline.
func (r *Report) WriteYAML(w io.Writer) error {
	const errCtx = "writing YAML inventory"

	out, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	return nil
}
```

- [ ] **Step 4: Generate golden files**

Run: `mkdir -p inventory/testdata && go test -update ./inventory/...`
Expected: tests pass (because `-update` writes goldens).

- [ ] **Step 5: Run without `-update` to confirm goldens match**

Run: `go test ./inventory/...`
Expected: tests pass.

- [ ] **Step 6: Inspect goldens visually**

Run: `cat inventory/testdata/sample.json && echo "---" && cat inventory/testdata/sample.yaml`
Expected: JSON and YAML look like the spec examples; durations stringified to `"30s"`.

- [ ] **Step 7: Add a writer-error test**

Append to `format_test.go`:

```go
type errWriter struct{ err error }

func (w errWriter) Write([]byte) (int, error) { return 0, w.err }

// TestWriteJSONPropagatesWriterError covers the error branch when the
// underlying io.Writer fails.
func TestWriteJSONPropagatesWriterError(t *testing.T) {
	t.Parallel()
	err := fixtureReport().WriteJSON(errWriter{err: assert.AnError})
	require.ErrorIs(t, err, assert.AnError)
}

// TestWriteYAMLPropagatesWriterError covers the same for WriteYAML.
func TestWriteYAMLPropagatesWriterError(t *testing.T) {
	t.Parallel()
	err := fixtureReport().WriteYAML(errWriter{err: assert.AnError})
	require.ErrorIs(t, err, assert.AnError)
}
```

Run: `go test -race ./inventory/...`
Expected: all pass.

- [ ] **Step 8: Commit**

```bash
git add inventory/format.go inventory/format_test.go inventory/testdata/sample.json inventory/testdata/sample.yaml
git commit -m "inventory: add WriteJSON and WriteYAML with golden coverage"
```

---

## Task 14: Implement WriteText with column layout

**Files:**
- Create: `inventory/text.go`
- Test: `inventory/text_test.go`
- Create: `inventory/testdata/sample.txt`

- [ ] **Step 1: Write the failing tests**

Create `inventory/text_test.go`:

```go
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

	r := &inventory.Report{
		Type: "Cfg",
		Fields: []inventory.Field{
			{Path: "X", GoType: "*string"},
		},
	}
	var buf bytes.Buffer
	require.NoError(t, r.WriteText(&buf))

	out := buf.String()
	assert.Contains(t, out, "—", "empty cells must use em-dash")
}

// TestWriteTextTruncatesLongDefaults verifies values >40 chars get
// truncated with ellipsis.
func TestWriteTextTruncatesLongDefaults(t *testing.T) {
	t.Parallel()

	long := strings.Repeat("x", 80)
	r := &inventory.Report{
		Type: "Cfg",
		Fields: []inventory.Field{
			{
				Path: "X", GoType: "*string",
				Satisfied: &inventory.Satisfaction{LayerID: "d", Value: long},
			},
		},
	}
	var buf bytes.Buffer
	require.NoError(t, r.WriteText(&buf))

	assert.Contains(t, buf.String(), "…", "long values must end with ellipsis")
	assert.NotContains(t, buf.String(), long, "full value must not appear")
}

// TestWriteTextPropagatesWriterError covers the error branch.
func TestWriteTextPropagatesWriterError(t *testing.T) {
	t.Parallel()
	err := fixtureReport().WriteText(errWriter{err: assert.AnError})
	require.ErrorIs(t, err, assert.AnError)
}
```

- [ ] **Step 2: Run tests to verify failure**

Run: `go test ./inventory/...`
Expected: build error — `WriteText` undefined.

- [ ] **Step 3: Implement text.go**

Create `inventory/text.go`:

```go
package inventory

import (
	"fmt"
	"io"
	"strings"
)

const (
	emDash             = "—"
	textMinPath        = 20
	textMinType        = 10
	textMinKey         = 20
	textMaxValueLength = 40
)

// WriteText writes a human-readable, fixed-width-column inventory to w.
// Columns: PATH | TYPE | KEY | DEFAULT. Empty cells render as "—".
// Long default values are truncated with ellipsis. Output ends with a
// trailing newline.
func (r *Report) WriteText(w io.Writer) error {
	const errCtx = "writing text inventory"

	rows := buildTextRows(r)

	pathW := columnWidth(rows, 0, textMinPath)
	typeW := columnWidth(rows, 1, textMinType)
	keyW := columnWidth(rows, 2, textMinKey)

	var b strings.Builder
	fmt.Fprintf(&b, "TYPE: %s\n\n", r.Type)
	fmt.Fprintf(&b, "%-*s  %-*s  %-*s  %s\n",
		pathW, "PATH", typeW, "TYPE", keyW, "KEY", "DEFAULT")
	for _, row := range rows {
		fmt.Fprintf(&b, "%-*s  %-*s  %-*s  %s\n",
			pathW, row[0], typeW, row[1], keyW, row[2], row[3])
	}

	if _, err := io.WriteString(w, b.String()); err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	return nil
}

// buildTextRows converts r.Fields into the [path, type, key, default]
// quadruples used by WriteText.
func buildTextRows(r *Report) [][4]string {
	rows := make([][4]string, 0, len(r.Fields))
	for _, f := range r.Fields {
		rows = append(rows, [4]string{
			f.Path,
			f.GoType,
			renderTextKey(f.Key),
			renderTextDefault(f.Satisfied),
		})
	}
	return rows
}

// renderTextKey renders a KeySpec as "<layer>: <key>" or em-dash.
func renderTextKey(k *KeySpec) string {
	if k == nil {
		return emDash
	}
	return k.Layer + ": " + k.Key
}

// renderTextDefault renders a Satisfaction as "defaults=<value>" with
// truncation, or em-dash when nil.
func renderTextDefault(s *Satisfaction) string {
	if s == nil {
		return emDash
	}
	val := fmt.Sprintf("%v", normalizeValue(s.Value))
	if len(val) > textMaxValueLength {
		val = val[:textMaxValueLength-1] + "…"
	}
	return s.LayerID + "=" + val
}

// columnWidth returns max(len(rows[i][col]), minimum) — the width
// needed to fit every value in the column.
func columnWidth(rows [][4]string, col, minimum int) int {
	w := minimum
	for _, row := range rows {
		if len(row[col]) > w {
			w = len(row[col])
		}
	}
	return w
}
```

- [ ] **Step 4: Generate golden file**

Run: `go test -update ./inventory/...`
Expected: tests pass; `inventory/testdata/sample.txt` written.

- [ ] **Step 5: Inspect golden visually**

Run: `cat inventory/testdata/sample.txt`
Expected: matches the spec text example (column layout, em-dashes, defaults annotation).

- [ ] **Step 6: Run tests without -update to confirm**

Run: `go test -race ./inventory/...`
Expected: all pass.

- [ ] **Step 7: Commit**

```bash
git add inventory/text.go inventory/text_test.go inventory/testdata/sample.txt
git commit -m "inventory: add WriteText with column layout and truncation"
```

---

## Task 15: Branch coverage push to 100%

**Files:**
- Modify: any test files needed to fill coverage gaps

- [ ] **Step 1: Run coverage**

Run: `go test -race -coverprofile=cover.out -covermode=atomic ./inventory/... ./...`
Expected: succeeds.

- [ ] **Step 2: Identify gaps**

Run: `go tool cover -func=cover.out | grep -Ev "100\.0%|_test\.go|/mock_"`
Expected: lists every function below 100%.

- [ ] **Step 3: Add tests for each gap**

For each function < 100%, add a focused test in the appropriate `*_test.go`. Common gaps to anticipate:

- Custom-provider note path in `StringBasedBuilder.ReportInventory`: write a test using `WithStringValueProvider` and a tiny `NamedStringValuesProvider` stub. Assert `LayerInventory.Note` is set and `Provides[*].Key` is empty.
- Nil-formatter defensive branch in `StringBasedBuilder.ReportInventory`: construct a builder via the public `NewStringBasedBuilder` (which leaves `keyFormatter` nil) and call `ReportInventory` directly; assert the wrapped error.
- `Satisfaction.MarshalJSON` / `MarshalYAML` via `inventory.Compute` end-to-end test that round-trips through both encoders.

- [ ] **Step 4: Re-run coverage until clean**

Run: `go test -race -coverprofile=cover.out -covermode=atomic ./... && go tool cover -func=cover.out | grep -Ev "100\.0%|_test\.go|/mock_"`
Expected: empty output for all new packages (`inventory/`, `key_formatter.go`, `inventory_walk.go`, `inventory_reporter.go`, and the `ReportInventory` methods in `sbased.go` / `structs.go`).

- [ ] **Step 5: Commit**

```bash
git add .
git commit -m "inventory: bring new code to 100% coverage"
```

---

## Task 16: Linting and formatting

**Files:** all changed files

- [ ] **Step 1: Run `go generate`**

Run: `go generate ./...`
Expected: no diff (no generators added).

- [ ] **Step 2: Auto-fix lint issues**

Run: `golangci-lint run --fix`
Expected: zero remaining issues, or only manual-fix items.

- [ ] **Step 3: Fix struct alignment**

Run: `betteralign -apply ./...`
Expected: no errors; possibly reorders fields in new structs.

- [ ] **Step 4: Enforce 80-char line length**

Run: `golines --shorten-comments --chain-split-dots --max-len=80 --base-formatter=gofumpt -w ./...`
Expected: rewrites long lines.

- [ ] **Step 5: Final lint pass**

Run: `golangci-lint run`
Expected: clean.

- [ ] **Step 6: Final test pass**

Run: `go test -race -cover ./...`
Expected: clean.

- [ ] **Step 7: Commit**

```bash
git add .
git commit -m "inventory: lint, format, and align"
```

---

## Task 17: Documentation

**Files:**
- Modify: `README.md`
- Modify: `QUICKSTART.md`
- Modify: `doc.go`
- Create: `examples/inventory/main.go`

- [ ] **Step 1: Add an Inventory section to README.md**

Insert before the "License" or final section, after the existing layered-fill examples:

```markdown
## Inventory

Need to know exactly which keys to wire up before deploying? `inventory.Compute`
walks your config struct and layers without performing any I/O and returns
the canonical key each layer would accept for every required field.

\`\`\`go
import (
    "os"
    "github.com/byte4ever/dsco"
    "github.com/byte4ever/dsco/inventory"
)

report, err := inventory.Compute(&config,
    dsco.WithStructLayer(defaults, "defaults"),
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithCmdlineLayer(),
)
if err != nil {
    log.Fatal(err)
}
report.WriteText(os.Stdout) // or WriteJSON / WriteYAML
\`\`\`
```

(Match existing README formatting.)

- [ ] **Step 2: Add a small Inventory example to QUICKSTART.md**

Append a one-paragraph note pointing readers to the README section, with a tiny one-liner.

- [ ] **Step 3: Extend `doc.go`**

Append a short paragraph after the existing package-level overview:

```go
// # Inventory
//
// The inventory sub-package (github.com/byte4ever/dsco/inventory) computes
// a static list of configuration keys a Fill call would expect, with no
// I/O, suitable for "what do I need to set" diagnostics in operator
// tooling.
```

- [ ] **Step 4: Create a runnable example**

Create `examples/inventory/main.go`:

```go
package main

import (
	"log"
	"os"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/inventory"
)

type Config struct {
	Host    *string `yaml:"host"`
	Port    *int    `yaml:"port"`
	Verbose *bool   `yaml:"verbose"`
}

func main() {
	defaults := &Config{Port: dsco.R(8080), Verbose: dsco.R(false)}

	var c *Config
	report, err := inventory.Compute(&c,
		dsco.WithStructLayer(defaults, "defaults"),
		dsco.WithEnvLayer("MYAPP"),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := report.WriteText(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
```

- [ ] **Step 5: Verify the example builds**

Run: `go build ./examples/inventory/...`
Expected: success.

- [ ] **Step 6: Commit**

```bash
git add README.md QUICKSTART.md doc.go examples/inventory/main.go
git commit -m "docs: document inventory sub-package and add runnable example"
```

---

## Self-Review Checklist

After implementing, verify each item:

1. **Spec coverage**
   - [ ] `Compute` exists in `inventory/` (Task 11) — spec §"Public API"
   - [ ] `Report` / `Field` / `Satisfaction` / `KeySpec` types defined (Task 2)
   - [ ] `WriteText` / `WriteJSON` / `WriteYAML` (Tasks 13–14)
   - [ ] `InventoryReporter` interface defined (Task 3)
   - [ ] `StringBasedBuilder.ReportInventory` (Task 6)
   - [ ] `StructBuilder.ReportInventory` (Task 7)
   - [ ] Per-layer key formatting (Tasks 4–5)
   - [ ] `LayerInventory.Note` for custom providers (Task 6 step 3, Task 15)
   - [ ] Path-lex sort across all formats (Task 11 + golden files)
   - [ ] Stringer normalization for durations (Task 12)
   - [ ] 100% coverage on new code (Task 15)
   - [ ] Mockery for `InventoryReporter` (Task 10)
   - [ ] Refactor of first two filler phases into shared helpers (Task 8)
   - [ ] Documentation updates (Task 17)

2. **Placeholder scan** — no `TBD`, no `// implement later`, no "similar to Task N".

3. **Type / signature consistency**
   - `Compute(cfg any, layers ...dsco.Layer) (*Report, error)` everywhere
   - `ReportInventory(model ModelInterface) (LayerInventory, error)` consistent
   - `KeyFormatter` methods (`LayerKind`, `LayerName`, `FormatKey`) match across all four implementations.

4. **Frequent commits** — 16 commits planned (one per task except Task 1 verification).
