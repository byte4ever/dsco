# dsco Quick Start Guide

A step-by-step tutorial to master dsco configuration management.

## Table of Contents

1. [Installation](#1-installation)
2. [Core Concept: Pointer-Based Configuration](#2-core-concept-pointer-based-configuration)
3. [Your First Configuration](#3-your-first-configuration)
4. [Understanding Layers](#4-understanding-layers)
5. [Struct Layers: Default Values](#5-struct-layers-default-values)
6. [Environment Variables](#6-environment-variables)
7. [Command Line Arguments](#7-command-line-arguments)
8. [Combining Multiple Layers](#8-combining-multiple-layers)
9. [Strict Mode: Configuration Safety](#9-strict-mode-configuration-safety)
10. [Aliases: Shortcut Names](#10-aliases-shortcut-names)
11. [Custom Providers](#11-custom-providers)
12. [Error Handling](#12-error-handling)
13. [Complete Example](#13-complete-example)

---

## 1. Installation

```bash
go get github.com/byte4ever/dsco
```

Requires Go 1.21 or later.

---

## 2. Core Concept: Pointer-Based Configuration

dsco uses pointer fields to distinguish "not configured" from "configured
with a value", so a missing value isn't masked by a silent zero-value default.

### The Problem with Traditional Configuration

```go
// Plain approach: zero values are ambiguous
type Config struct {
    Port    int    // Is 0 intentional or missing?
    Host    string // Is "" intentional or missing?
    Verbose bool   // Is false intentional or missing?
}
```

### The dsco Solution

```go
// dsco approach: nil is distinct from a zero value
type Config struct {
    Port    *int    `yaml:"port"`    // nil = not configured
    Host    *string `yaml:"host"`    // nil = not configured
    Verbose *bool   `yaml:"verbose"` // nil = not configured
}
```

**Key insight**: `nil` means "not configured", any value means "explicitly set".

### The `R()` Helper

dsco provides `R[T](value T) *T` to easily create pointers:

```go
import "github.com/byte4ever/dsco"

// Instead of this:
port := 8080
config.Port = &port

// Use this:
config.Port = dsco.R(8080)
```

---

## 3. Your First Configuration

Let's create a minimal working example:

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

    _, err := dsco.Fill(
        &config,
        dsco.WithStructLayer(&Config{
            Host: dsco.R("localhost"),
            Port: dsco.R(8080),
        }, "defaults"),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Server: %s:%d\n", *config.Host, *config.Port)
}
```

Output:
```
Server: localhost:8080
```

---

## 4. Understanding Layers

dsco uses a **layered configuration system**. Each layer provides values, and
**the first layer to supply a field wins**. Later layers only fill fields left
nil by all preceding layers.

```
Layer 1 (first)  → highest priority (wins)
Layer 2          → fills fields Layer 1 left nil
Layer 3 (last)   → lowest priority
```

Think of it like stacking transparencies: you see the topmost non-transparent
layer first. Lower layers show through only where upper layers are clear (nil).

### Layer Types

| Layer Type | Source | Use Case |
|------------|--------|----------|
| `WithStructLayer` | Go struct | Default values |
| `WithEnvLayer` | Environment variables | Container/K8s config |
| `WithCmdlineLayer` | Command line args | Runtime overrides |
| `WithStringValueProvider` | Custom provider | Secrets, files, etc. |

---

## 5. Struct Layers: Default Values

Struct layers provide hardcoded values, typically used for defaults:

```go
type DatabaseConfig struct {
    Host    *string `yaml:"host"`
    Port    *int    `yaml:"port"`
    Timeout *int    `yaml:"timeout"`
}

defaults := &DatabaseConfig{
    Host:    dsco.R("localhost"),
    Port:    dsco.R(5432),
    Timeout: dsco.R(30),
}

_, err := dsco.Fill(
    &config,
    dsco.WithStructLayer(defaults, "defaults"),
)
```

The second argument (`"defaults"`) is an identifier for error messages.

### Partial Defaults

You don't need to provide all fields:

```go
// Only provide some defaults - others must come from other layers
defaults := &DatabaseConfig{
    Port:    dsco.R(5432),  // Default port
    Timeout: dsco.R(30),    // Default timeout
    // Host intentionally nil - must be provided elsewhere
}
```

---

## 6. Environment Variables

### Why Prefixes Matter

Environment variable prefixes are essential in production environments:

**1. Multi-Container Pods (Kubernetes)**

When running multiple containers in a single pod, all containers share the same
environment. Prefixes let you target configuration to specific containers:

```yaml
# Kubernetes pod with two containers
env:
  # Frontend container reads FRONTEND-* variables
  - name: FRONTEND-PORT
    value: "8080"
  - name: FRONTEND-API-URL
    value: "http://localhost:3000"

  # Backend container reads BACKEND-* variables
  - name: BACKEND-PORT
    value: "3000"
  - name: BACKEND-DATABASE-HOST
    value: "postgres.default.svc"
```

```go
// frontend/main.go
dsco.Fill(&config, dsco.WithEnvLayer("FRONTEND"))

// backend/main.go
dsco.Fill(&config, dsco.WithEnvLayer("BACKEND"))
```

**2. Avoiding Conflicts**

System and third-party tools define many environment variables (`PATH`, `HOME`,
`HTTP_PROXY`, `DATABASE_URL`, etc.). Prefixes prevent your application from
accidentally reading unrelated variables or conflicting with existing ones:

```bash
# Without prefix - risky! Could conflict with system vars
HOST=localhost        # Might conflict with other tools
PORT=8080             # Common variable name

# With prefix - safe and explicit
MYAPP-HOST=localhost  # Clearly belongs to your app
MYAPP-PORT=8080       # No ambiguity
```

**3. Multiple Instances**

Run multiple instances of the same application with different configurations:

```bash
# Instance 1
WORKER1-QUEUE=high-priority WORKER1-CONCURRENCY=10 ./worker

# Instance 2
WORKER2-QUEUE=low-priority WORKER2-CONCURRENCY=5 ./worker
```

### Choosing Good Prefixes

**Avoid generic prefixes** - they're too common and cause confusion:

```bash
# BAD: Generic, ambiguous prefixes
APP-HOST=...        # Which app? Every service is an "app"
SERVER-PORT=...     # Which server? Too vague
SERVICE-URL=...     # Meaningless in a microservices environment
CONFIG-TIMEOUT=...  # Everything has config
```

**Use specific, role-based prefixes** that identify the component:

```bash
# GOOD: Clear, distinguishable prefixes
API-HOST=...              # The API gateway
WORKER-CONCURRENCY=...    # Background job worker
CONSUMER-BATCH-SIZE=...   # Message queue consumer
SCHEDULER-INTERVAL=...    # Cron/scheduler service
GATEWAY-RATE-LIMIT=...    # API gateway
INDEXER-CHUNK-SIZE=...    # Search indexer
NOTIFIER-SMTP-HOST=...    # Notification service
```

**Why this matters:**

1. **Debugging**: When you see `WORKER-TIMEOUT=30` in logs, you instantly know
   which service it belongs to

2. **Kubernetes manifests**: Clear prefixes make YAML files self-documenting:
   ```yaml
   env:
     - name: ORDERAPI-DATABASE-HOST    # Obviously for order API
     - name: PAYMENTWORKER-RETRY-MAX   # Obviously for payment worker
   ```

3. **Shared environments**: In dev/staging where multiple services share
   infrastructure, specific prefixes prevent accidental cross-contamination

4. **Team communication**: "Check the INDEXER config" is clearer than
   "check the APP config for the indexer service"

**Naming conventions by service type:**

| Service Type | Prefix Examples |
|--------------|-----------------|
| HTTP APIs | `USERAPI`, `ORDERAPI`, `AUTHAPI` |
| Background workers | `EMAILWORKER`, `PAYMENTWORKER` |
| Message consumers | `ORDERCONSUMER`, `EVENTCONSUMER` |
| Scheduled jobs | `REPORTSCHEDULER`, `CLEANUPJOB` |
| Gateways/proxies | `APIGATEWAY`, `AUTHPROXY` |

### Format Overview

```
PREFIX-KEY=value
│      │
│      └─ Key in UPPERCASE (hyphens and underscores allowed)
└─ Prefix (UPPERCASE letters and digits only)
```

### How Parsing Works

1. **Prefix**: Must match `^[A-Z][A-Z0-9]*$` (uppercase letters/digits, starts with letter)
2. **Separator**: Single hyphen (`-`) between prefix and key
3. **Key**: Must match `^[A-Z][A-Z0-9]*(?:[-_][A-Z][A-Z0-9]*)*$`
   - Starts with uppercase letter
   - Can contain hyphens (`-`) or underscores (`_`) as word separators
   - Each word segment starts with uppercase letter

### Struct to Environment Variable Mapping

Given this struct:

```go
type Config struct {
    Host        *string         `yaml:"host"`
    Port        *int            `yaml:"port"`
    MaxRetry    *int            `yaml:"max_retry"`
    Database    *DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
    Host     *string `yaml:"host"`
    Port     *int    `yaml:"port"`
    PoolSize *int    `yaml:"pool_size"`
}
```

With `dsco.WithEnvLayer("MYAPP")`, the mapping is:

| Struct Path | YAML Key | Internal Key | Environment Variable |
|-------------|----------|--------------|---------------------|
| `Config.Host` | `host` | `host` | `MYAPP-HOST` |
| `Config.Port` | `port` | `port` | `MYAPP-PORT` |
| `Config.MaxRetry` | `max_retry` | `max_retry` | `MYAPP-MAX_RETRY` |
| `Config.Database.Host` | `database.host` | `database-host` | `MYAPP-DATABASE-HOST` |
| `Config.Database.Port` | `database.port` | `database-port` | `MYAPP-DATABASE-PORT` |
| `Config.Database.PoolSize` | `database.pool_size` | `database-pool_size` | `MYAPP-DATABASE-POOL_SIZE` |

**Key transformation rules:**
- Nested struct paths use **hyphen** (`-`) as level separator
- Field names preserve their **underscores** (`_`) from yaml tags
- Everything after prefix is **UPPERCASE** in env var, **lowercase** internally

### Valid vs Invalid Examples

```bash
# Valid environment variables
MYAPP-HOST=localhost           # Simple key
MYAPP-MAX-RETRY=5              # Hyphen in key (nested or kebab-case)
MYAPP-DB_POOL_SIZE=10          # Underscore in key (from yaml tag)
MYAPP-DATABASE-HOST=postgres   # Nested struct field

# Invalid environment variables
MYAPP_HOST=localhost           # Wrong: underscore instead of hyphen after prefix
myapp-HOST=localhost           # Wrong: lowercase prefix
MYAPP-host=localhost           # Wrong: lowercase key
MYAPP--HOST=localhost          # Wrong: double hyphen
MYAPP-123KEY=value             # Wrong: key starts with digit
```

### Usage Example

```go
type Config struct {
    Host     *string         `yaml:"host"`
    Database *DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
    Host *string `yaml:"host"`
    Port *int    `yaml:"port"`
}

var config *Config

_, err := dsco.Fill(
    &config,
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithStructLayer(&Config{
        Host: dsco.R("localhost"),
        Database: &DatabaseConfig{
            Port: dsco.R(5432),
        },
    }, "defaults"),
)
```

Run with:
```bash
MYAPP-HOST=api.example.com MYAPP-DATABASE-HOST=db.example.com MYAPP-DATABASE-PORT=5433 ./myapp
```

---

## 7. Command Line Arguments

Command line arguments use `--key=value` format:

```go
_, err := dsco.Fill(
    &config,
    dsco.WithCmdlineLayer(),
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithStructLayer(defaults, "defaults"),
)
```

Run with:
```bash
./myapp --host=production.example.com --port=9000
```

### Key Format Rules

Keys must be lowercase with hyphens or underscores:
- `--host=value` (simple key)
- `--max-connections=100` (kebab-case)
- `--db_host=localhost` (snake_case)

**Invalid**: `--Host=value` (uppercase not allowed)

### Nested Fields

For nested structs, keys are joined with **hyphens** (not dots):

```go
type Config struct {
    Database *DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
    Host *string `yaml:"host"`
    Port *int    `yaml:"port"`
}
```

```bash
# Correct: hyphen-separated
./myapp --database-host=db.example.com --database-port=5432

# Wrong: dots are NOT supported
./myapp --database.host=db.example.com  # Invalid!
```

---

## 8. Combining Multiple Layers

### Quick Local Development

Combining **struct layers** (defaults) with **command line layer** enables
rapid local development without any external configuration:

```go
var config *Config

_, err := dsco.Fill(
    &config,
    dsco.WithCmdlineLayer(),
    dsco.WithStructLayer(&Config{
        Host:     dsco.R("localhost"),
        Port:     dsco.R(8080),
        Database: &DatabaseConfig{
            Host: dsco.R("localhost"),
            Port: dsco.R(5432),
            Name: dsco.R("devdb"),
            User: dsco.R("devuser"),
        },
        LogLevel: dsco.R("debug"),
    }, "dev-defaults"),
)
```

**Benefits for local development:**

1. **Zero configuration to start**: Just run `./myapp` - all defaults are
   embedded in code, no config files or env vars needed

2. **Quick overrides**: Test different scenarios without editing files:
   ```bash
   # Test with different port
   ./myapp --port=9000

   # Test against staging database
   ./myapp --database-host=staging-db.example.com

   # Test with production-like logging
   ./myapp --log-level=info
   ```

3. **Self-documenting**: Defaults in code show what values are expected
   and what the development configuration looks like

4. **No environment pollution**: Unlike env vars, command line args don't
   persist or affect other processes

**Typical workflow:**
```bash
# Day-to-day development - just run it
./myapp

# Quick test with one change
./myapp --port=9000

# Test specific scenario
./myapp --database-host=testdb --log-level=trace

# Share exact reproduction steps with team
./myapp --feature-flag=true --timeout=5s
```

This pattern is especially useful for:
- Microservices with many configuration options
- Quick debugging sessions
- Reproducing issues with specific settings
- Onboarding new developers (no setup required)

### Full Layer Stack

The power of dsco comes from combining layers with clear precedence:

```go
type Config struct {
    Host    *string `yaml:"host"`
    Port    *int    `yaml:"port"`
    Debug   *bool   `yaml:"debug"`
    Timeout *int    `yaml:"timeout"`
}

var config *Config

_, err := dsco.Fill(
    &config,
    // Layer 1: Command line (highest priority)
    dsco.WithCmdlineLayer(),

    // Layer 2: Environment variables (middle priority)
    dsco.WithEnvLayer("MYAPP"),

    // Layer 3: Hardcoded defaults (lowest priority)
    dsco.WithStructLayer(&Config{
        Host:    dsco.R("localhost"),
        Port:    dsco.R(8080),
        Debug:   dsco.R(false),
        Timeout: dsco.R(30),
    }, "defaults"),
)
```

### Precedence Example

Given:
- Command line: `--host=production.example.com`
- Environment: `MYAPP-HOST=staging.example.com`
- Struct layer: `Host="localhost"`, `Port=8080`

Result:
- `Host` = `"production.example.com"` (from cmdline, the first layer to supply it)
- `Port` = `8080` (from struct; cmdline and env left it nil)

---

## 9. Strict Mode: Configuration Safety

**Strict mode** ensures all values from a strict layer are consumed during
filling. A value is "not consumed" if:

1. **It doesn't match any config field** (typo detection)
2. **It was overridden by an earlier layer** (override detection)

### Normal vs Strict

| Mode | Behavior |
|------|----------|
| Normal | Unmatched values and overrides silently ignored |
| Strict | Unmatched values and overrides cause an error |

### Understanding Override Detection

Remember: **the first layer to supply a field wins**. If a strict layer's value
is overridden by an earlier layer in the list, that's an error.

```go
_, err := dsco.Fill(
    &config,
    dsco.WithCmdlineLayer(),            // Earlier layer (position 0)
    dsco.WithStrictEnvLayer("MYAPP"),  // Strict layer (position 1)
)
```

If both `--port` and `MYAPP-PORT` are provided, the cmdline value wins
(earlier layer). But since the env layer is strict, its overridden value
causes an `OverriddenKeyError`.

### Typo Detection

Strict mode also catches typos - values that don't match any config field:

```bash
# Typo: HOOST instead of HOST
MYAPP-HOOST=localhost ./myapp
# Error: "HOOST" doesn't match any field, remains unused
```

### Placement Matters

Since the first layer to supply a field wins:

- **Strict layer early** → its values win; errors only for unmatched fields
- **Strict layer late** → errors if earlier layers have already supplied its values

```go
// Pattern 1: Strict cmdline at the front (typo detection only)
dsco.Fill(&config,
    dsco.WithStrictCmdlineLayer(),  // Errors only for unmatched flags
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithStructLayer(defaults, "defaults"),
)

// Pattern 2: Strict env, ensure env vars aren't overridden by cmdline
dsco.Fill(&config,
    dsco.WithCmdlineLayer(),
    dsco.WithStrictEnvLayer("MYAPP"),  // Errors if cmdline already supplied it
)

// Pattern 3: Immutable values locked at the front
dsco.Fill(&config,
    dsco.WithStrictStructLayer(&Config{
        APIEndpoint: dsco.R("https://api.production.com"),
    }, "immutable"),  // Highest priority, errors if not consumed
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithCmdlineLayer(),
)
```

### When to Use Strict Mode

**Use strict for**:
- Enforcing that certain layer values cannot be overridden
- Detecting typos in environment variable names
- Detecting unknown command line flags
- Ensuring immutable configuration values are used

**Use normal for**:
- Default values (may be skipped when earlier layers supply the field)
- Lower-priority layers where being superseded by earlier layers is expected

---

## 10. Aliases: Shortcut Names

Aliases provide short names for configuration keys:

```go
type Config struct {
    Database *DatabaseConfig `yaml:"database"`
    Server   *ServerConfig   `yaml:"server"`
    Logging  *LoggingConfig  `yaml:"logging"`
}

_, err := dsco.Fill(
    &config,
    dsco.WithCmdlineLayer(
        dsco.WithAliases(map[string]string{
            // Format: "alias": "internal.path"
            "db-host": "database.host",   // --db-host → database-host
            "db-port": "database.port",   // --db-port → database-port
            "port":    "server.port",     // --port → server-port
            "v":       "logging.verbose", // --v → logging-verbose
        }),
    ),
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithStructLayer(defaults, "defaults"),
)
```

Now you can use:
```bash
./myapp --db-host=localhost --port=9000 --v=true
```

Instead of:
```bash
./myapp --database-host=localhost --server-port=9000 --logging-verbose=true
```

**Note**: The right side of the alias mapping uses dots (internal path format),
but actual command line keys use hyphens.

---

## 11. Custom Providers

For configuration sources beyond env/cmdline/structs, implement a custom
provider:

### The Interface

```go
type NamedStringValuesProvider interface {
    GetName() string
    GetStringValues() svalue.Values
}
```

### Example: File Provider

```go
import (
    "os"

    "github.com/byte4ever/dsco/svalue"
    "gopkg.in/yaml.v3"
)

type FileProvider struct {
    name   string
    values svalue.Values
}

func NewFileProvider(path string) (*FileProvider, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var raw map[string]string
    if err := yaml.Unmarshal(data, &raw); err != nil {
        return nil, err
    }

    values := make(svalue.Values)
    for k, v := range raw {
        values[k] = &svalue.Value{
            Value:    v,
            Location: svalue.NewLocation("file", path, k),
        }
    }

    return &FileProvider{name: path, values: values}, nil
}

func (f *FileProvider) GetName() string              { return f.name }
func (f *FileProvider) GetStringValues() svalue.Values { return f.values }
```

### Usage

```go
fileProvider, err := NewFileProvider("config.yaml")
if err != nil {
    log.Fatal(err)
}

_, err = dsco.Fill(
    &config,
    dsco.WithCmdlineLayer(),                     // CLI (highest priority)
    dsco.WithEnvLayer("MYAPP"),                  // Env
    dsco.WithStringValueProvider(fileProvider),  // File config (lowest priority)
)
```

### Example: Secrets Provider

```go
type VaultProvider struct {
    client *vault.Client
}

func (v *VaultProvider) GetName() string { return "vault" }

func (v *VaultProvider) GetStringValues() svalue.Values {
    values := make(svalue.Values)

    // Fetch secrets from Vault
    secret, _ := v.client.Read("secret/myapp")

    for k, val := range secret.Data {
        values[k] = &svalue.Value{
            Value:    val.(string),
            Location: svalue.NewLocation("vault", "secret/myapp", k),
        }
    }

    return values
}
```

---

## 12. Error Handling

dsco provides detailed errors with location tracking.

### Error Types

| Error Type | Cause |
|------------|-------|
| `LayerErrors` | Layer registration problems |
| `FillerErrors` | Struct filling problems |
| `InvalidInputError` | Invalid target type |
| `CmdlineAlreadyUsedError` | Multiple cmdline layers |
| `OverriddenKeyError` | Strict value was overridden |

### Checking Errors

```go
_, err := dsco.Fill(&config, layers...)
if err != nil {
    var layerErr dsco.LayerErrors
    if errors.As(err, &layerErr) {
        for _, e := range layerErr.Errors() {
            log.Printf("Layer error: %v", e)
        }
    }

    var fillerErr dsco.FillerErrors
    if errors.As(err, &fillerErr) {
        for _, e := range fillerErr.Errors() {
            log.Printf("Fill error: %v", e)
        }
    }

    log.Fatal(err)
}
```

### Location Tracking

dsco tracks where each value came from:

```go
locations, err := dsco.Fill(&config, layers...)
if err != nil {
    log.Fatal(err)
}

// Print where each value originated
for path, loc := range locations {
    fmt.Printf("%s: %s\n", path, loc)
}
```

Output:
```
host: env[MYAPP-HOST]
port: cmdline[--port]
timeout: struct[defaults]
```

---

## 13. Complete Example

Here's a production-ready example combining all concepts:

### config.go

```go
package main

import (
    "errors"
    "time"

    "github.com/byte4ever/dsco"
)

// Config represents the application configuration.
type Config struct {
    Server   *ServerConfig   `yaml:"server"`
    Database *DatabaseConfig `yaml:"database"`
    Logging  *LoggingConfig  `yaml:"logging"`
}

type ServerConfig struct {
    Host         *string        `yaml:"host"`
    Port         *int           `yaml:"port"`
    ReadTimeout  *time.Duration `yaml:"read_timeout"`
    WriteTimeout *time.Duration `yaml:"write_timeout"`
}

type DatabaseConfig struct {
    Host     *string `yaml:"host"`
    Port     *int    `yaml:"port"`
    Name     *string `yaml:"name"`
    User     *string `yaml:"user"`
    Password *string `yaml:"password"`
    SSLMode  *string `yaml:"ssl_mode"`
}

type LoggingConfig struct {
    Level   *string `yaml:"level"`
    Format  *string `yaml:"format"`
    Verbose *bool   `yaml:"verbose"`
}

// Validate checks required fields and constraints.
func (c *Config) Validate() error {
    if c.Server == nil || c.Server.Port == nil {
        return errors.New("server.port is required")
    }
    if c.Database == nil || c.Database.Host == nil {
        return errors.New("database.host is required")
    }
    if c.Database.Password == nil {
        return errors.New("database.password is required")
    }
    return nil
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() *Config {
    return &Config{
        Server: &ServerConfig{
            Host:         dsco.R("0.0.0.0"),
            Port:         dsco.R(8080),
            ReadTimeout:  dsco.R(30 * time.Second),
            WriteTimeout: dsco.R(30 * time.Second),
        },
        Database: &DatabaseConfig{
            Port:    dsco.R(5432),
            SSLMode: dsco.R("require"),
        },
        Logging: &LoggingConfig{
            Level:   dsco.R("info"),
            Format:  dsco.R("json"),
            Verbose: dsco.R(false),
        },
    }
}
```

### main.go

```go
package main

import (
    "fmt"
    "log"

    "github.com/byte4ever/dsco"
)

func main() {
    var config *Config

    locations, err := dsco.Fill(
        &config,
        // 1. Command line (highest priority)
        dsco.WithCmdlineLayer(
            dsco.WithAliases(map[string]string{
                "db-host":     "database.host",
                "db-port":     "database.port",
                "db-name":     "database.name",
                "db-user":     "database.user",
                "db-password": "database.password",
                "port":        "server.port",
                "v":           "logging.verbose",
            }),
        ),

        // 2. Environment variables
        dsco.WithStrictEnvLayer("APP"),

        // 3. Defaults (lowest priority)
        dsco.WithStructLayer(DefaultConfig(), "defaults"),
    )
    if err != nil {
        log.Fatalf("Configuration error: %v", err)
    }

    // Validate required fields
    if err := config.Validate(); err != nil {
        log.Fatalf("Validation error: %v", err)
    }

    // Print configuration sources
    fmt.Println("Configuration loaded from:")
    for path, loc := range locations {
        fmt.Printf("  %s: %s\n", path, loc)
    }

    // Start application
    fmt.Printf("\nStarting server on %s:%d\n",
        *config.Server.Host,
        *config.Server.Port,
    )
}
```

### Running the Example

```bash
# With defaults only (will fail validation - no db password)
./myapp

# With required values
APP-DATABASE-HOST=db.example.com \
APP-DATABASE-USER=appuser \
APP-DATABASE-PASSWORD=secret123 \
APP-DATABASE-NAME=mydb \
./myapp

# Override port via command line
APP-DATABASE-HOST=db.example.com \
APP-DATABASE-USER=appuser \
APP-DATABASE-PASSWORD=secret123 \
./myapp --port=9000 --db-name=production -v=true
```

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Pointer fields | `nil` = not configured, value = explicitly set |
| `R()` helper | Creates pointers easily: `dsco.R(8080)` |
| Layer precedence | First layer to supply a field wins |
| Struct layers | Hardcoded defaults |
| Env variables | Format: `PREFIX-KEY=value` |
| Command line | Format: `--key=value` |
| Strict mode | Errors on unused values |
| Aliases | Short names for nested paths |
| Custom providers | Implement `NamedStringValuesProvider` |

**Best Practice**: Order layers from highest to lowest priority; the first
layer to supply a field wins:
```go
dsco.Fill(&config,
    dsco.WithCmdlineLayer(),        // 1. Command line (highest priority)
    dsco.WithEnvLayer(...),         // 2. Environment
    dsco.WithStringValueProvider(), // 3. Files/Secrets
    dsco.WithStructLayer(...),      // 4. Defaults (lowest priority)
)
```

---

## What keys do I need to provide?

`inventory.Compute` lists every key your layered config expects, without
reading anything from env, flags, or files:

```go
report, _ := inventory.Compute(&config, layers...)
report.WriteText(os.Stdout)
```

See the [Inventory section in the README](README.md#inventory) for the full
walkthrough and the three runnable examples (text dump, JSON for tooling, and
a preflight check that fails CI when required keys are missing).
