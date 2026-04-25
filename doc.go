/*
Package dsco (pronounce /ˈdɪskoʊ/) provides a layered configuration system
supporting command line arguments, environment variables, YAML files, and
struct-based configurations with strict validation.

# Overview

dsco implements a layered configuration system where different configuration
sources are organized into layers with configurable precedence. Earlier layers
override later ones (first-layer wins); a layer that leaves a field nil falls
through to the next layer. Strict mode detects conflicts and unused values.

The library is designed for microservices environments where configuration
safety is critical. It enforces explicit configuration through pointer-based
fields, ensuring no silent defaults or ambiguous states.

# Quick Start

	type Config struct {
		Host     *string        `yaml:"host"`
		Port     *int           `yaml:"port"`
		Timeout  *time.Duration `yaml:"timeout"`
		Verbose  *bool          `yaml:"verbose"`
	}

	var config *Config
	_, err := dsco.Fill(
		&config,
		dsco.WithCmdlineLayer(),
		dsco.WithEnvLayer("MYAPP"),
		dsco.WithStructLayer(&Config{
			Host:    dsco.R("localhost"),
			Port:    dsco.R(8080),
			Timeout: dsco.R(30 * time.Second),
			Verbose: dsco.R(false),
		}, "defaults"),
	)

# Layer Types

The library supports multiple configuration sources:

- Command line arguments (WithCmdlineLayer, WithStrictCmdlineLayer)
- Environment variables (WithEnvLayer, WithStrictEnvLayer)
- Go structs (WithStructLayer, WithStrictStructLayer)
- Custom string providers (WithStringValueProvider, WithStrictStringValueProvider)
- File-based providers (via external packages like kfile)

# Safety Design

dsco enforces explicit configuration through pointer-based fields:

- nil clearly means "not configured"
- Non-nil values mean "explicitly configured"
- Zero values (0, "", false) don't mask missing configuration
- Services fail fast with clear error messages
- No hidden defaults or silent failures

# Strict Mode

Strict mode layers detect unused configuration values and conflicts:

	_, err := dsco.Fill(
		&config,
		dsco.WithStrictCmdlineLayer(),
		dsco.WithStrictEnvLayer("MYAPP"),
	)

This ensures all provided configuration values are actually used, preventing
configuration drift and accidental misconfigurations.

# Error Handling

The library provides comprehensive error types with location tracking:

	if err != nil {
		var layerErr LayerErrors
		if errors.As(err, &layerErr) {
			// Handle layer registration errors
		}

		var fillerErr FillerErrors
		if errors.As(err, &fillerErr) {
			// Handle field filling errors
		}
	}

# Testing

The dsco library maintains 100% test coverage across all core packages,
ensuring reliability and stability for production use. The test suite includes:

- Unit tests for all public APIs
- Integration tests for multi-layer scenarios
- Error path testing with comprehensive validation
- Concurrent access testing for thread safety
- Edge case handling with boundary conditions

The testing infrastructure follows project standards:
- Parallel test execution with t.Parallel()
- Mock-based testing using testify framework
- Table-driven test patterns for comprehensive coverage
- Error assertion using testify.ErrorIs, ErrorAs, and Contains

# Architecture

The library implements a sophisticated processing pipeline:

1. Layer Registration: Different configuration sources register as layers
2. Model Generation: Target struct is analyzed via reflection
3. Value Collection: Each layer provides field values from its source
4. Precedence Resolution: Earlier layers override later ones (first-layer wins)
5. Type Conversion: String values converted to target types via YAML
6. Validation: Required fields and custom validation applied
7. Struct Filling: Target struct populated with resolved values

For complete documentation and examples, see:
https://pkg.go.dev/github.com/byte4ever/dsco

# Inventory

The inventory sub-package (github.com/byte4ever/dsco/inventory) computes
a static list of configuration keys a Fill call would expect, with no
I/O, suitable for "what do I need to set" diagnostics in operator
tooling. See the Inventory section in README.md for examples.
*/
package dsco
