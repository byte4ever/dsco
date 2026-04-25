# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

dsco (pronounced /ˈdɪskoʊ/) is a Go configuration library implementing a layered configuration system for microservices. It supports command line arguments, environment variables, YAML files, and struct-based configurations with strict validation.

**Core Philosophy**: Enforces explicit configuration through pointer-based fields to prevent silent defaults. `nil` means "not configured", non-nil means "explicitly configured".

**User-facing docs**: `README.md` / `README_fr.md` (overview and motivation), `QUICKSTART.md` / `QUICKSTART_fr.md` (getting started), `WARP.md` (Warp terminal integration notes), and `doc.go` (Go package doc).

## Development Commands

```bash
# Build
go build ./...

# Run all tests with race detection and coverage
go test -race -cover ./...

# Run a single test
go test -v -run TestName ./path/to/package

# Run tests for a specific package
go test ./internal/model/...

# Linting (extensive golangci-lint configuration)
golangci-lint run

# Auto-fix lint issues
golangci-lint run --fix

# Formatting
gofumpt -w .
golines --max-len=80 --base-formatter=gofumpt -w .
```

## Architecture

### Core Data Flow

```
Fill(&config, layers...) → Layer Registration → Model Generation (reflection)
    → Value Collection → Precedence Resolution → Type Conversion (YAML) → Validation
```

Earlier layers override later ones. A layer that leaves a field nil falls through
to the next layer. Strict mode layers error if their values are not consumed
(either unmatched to config fields, or overridden by earlier layers).

### Fill Pipeline (`filler.go`)

`dscoContext` orchestrates filling through sequential phases. Each phase checks
`c.err.None()` before proceeding — early errors short-circuit later phases:

1. `generateModel()` — reflect on target struct type, build field tree
2. `generateBuilders()` — convert `Layer` slice into `constraintLayerPolicies`
3. `generateFieldValues()` — each policy extracts values from the model; strict layers are tracked in `mustBeUsed`
4. `fillIt()` — model merges layer values by precedence, fills target struct
5. `checkUnused()` — verify strict layers had all values consumed

### Type Conversion Mechanism

All string-based layers (cmdline, env, file, custom providers) convert strings to
Go types via `yaml.Unmarshal` in `sbased.go`. This means any type that YAML can
parse is automatically supported (`time.Duration`, `net.URL`, etc.). The
`StringBasedBuilder.Get()` method handles pointer and slice kinds separately.

### Key Interfaces

| Interface | Location | Purpose |
|-----------|----------|---------|
| `Layer` | `builders.go` | Registration via `register(*layerBuilder)` — unexported method |
| `constraintLayerPolicy` | `policy.go` | Wraps `FieldValuesGetter` with `isStrict()` flag |
| `FieldValuesGetter` | `oiface.go` | Extracts `fvalue.Values` from a `ModelInterface` |
| `ModelInterface` | `oiface.go` | Bridge between model and layers (Fill, ApplyOn, Expand) |
| `ValueGetter` | `internal/iface.go` | `Get(path, type)` — used by model to pull values from layers |
| `StructExpander` | `internal/iface.go` | `ExpandStruct(path, type)` — handles nested struct expansion |
| `StringValuesProvider` | `iface.go` | Public interface for custom providers |
| `NamedStringValuesProvider` | `iface.go` | Extends `StringValuesProvider` with `GetName()` |

### Key Files

| File | Purpose |
|------|---------|
| `filler.go` | Core orchestration, `Fill()` function, `dscoContext` |
| `builders.go` | Layer construction (`WithCmdlineLayer`, `WithEnvLayer`, etc.) and dedup logic |
| `sbased.go` | `StringBasedBuilder` — YAML-based type conversion, alias handling, struct expansion |
| `structs.go` | `StructBuilder` — backs `WithStructLayer`, validates source struct type matches model |
| `oiface.go` | `FieldValuesGetter` and `ModelInterface` — bridge interfaces between layers and model |
| `iface.go` | Public provider interfaces (`StringValuesProvider`, `NamedStringValuesProvider`) |
| `policy.go` | `strictLayer` / `normalLayer` wrappers |
| `convert.go` | Path conversion: dot-separated → dash-separated snake_case |
| `errors.go` | Error type definitions and layer-level dedup errors |
| `doc.go` | Package-level Go doc with quickstart and overview |

