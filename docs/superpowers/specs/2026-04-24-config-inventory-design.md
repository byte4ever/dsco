# Config Inventory Design

**Date:** 2026-04-24
**Status:** Approved (brainstorm phase) — pending implementation plan

## Problem

A dsco-managed service fails fast at startup when a required pointer field is
left `nil`. Operators currently learn what keys to wire up by reading the Go
config struct, mentally translating each field through every layer's naming
convention (`MYAPP-DATABASE-HOST`, `--database-host=`, `database.host`), and
guessing whether a defaults layer already covers it.

We want a static, no-I/O answer to "given this `Fill(&cfg, layers...)` call,
what is the exact list of keys an operator must provide for the service to
start?"

## Goals

- Enumerate every leaf field of a config struct alongside the canonical key
  each layer would accept.
- Honour layer precedence: when multiple string-based layers can supply the
  same field, surface the highest-precedence one.
- Mark fields already satisfied by a struct (defaults) layer.
- Output as JSON, YAML, or human-readable text.
- Perform no I/O — never read environment variables, command-line arguments,
  or files.
- Reach 100% test coverage on new code.

## Non-Goals

- A live status check ("is `MYAPP-DATABASE-HOST` currently set?"). Inventory
  is static; live checks remain a separate `Fill`-error concern.
- Auto-injecting a CLI flag into the user's program. The library exposes a
  function; the user wires their own subcommand.
- Enumerating slice/map elements (impossible statically).
- Surfacing alias keys. v1 shows the canonical key only.
- Custom string-provider key enumeration. v1 reports custom providers by name
  with a "keys not enumerable" note.

## Public API

A new sub-package `inventory/`, mirroring `dsco/svalue/`, `dsco/registry/`,
`dsco/url/`.

```go
// Package inventory computes the static list of configuration keys a
// dsco-managed struct expects, given a set of layers, without performing
// any I/O.
package inventory

import (
    "io"

    "github.com/byte4ever/dsco"
)

// Compute walks the model and layer builders exactly like dsco.Fill, but
// instead of loading values it returns the canonical key each string-based
// layer would accept for every required leaf field of cfg. No environment
// variables, command-line arguments, or files are read.
//
// Pattern: Factory — assembles a Report from model + layers without I/O.
func Compute(cfg any, layers ...dsco.Layer) (*Report, error)

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

// WriteText writes a flat human-readable inventory to w.
func (r *Report) WriteText(w io.Writer) error

// WriteJSON writes the inventory as JSON via github.com/goccy/go-json.
func (r *Report) WriteJSON(w io.Writer) error

// WriteYAML writes the inventory as YAML via github.com/goccy/go-yaml.
func (r *Report) WriteYAML(w io.Writer) error
```

### API rationale

- Sub-package avoids type-name stutter (`dsco.InventoryReport`...) and keeps
  the dsco root focused on the layered-fill API.
- `Compute` is a verb returning the noun `*Report`, avoiding a function-vs-type
  name collision.
- All exported types carry `json` and `yaml` struct tags so callers may also
  marshal directly with their preferred encoder if they don't use `Write*`.
- `Satisfied *Satisfaction` and `Key *KeySpec` are pointer-nullable per dsco's
  own "nil = absent" philosophy; `omitempty` removes them cleanly when
  inapplicable.
- Three concrete `Write*` methods rather than a `Formatter` strategy interface
  — three formats with distinct serialization shapes do not share enough
  behaviour to justify the abstraction (YAGNI).

## Internal Architecture

### New layer-builder interface

```go
// InventoryReporter is implemented by layer builders that can describe what
// they would contribute to an inventory without performing I/O.
//
// Pattern: Strategy — each layer reports its own contribution shape (key
// form for string-based layers, baked-in values for struct layers).
type InventoryReporter interface {
    ReportInventory(model ModelInterface) (LayerInventory, error)
}

// LayerInventory is one layer's contribution to a Report.
type LayerInventory struct {
    Name     string             // "env:MYAPP", "cmdline", "file:<id>", "struct:<id>", "<custom>"
    Note     string             // optional — e.g. "custom provider — keys not enumerable"
    Provides []FieldProvision   // one entry per field this layer can populate
}

// FieldProvision is one (field, layer) pair from a layer's perspective.
type FieldProvision struct {
    FieldUID string // matches the model's field UID
    Key      string // canonical key for string-based layers; empty for struct layers
    Value    any    // baked-in value for struct layers; nil for string-based layers
}
```

The interface is defined in the dsco root (consumed there) and implemented by
the existing builders in `builders.go`, `sbased.go`, and `structs.go`. One
method, well under the 5-method cap, single-responsibility.

### Pipeline

`Compute` mirrors `dscoContext` with a small parallel `inventoryContext`. The
first two phases (`generateModel`, `generateBuilders`) are extracted into
shared helpers used by both contexts. `Compute` deliberately does **not**
call `generateFieldValues` — that phase is what reads env/cmdline/files.

