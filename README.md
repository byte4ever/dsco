# dsco

**Stop deploying microservices with broken configuration.**

dsco is a Go configuration library that makes misconfiguration impossible.
No more silent defaults. No more "it works on my machine." No more 3 AM pages
because someone forgot to set `DATABASE_PASSWORD` in production.

```go
// 30 seconds to bulletproof configuration
var config *Config
dsco.Fill(&config,
    dsco.WithCmdlineLayer(),                     // Quick local overrides
    dsco.WithEnvLayer("MYAPP"),                  // Container/K8s config
    dsco.WithStructLayer(defaults, "defaults"),  // Dev defaults baked in
)
// Missing config? App won't start. You'll know immediately.
```

[![Go](https://github.com/byte4ever/dsco/actions/workflows/go.yml/badge.svg)](https://github.com/byte4ever/dsco/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/byte4ever/dsco.svg)](https://pkg.go.dev/github.com/byte4ever/dsco)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.21-61CFDD.svg?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/byte4ever/dsco?style=flat-square)](https://goreportcard.com/report/github.com/byte4ever/dsco)
[![Test Coverage](https://api.codeclimate.com/v1/badges/c64776c8e19d20057719/test_coverage)](https://codeclimate.com/github/byte4ever/dsco/test_coverage)

[Français](README_fr.md) | English

---

## Why dsco?

**Traditional Go configuration is dangerous:**

```go
type Config struct {
    Host string  // Is "" intentional or did someone forget to set it?
    Port int     // Is 0 a valid port or a missing value?
}
```

**dsco makes intent explicit:**

```go
type Config struct {
    Host *string `yaml:"host"`  // nil = not configured (fail fast)
    Port *int    `yaml:"port"`  // nil = not configured (fail fast)
}
```

| Problem | dsco Solution |
|---------|---------------|
| Service starts with missing DB password | Fails immediately with clear error |
| Zero value `0` masks missing port config | `nil` explicitly means "not set" |
| Config works locally, breaks in prod | Same validation everywhere |
| "Which env var overrode what?" | Full audit trail with source tracking |

---

## Quick Start

```bash
go get github.com/byte4ever/dsco
```

```go
package main

import (
    "fmt"
    "log"

    "github.com/byte4ever/dsco"
)

type Config struct {
    Host *string `yaml:"host"`
    Port *int    `yaml:"port"`
}

func main() {
    var config *Config

    _, err := dsco.Fill(&config,
        // Layer 1: command line (highest priority)
        dsco.WithCmdlineLayer(),
        // Layer 2: environment variables
        dsco.WithEnvLayer("MYAPP"),
        // Layer 3: defaults (lowest priority)
        dsco.WithStructLayer(&Config{
            Host: dsco.R("localhost"),
            Port: dsco.R(8080),
        }, "defaults"),
    )
    if err != nil {
        log.Fatal(err)  // Missing config? Crash here, not in production.
    }

    fmt.Printf("Server: %s:%d\n", *config.Host, *config.Port)
}
```

```bash
# Just works with defaults
./myapp

# Override via environment (Kubernetes/Docker)
MYAPP-HOST=api.prod.internal MYAPP-PORT=9000 ./myapp

# Override via command line (local dev)
./myapp --host=staging.example.com --port=9000
```

**New to dsco?** The [Quick Start Guide](QUICKSTART.md) covers all concepts
step-by-step.

---

## Table of Contents

- [Key Features](#key-features)
- [The Safety Design](#the-safety-design)
- [You're In Control](#youre-in-control)
- [Layer Types](#layer-types)
- [Environment Variables](#environment-variables)
- [Architecture](#architecture)
- [Configuration Patterns](#configuration-patterns)
- [Error Handling](#error-handling)
- [Advanced Usage](#advanced-usage)
- [Inventory](#inventory)
- [API Reference](#api-reference)
- [Examples](#examples)
- [Contributing](#contributing)

---

## Key Features

| Feature | Benefit |
|---------|---------|
| **Layered Priority** | Cmdline → env vars → struct defaults. First wins. |
| **Pointer-Based Safety** | `nil` = not configured. No silent zero values. |
| **Strict Mode** | Catch typos and unwanted overrides immediately. |
| **Source Tracking** | Know exactly where every value came from. |
| **Multi-Source** | Cmdline, env vars, files, structs, custom providers. |
| **Type Safety** | Automatic conversion with clear parse errors. |
| **Alias Support** | `--db-host` instead of `--database-host`. |
| **Minimal Deps** | Only `yaml.v3` and `afero`. |

---

## The Safety Design

### Why Pointers?

```go
// DANGEROUS: Is Port 0 intentional or missing?
type Config struct {
    Port int
}

// SAFE: nil clearly means "not configured"
type Config struct {
    Port *int `yaml:"port"`
}
```

**The `dsco.R()` helper makes pointer creation painless:**

```go
config := &Config{
    Host:    dsco.R("localhost"),   // dsco.R[T](v T) *T
    Port:    dsco.R(8080),
    Timeout: dsco.R(30 * time.Second),
}
```

### Fail-Fast Guarantee

dsco ensures **all configuration is complete before your app starts**:

```go
// This FAILS - Password is nil
dsco.Fill(&config,
    dsco.WithStructLayer(&DatabaseConfig{
        Host: dsco.R("localhost"),
        Port: dsco.R(5432),
        // Password not set - nil
    }, "defaults"),
)
// Error: "password" is not configured

// This SUCCEEDS - all fields explicitly set
dsco.Fill(&config,
    dsco.WithEnvLayer("DB"),  // DB-PASSWORD must be set
    dsco.WithStructLayer(&DatabaseConfig{
        Host: dsco.R("localhost"),
        Port: dsco.R(5432),
    }, "defaults"),
)
```

### Production Example

```go
type DatabaseConfig struct {
    Host     *string `yaml:"host"`
    Port     *int    `yaml:"port"`
    Username *string `yaml:"username"`
    Password *string `yaml:"password"`
    SSLMode  *string `yaml:"ssl_mode"`
}

_, err := dsco.Fill(&dbConfig,
    // Secrets from Vault/external system
    dsco.WithStringValueProvider(secretProvider),
    // Environment overrides
    dsco.WithStrictEnvLayer("DB"),
    // Base configuration
    dsco.WithStructLayer(&DatabaseConfig{
        Host:    dsco.R("postgres.prod.internal"),
        Port:    dsco.R(5432),
        SSLMode: dsco.R("require"),
        // Username/Password MUST come from higher layers
    }, "base"),
)

if err != nil {
    // Clear error: "username is not configured"
    log.Fatal("Configuration incomplete:", err)
}
```

---

## You're In Control

dsco gives you **complete control** over what's configurable, when, and by whom.

### The Progressive Exposure Pattern

Start with everything hardcoded, then progressively expose parameters as needed:

**Phase 1: All defaults in code**

```go
// Initial deployment - everything hardcoded, nothing external
dsco.Fill(&config,
    dsco.WithStructLayer(&Config{
        Host:       dsco.R("api.internal"),
        Port:       dsco.R(8080),
        MaxRetries: dsco.R(3),
        Timeout:    dsco.R(30 * time.Second),
        BatchSize:  dsco.R(100),
    }, "defaults"),
)
```

Your service runs perfectly. No external configuration needed. No environment
variables to forget. No config files to deploy.

**Phase 2: Expose what matters**

Later, you realize `Timeout` needs adjustment per environment:

```go
// Now Timeout can be overridden via environment, everything else stays fixed
dsco.Fill(&config,
    dsco.WithEnvLayer("MYSERVICE"),  // Only MYSERVICE-TIMEOUT needs to exist
    dsco.WithStructLayer(&Config{
        Host:       dsco.R("api.internal"),
        Port:       dsco.R(8080),
        MaxRetries: dsco.R(3),
        Timeout:    dsco.R(30 * time.Second),  // Default, but overridable
        BatchSize:  dsco.R(100),
    }, "defaults"),
)
```

**No recompilation required.** The code didn't change - you just added an env
layer. Operations can now tune `MYSERVICE-TIMEOUT=60s` without touching the
binary.

**Phase 3: Protect critical values**

Some defaults should **never** be overridden in production:

```go
dsco.Fill(&config,
    // These values are LOCKED - highest priority, strict enforcement
    dsco.WithStrictStructLayer(&Config{
        APIEndpoint: dsco.R("https://api.production.com"),
        AuditMode:   dsco.R(true),
    }, "immutable"),

    // Operational overrides allowed
    dsco.WithEnvLayer("MYSERVICE"),
    dsco.WithCmdlineLayer(),
)
```

Even if someone sets `MYSERVICE-API-ENDPOINT`, the strict struct layer wins
**and** raises an error about the attempted override.

### Why This Matters

| Traditional Approach | dsco Approach |
|---------------------|---------------|
| Expose everything upfront "just in case" | Start minimal, expose on demand |
| Config sprawl - hundreds of env vars | Only what's actually needed |
| No protection - any value can be overridden | Lock critical values with strict layers |
| Must redeploy to change exposure | Add layers without code changes |
| "What's the default?" - check docs/code | Defaults visible in layer definition |

### Real-World Scenarios

**Scenario 1: New service deployment**

```go
// Week 1: Ship with safe defaults, zero external config
dsco.Fill(&config, dsco.WithStructLayer(productionDefaults, "defaults"))
```

**Scenario 2: Ops needs to tune performance**

```go
// Week 3: Add env layer - ops can now adjust without new release
dsco.Fill(&config,
    dsco.WithEnvLayer("SVC"),
    dsco.WithStructLayer(productionDefaults, "defaults"),
)
// Ops sets SVC-CONNECTION-POOL-SIZE=50
```

**Scenario 3: Prevent accidental security misconfiguration**

```go
// Security audit: ensure TLS and audit logging can't be disabled
dsco.Fill(&config,
    dsco.WithStrictStructLayer(&Config{
        TLSEnabled:    dsco.R(true),
        AuditLogging:  dsco.R(true),
        MinTLSVersion: dsco.R("1.3"),
    }, "security"),
    dsco.WithEnvLayer("SVC"),
)
```

**You decide** what's flexible and what's fixed. dsco enforces your decisions.

---

## Layer Types

### Struct Layers (Defaults)

```go
dsco.WithStructLayer(&Config{
    Host: dsco.R("localhost"),
    Port: dsco.R(8080),
}, "defaults")

// Strict: errors if values are overridden
dsco.WithStrictStructLayer(&Config{
    APIEndpoint: dsco.R("https://api.prod.com"),
}, "immutable")
```

**Local development pattern** - zero config to start:

```go
dsco.Fill(&config,
    dsco.WithCmdlineLayer(),
    dsco.WithStructLayer(devDefaults, "dev"),
)
```

```bash
./myapp                    # Just works
./myapp --port=9000        # Quick override
./myapp --database-host=staging-db
```

### Command Line Layers

```go
dsco.WithCmdlineLayer()
dsco.WithStrictCmdlineLayer()  // Error on unknown flags

// With aliases
dsco.WithCmdlineLayer(
    dsco.WithAliases(map[string]string{
        "v": "verbose",
        "p": "port",
    }),
)
```

**Format**: `--key=value` (lowercase, hyphens for nested fields)

```bash
./myapp --host=localhost --database-port=5432
```

### Environment Variable Layers

```go
dsco.WithEnvLayer("MYAPP")
dsco.WithStrictEnvLayer("MYAPP")  // Error on unmatched vars
```

### Custom Providers

```go
type SecretProvider struct{}

func (s SecretProvider) GetName() string { return "vault" }
func (s SecretProvider) GetStringValues() svalue.Values {
    return svalue.Values{
        "database-password": &svalue.Value{
            Value:    fetchFromVault("db-password"),
            Location: "vault:db-password",
        },
    }
}

dsco.WithStringValueProvider(&SecretProvider{})
```

---

## Environment Variables

### Why Prefixes Matter

**Multi-container pods (Kubernetes):**

All containers in a pod share environment variables. Prefixes target specific
containers:

```yaml
env:
  - name: FRONTEND-PORT
    value: "8080"
  - name: BACKEND-PORT
    value: "3000"
```

**Avoid conflicts:**

Prevents collision with `PATH`, `HOME`, `HTTP_PROXY`, `DATABASE_URL`, etc.

**Multiple instances:**

```bash
WORKER1-QUEUE=high-priority ./worker &
WORKER2-QUEUE=low-priority ./worker &
```

### Choosing Good Prefixes

**Avoid generic prefixes** that cause confusion:

```bash
# BAD: Too generic
APP-HOST=...       # Which app?
SERVER-PORT=...    # Which server?
SERVICE-URL=...    # Meaningless
```

**Use specific, role-based prefixes:**

```bash
# GOOD: Clear and distinguishable
ORDERAPI-HOST=...           # Order API service
PAYMENTWORKER-TIMEOUT=...   # Payment background worker
EVENTCONSUMER-BATCH=...     # Event queue consumer
```

This makes debugging easier ("check INDEXER config"), Kubernetes manifests
self-documenting, and prevents cross-contamination in shared environments.

### Format

```
PREFIX-KEY=value
│      │
│      └─ UPPERCASE key (hyphens/underscores allowed)
└─ UPPERCASE prefix
```

### Mapping Examples

| Struct Field | YAML Tag | Environment Variable |
|--------------|----------|---------------------|
| `Host` | `host` | `MYAPP-HOST` |
| `MaxRetry` | `max_retry` | `MYAPP-MAX_RETRY` |
| `Database.Host` | `database.host` | `MYAPP-DATABASE-HOST` |
| `Database.PoolSize` | `database.pool_size` | `MYAPP-DATABASE-POOL_SIZE` |

**Rules:**
- Prefix and keys: UPPERCASE
- Prefix-to-key separator: hyphen (`-`)
- Nested struct separator: hyphen (`-`)
- Underscores in yaml tags: preserved

---

## Architecture

```mermaid
graph TB
    A[Configuration Sources] --> B[Layer Registration]
    B --> C[Model Generation]
    C --> D[Value Collection]
    D --> E[Precedence Resolution]
    E --> F[Type Conversion]
    F --> G[Validation]
    G --> H[Struct Filling]

    A1[Command Line] --> B1[CmdlineLayer]
    A2[Environment] --> B2[EnvLayer]
    A3[Files] --> B3[FileLayer]
    A4[Structs] --> B4[StructLayer]
    A5[Custom] --> B5[StringProviderLayer]

    B1 --> B
    B2 --> B
    B3 --> B
    B4 --> B
    B5 --> B
```

**Flow:**
1. **Layer Registration** - Sources register as layers
2. **Model Generation** - Struct analyzed via reflection
3. **Value Collection** - Each layer provides values
4. **Precedence Resolution** - Earlier layers override later (first-layer wins)
5. **Type Conversion** - Strings → target types via YAML
6. **Validation** - Required fields checked
7. **Struct Filling** - Target populated with resolved values

---

## Configuration Patterns

### Field Rules

```go
type DatabaseConfig struct {
    // Pointers for scalar types
    Host    *string `yaml:"host"`
    Port    *int    `yaml:"port"`
    Timeout *int    `yaml:"timeout"`

    // Slices and maps: non-pointer OK
    Tables  []string          `yaml:"tables"`
    Options map[string]string `yaml:"options"`
}
```

### Validation Pattern

dsco fills structs; you validate:

```go
func (c *Config) Validate() error {
    if c.Port == nil {
        return errors.New("port is required")
    }
    if *c.Port < 1 || *c.Port > 65535 {
        return errors.New("port must be 1-65535")
    }
    return nil
}

// Usage
_, err := dsco.Fill(&config, layers...)
if err != nil {
    log.Fatal(err)
}
if err := config.Validate(); err != nil {
    log.Fatal("validation:", err)
}
```

---

## Error Handling

### Error Types

| Error | Cause |
|-------|-------|
| `LayerErrors` | Layer registration issues |
| `FillerErrors` | Struct filling issues |
| `InvalidInputError` | Target not `*Config` pointer |
| `CmdlineAlreadyUsedError` | Multiple cmdline layers |
| `OverriddenKeyError` | Strict layer value overridden |

### Checking Errors

```go
_, err := dsco.Fill(&config, layers...)
if err != nil {
    var layerErr LayerErrors
    if errors.As(err, &layerErr) {
        for _, e := range layerErr.Errors() {
            log.Printf("Layer: %v", e)
        }
    }
}
```

---

## Advanced Usage

### Strict Mode

Strict layers error when values are **not consumed**:

1. Value doesn't match any field (typo detection)
2. Value overridden by an earlier layer (override detection)

```go
_, err := dsco.Fill(&config,
    dsco.WithCmdlineLayer(),            // Earlier layer — its values win
    dsco.WithStrictEnvLayer("MYAPP"),  // Strict — errors if cmdline already supplied field
)
// --port=9000 + MYAPP-PORT=8080 → Error!
// Env value was overridden by cmdline.
```

### Aliases

```go
dsco.WithCmdlineLayer(
    dsco.WithAliases(map[string]string{
        "db-host": "database.host",
        "db-port": "database.port",
        "v":       "verbose",
    }),
)
```

```bash
./myapp --db-host=localhost --v=true
# Instead of: --database-host=localhost --verbose=true
```

### File-Based Configuration

```go
type FileProvider struct {
    name   string
    values svalue.Values
}

func NewFileProvider(path string) (*FileProvider, error) {
    data, _ := os.ReadFile(path)
    var raw map[string]string
    yaml.Unmarshal(data, &raw)

    values := make(svalue.Values)
    for k, v := range raw {
        values[k] = &svalue.Value{Value: v, Location: "file:" + path}
    }
    return &FileProvider{name: path, values: values}, nil
}

func (f *FileProvider) GetName() string              { return f.name }
func (f *FileProvider) GetStringValues() svalue.Values { return f.values }
```

---

## API Reference

### Core

```go
Fill(target any, layers ...Layer) (plocation.Locations, error)
```

### Layer Builders

| Function | Description |
|----------|-------------|
| `WithCmdlineLayer(opts...)` | Command line arguments |
| `WithStrictCmdlineLayer(opts...)` | Strict command line |
| `WithEnvLayer(prefix, opts...)` | Environment variables |
| `WithStrictEnvLayer(prefix, opts...)` | Strict environment |
| `WithStructLayer(input, id)` | Struct defaults |
| `WithStrictStructLayer(input, id)` | Immutable struct values |
| `WithStringValueProvider(provider, opts...)` | Custom provider |
| `WithStrictStringValueProvider(provider, opts...)` | Strict custom provider |

### Helpers

```go
R[T any](value T) *T              // Create pointer
WithAliases(map[string]string)    // Define aliases
```

### Interfaces

```go
type StringValuesProvider interface {
    GetStringValues() svalue.Values
}

type NamedStringValuesProvider interface {
    StringValuesProvider
    GetName() string
}
```

Full API docs: [pkg.go.dev/github.com/byte4ever/dsco](https://pkg.go.dev/github.com/byte4ever/dsco)

---

## Inventory

Want to know which keys an operator must set before the service will start?
`inventory.Compute` walks your config struct and the layers you plan to
register, then reports the canonical key each layer would accept for every
leaf field. It reads nothing: no env vars, no flags, no files.

```go
import (
    "os"

    "github.com/byte4ever/dsco"
    "github.com/byte4ever/dsco/inventory"
)

var config *Config
report, err := inventory.Compute(&config,
    dsco.WithCmdlineLayer(),
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithStructLayer(defaults, "defaults"),
)
if err != nil {
    log.Fatal(err)
}
report.WriteText(os.Stdout) // or WriteJSON / WriteYAML
```

Sample text output:

```
TYPE: github.com/example/myapp.Config

PATH                  TYPE             KEY                              DEFAULT
Database.Host         *string          cmdline: --database-host=        —
Database.Port         *int             cmdline: --database-port=        defaults=5432
Server.Timeout        *time.Duration   —                                defaults=30s
```

A `—` in the DEFAULT column means no layer bakes in a value, so the operator
must supply that key. Anything with `defaults=...` is already covered.
The KEY column shows the canonical key from the first layer that can supply
the field — here cmdline, since it is listed first (highest priority).

Three runnable examples ship in the repo:

- [examples/inventory](examples/inventory/) — text dump for human eyeballing.
- [examples/inventory/json](examples/inventory/json/) — JSON output, the format
  you'd pipe into `jq` or your CI.
- [examples/inventory/preflight](examples/inventory/preflight/) — preflight
  check that exits non-zero if any key has no default, so an orchestrator
  can fail the deploy before the service even tries to start.

---

## Examples

- **[Quick Start Guide](QUICKSTART.md)** - Step-by-step tutorial
- **[examples/deadsimple](examples/deadsimple/)** - Basic multi-layer config
- **[examples/simplemain](examples/simplemain/)** - Command-line application
- **[examples/inventory](examples/inventory/)** - Inventory text dump
- **[examples/inventory/json](examples/inventory/json/)** - Inventory as JSON
- **[examples/inventory/preflight](examples/inventory/preflight/)** - Preflight check that fails the deploy when required keys are missing

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Follow project coding standards
4. Add tests
5. Run `go test -race -cover ./...`
6. Run `golangci-lint run`
7. Submit PR

```bash
go build ./...
go test -race -cover ./...
golangci-lint run
gofumpt -w .
golines --max-len=80 --base-formatter=gofumpt -w .
```

---

## License

MIT License - see [LICENSE](LICENSE)