### Internal Packages

| Package | Purpose |
|---------|---------|
| `internal/model/` | Reflection-based struct scanning, field tree (`snode`/`vnode`), get/expand operations |
| `internal/cmdline/` | Command line parsing (`--key=value`) |
| `internal/env/` | Environment variable parsing (`PREFIX-KEY=value`) |
| `internal/kfile/` | File-based configuration (YAML/JSON) via `afero` filesystem |
| `internal/fvalue/` | Field value representation (`Value` = `reflect.Value` + location) |
| `internal/merror/` | Multi-error aggregation (`MError` — embedded by all aggregate error types) |
| `internal/plocation/` | Path location tracking for source attribution |
| `internal/ierror/` | Indexed error wrapping (layer index + context info) |
| `internal/utils/` | Key name conversion, snake_case utilities |

### Public Packages

| Package | Purpose |
|---------|---------|
| `svalue/` | String value container with location tracking (input to `StringValuesProvider`) |
| `ref/` | Pointer helper (`R[T](value T) *T`) — also re-exported as `dsco.R()` |
| `registry/` | Type registration for runtime type name info |
| `hit/` | Hash-based integer types and providers |
| `url/` | YAML-unmarshallable `net/url.URL` wrapper |

## Error Pattern

All error types follow the sentinel + typed error pattern:

```go
var ErrOverriddenKey = errors.New("overridden key")  // sentinel

type OverriddenKeyError struct {                       // typed, carries context
    Path, Location, OverrideLocation string
}

func (OverriddenKeyError) Is(err error) bool {         // matches sentinel
    return errors.Is(err, ErrOverriddenKey)
}
```

Aggregate errors embed `merror.MError` (e.g., `FillerErrors`, `LayerErrors`,
`ModelError`, `GetError`). Each has its own sentinel and `Is()` method.

## Layer Types and Precedence

Order layers from highest to lowest priority; the first layer to supply a field
wins. Later layers only fill fields left nil by all preceding layers:

```go
dsco.Fill(&config,
    dsco.WithCmdlineLayer(),                     // highest priority
    dsco.WithEnvLayer("MYAPP"),                  // middle priority
    dsco.WithStructLayer(defaults, "defaults"),  // lowest priority
)
```

**Strict layers** (`WithStrictEnvLayer`, etc.) error if values don't match config
fields or are overridden by an earlier layer.

**Dedup rules** (in `builders.go`): only one cmdline layer allowed; env layers
deduplicated by prefix; struct layers by pointer address and string ID; string
providers by name.

## Configuration Struct Pattern

All configuration fields must be pointers (except slices/maps):

```go
type Config struct {
    Host    *string        `yaml:"host"`
    Port    *int           `yaml:"port"`
    Timeout *time.Duration `yaml:"timeout"`
}

// Use dsco.R() helper to create pointer values
defaults := &Config{
    Host: dsco.R("localhost"),
    Port: dsco.R(8080),
}
```

## Testing Patterns

- All tests use `t.Parallel()`
- Table-driven tests with `testify/require` and `testify/assert`
- Tests in same package (not `_test` suffix package)
- Mocks generated via mockery (`.mockery.yaml`) — in-package, with expecter pattern

## Linting

The `.golangci.yaml` enables ~60+ linters with strict settings:
- **Line length**: 80 characters max (enforced by `lll` and `golines`)
- **Formatter**: `gofumpt` (stricter than `gofmt`)
- **Notable enabled**: `err113` (wrapped errors), `ireturn` (interface returns), `exhaustive` (switch completeness), `gochecknoglobals`, `forcetypeassert`, `godot` (comment periods)
- **`nolint` directives** are used sparingly with reason comments (e.g., `//nolint:ireturn // this is required`)

## Environment Variable Mapping

Format: `PREFIX-KEY=value`
- Hyphen (`-`) separates prefix from key
- Hyphen (`-`) separates nested struct levels
- Underscores (`_`) in yaml tags are preserved
- Keys must be UPPERCASE

Examples with `WithEnvLayer("MYAPP")`:
- `Host` (`yaml:"host"`) → `MYAPP-HOST`
- `Database.Host` (`yaml:"database"` + `yaml:"host"`) → `MYAPP-DATABASE-HOST`
- `MaxRetry` (`yaml:"max_retry"`) → `MYAPP-MAX_RETRY`