```
Compute(cfg, layers...)
  → newInventoryContext
  → .generateModel()       // shared with dscoContext
  → .generateBuilders()    // shared with dscoContext
  → .collectInventory()    // NEW — calls ReportInventory on each builder
  → .reduce()              // NEW — walks model leaves, picks canonical Key per field
  → *Report
```

`reduce` is pure data shaping: no builder calls, no I/O.

### Where the key strings come from

`env`, `cmdline`, and `file` layers all wrap `dsco.StringBasedBuilder`
(in `sbased.go`), which already holds the field→key alias map needed to
recognise incoming strings at fill time. We add a small unexported accessor
(e.g., `func (b *StringBasedBuilder) aliases() map[string]string`) used only
by the layer wrappers in the same package; each wrapper formats its own key
form (uppercase + dashes for env, `--name=` for cmdline, raw path for file).
Naming logic stays in **one** place per layer kind — the inventory phase
never reinvents it.

## Per-Layer Key Forms

| Layer | `LayerInventory.Name` | `KeySpec.Layer` | `KeySpec.Key` for `Database.Host` (yaml: `database.host`) |
|---|---|---|---|
| `WithEnvLayer("MYAPP")` | `env:MYAPP` | `env` | `MYAPP-DATABASE-HOST` |
| `WithCmdlineLayer()` | `cmdline` | `cmdline` | `--database-host=` |
| `WithFileLayer(...)` | `file:<id>` | `file` | `database.host` |
| `WithStructLayer(v, "defaults")` | `struct:defaults` | — | — (records `Satisfaction`, not a key) |
| Custom string provider | `<name>` | — | empty + `LayerInventory.Note: "custom provider — keys not enumerable"` |

`KeySpec.Layer` is the layer **kind** (suitable for grouping in tooling); the
full identifier lives in `LayerInventory.Name` for attribution.

## Edge Cases

- **Nested structs.** Recurse fully. Each scalar leaf becomes one `Field`
  entry; `Path` is dot-separated using Go field names. Nested structs never
  appear as their own entry, only their leaves do.
- **Slices / maps / arrays.** A single `Field` entry per such field. `GoType`
  is verbose (`"[]string"`, `"map[string]int"`, `"[3]int"`); `KeySpec.Key` is
  the string the layer accepts for whole-collection assignment (e.g.,
  `MYAPP-PORTS=8080,8081`). No element enumeration.
- **Aliases.** Surface only the canonical (primary) key in v1. A
  `KeySpec.Aliases []string` field can be added later if needed.
- **`GoType` formatting.** `reflect.Type.String()` directly:
  `"string"`, `"*int"`, `"[]string"`, `"time.Duration"`, `"*url.URL"`. The
  declared type is reported (so `*string` appears as `"*string"`).
- **Strict layers.** Strict-ness is irrelevant — strict only affects
  unconsumed-value errors at `Fill` time. The keys a strict layer accepts
  match a non-strict one. Not surfaced in the report.
- **Precedence.** When multiple string-based layers can supply the same field,
  the highest-precedence wins (the last layer in `Compute(...)` arg list that
  can supply it). Same model dsco already uses for value resolution.

## Errors

Reuse the existing `FillerErrors` family. Failures from `Compute`:

- Non-pointer config → wrapped `ErrFiller`.
- Duplicate cmdline layers, etc. → wrapped `ErrFiller` (same path as `Fill`).
- A reporter returning an error → wrapped with the layer index in the wrap
  string, mirroring `generateFieldValues`.

Callers can `errors.Is(err, dsco.ErrFiller)`.

## Output Formats

All three formats sort `Fields` by `Path` lexicographically — consistent and
grep-friendly.

### JSON (`WriteJSON`)

`goccy/go-json` marshal of `*Report`, indent two spaces.

```json
{
  "type": "github.com/example/myapp.Config",
  "fields": [
    {
      "path": "Database.Host",
      "go_type": "*string",
      "key": { "layer": "env", "key": "MYAPP-DATABASE-HOST" }
    },
    {
      "path": "Database.Port",
      "go_type": "*int",
      "satisfied": { "layer_id": "defaults", "value": 5432 },
      "key": { "layer": "cmdline", "key": "--database-port=" }
    },
    {
      "path": "Server.Timeout",
      "go_type": "*time.Duration",
      "satisfied": { "layer_id": "defaults", "value": "30s" }
    }
  ]
}
```

### YAML (`WriteYAML`)

`goccy/go-yaml` marshal of `*Report`. Same shape as JSON.

```yaml
type: github.com/example/myapp.Config
fields:
  - path: Database.Host
    go_type: '*string'
    key:
      layer: env
      key: MYAPP-DATABASE-HOST
  - path: Database.Port
    go_type: '*int'
    satisfied:
      layer_id: defaults
      value: 5432
    key:
      layer: cmdline
      key: --database-port=
  - path: Server.Timeout
    go_type: '*time.Duration'
    satisfied:
      layer_id: defaults
      value: 30s
```

### Text (`WriteText`)

Fixed-width columns, em-dash for empty cells.

```
TYPE: github.com/example/myapp.Config

PATH                  TYPE             KEY                              DEFAULT
Database.Host         *string          env: MYAPP-DATABASE-HOST         —
Database.Port         *int             cmdline: --database-port=        defaults=5432
Server.Timeout        *time.Duration   —                                defaults=30s
```

Rules:
- Column widths sized to longest entry; minimum widths `path=20`, `type=10`,
  `key=20`.
- Em-dash `—` (U+2014) for empty cells.
- `defaults=<value>` uses `fmt.Sprintf("%v", value)`; long values (>40 chars)
  truncated with `…`.
- All three writers terminate output with a single `\n`.

### Value normalization

For `Satisfaction.Value`, a `normalizeValue(any) any` helper converts any
`fmt.Stringer` (`time.Duration`, `time.Time`, `*url.URL`, …) to its `String()`
form before serialisation. Primitives and plain structs pass through.
Implemented as custom `MarshalJSON` / `MarshalYAML` on `Satisfaction` so all
three formats share the logic.

## Testing Strategy

**Coverage target: 100%** for all new code (the `inventory/` sub-package, new
`ReportInventory` methods on every builder, new internal accessors). Verified
via `go test -race -coverprofile=cover.out -covermode=atomic ./...` — the spec
is incomplete until `go tool cover -func=cover.out | grep -v 100.0%` prints
nothing for new packages.

### Mockery (testing only)

Three mocks added to `.mockery.yaml`, all in-package, expecter pattern (per
CLAUDE.md convention).

| Mock | Location | Purpose |
|---|---|---|
| `MockInventoryReporter` | `mock_InventoryReporter_test.go` (dsco root) | Isolates `collectInventory` and `reduce` for branch-coverage tests (overlapping FieldUIDs, empty Provides, reporter errors) without spinning up a real env+cmdline+file stack. |
| `MockModelInterface` (extend existing) | `mock_ModelInterface_test.go` | Reuse for `Compute` error-path tests (`TypeName`, `Expand` failures). |
| `MockStringBasedAccessor` | `mock_StringBasedAccessor_test.go` (dsco root) | Only added if the sbased→layer accessor becomes its own interface. Lets each layer wrapper test key formatting independently. |

**No mocks for:** `Layer`/`Option` builders (real builders are cheap, mocking
hides actual key-format behaviour), `io.Writer` (use `bytes.Buffer`),
`Stringer` types (use real `time.Duration`, `time.Time`, `*url.URL`).

### Test layout

1. **Per-builder `ReportInventory` tests** — table-driven, one row per model
   shape (simple field, nested struct, slice, map, mixed nil/set struct
   values, empty struct).
2. **Pipeline tests** (`inventory/inventory_test.go`) — drive every branch in
   `reduce` via `MockInventoryReporter`: precedence resolution,
   satisfaction-only fields, key-only fields, empty inventory, reporter error
   propagation.
3. **`normalizeValue` tests** (`inventory/normalize_test.go`) — every type
   class: Stringer, primitive, struct, nil, pointer-to-Stringer,
   slice-of-Stringer.
4. **Format tests** (`inventory/format_test.go`):
   - Golden files in `inventory/testdata/` for JSON / YAML / text, regenerated
     via `-update` flag.
   - Branch coverage for `WriteText`: long-value truncation, all-empty cells,
     mixed widths, Unicode in paths.
   - Failing-writer test using a tiny in-test `errWriter` struct (no mock
     needed) to confirm `Write*` propagates errors.
5. **Error-path tests**:
   - Non-pointer config → wrapped `ErrFiller`.
   - Duplicate cmdline layers → wrapped `ErrFiller`.
   - Reporter returns error → propagated with layer index in wrap string.
6. **Custom-provider note test** — assert `LayerInventory.Note` is set and
   survives JSON/YAML roundtrip.

### Conventions (per CLAUDE.md + golang skill)

- All tests `t.Parallel()`.
- Tests in `_test` package (e.g., `inventory_test`) to verify the public API.
- `testify/require` for nil/error checks; `testify/assert` for value
  equality.
- Each mock's purpose documented at the top of its test file ("why this mock
  exists" line).

## Out of Scope

- A CLI subcommand or auto-injected flag in the user's binary.
- Live state checks against the actual environment.
- Element enumeration for slice/map fields.
- Custom string-provider key enumeration.
- Alias surfacing.

## Risks

- **Accessor surface growth.** Adding `aliases()` (or similar) on
  `StringBasedBuilder` widens its internal API. Mitigation: keep the accessor
  unexported; only the layer wrappers in the same `dsco` package consume it;
  no external users.
- **Pipeline duplication drift.** Extracting the first two phases into shared
  helpers risks the inventory and fill paths diverging silently. Mitigation:
  one shared helper per phase, used identically by both contexts; covered by
  pipeline tests on both sides.
- **Golden-file churn.** Format tests against goldens are sensitive to
  formatting tweaks. Mitigation: structure goldens to minimize cosmetic noise
  (sorted, deterministic widths); regenerate via `-update` flag is one-line.
