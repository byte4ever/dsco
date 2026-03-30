/*
Package env provides environment variable parsing functionality for dsco's
configuration system.

# Overview

The env package implements an environment variable parser that extracts
configuration values from system environment variables using a structured
prefix-based naming convention. It serves as one of the configuration layers
in dsco's layered configuration system.

# Environment Variable Naming Convention

Environment variables must follow a specific naming pattern:
PREFIX_SUBKEY=value

Where:
- PREFIX: Uppercase letters and digits, matching ^[A-Z][A-Z\d]*$
- SUBKEY: Dash-prefixed key matching ^-[A-Z][A-Z\d]*(?:[-_][A-Z][A-Z\d]*)*$
- value: Any string value

# Examples

With prefix "MYAPP":

	MYAPP_-HOST=localhost          # Maps to "host"
	MYAPP_-PORT=8080              # Maps to "port"
	MYAPP_-MAX-CONNECTIONS=100    # Maps to "max-connections"
	MYAPP_-DB_HOST=postgres       # Maps to "db_host"
	MYAPP_-API-KEY-V2=secret      # Maps to "api-key-v2"

The resulting configuration keys are lowercase versions of the SUBKEY (without
the leading dash):
- MYAPP_-HOST → "host"
- MYAPP_-MAX-CONNECTIONS → "max-connections"
- MYAPP_-DB_HOST → "db_host"

# Prefix Validation

Prefixes must be valid uppercase identifiers:

Valid prefixes:

	MYAPP, APP, SERVICE1, API2, CONFIGV3

Invalid prefixes:

	myapp         # Lowercase not allowed
	123APP        # Cannot start with digit
	MY-APP        # Hyphens not allowed in prefix
	MY_APP        # Underscores not allowed in prefix

# Key Format Rules

Environment variable keys (after prefix) must follow specific patterns:

## Valid Key Examples

	MYAPP_-HOST=localhost              # Simple key
	MYAPP_-MAX-RETRY-COUNT=5          # Kebab-case with hyphens
	MYAPP_-DB_CONNECTION_POOL=10      # Snake_case with underscores
	MYAPP_-API-VERSION-2=v2.1         # Mixed alphanumeric

## Invalid Key Examples

	MYAPP_HOST=value                  # Missing dash prefix
	MYAPP_-host=value                 # Lowercase not allowed
	MYAPP_-123=value                  # Cannot start with digit
	MYAPP_-=value                     # Empty key after dash

# Value Format

Values can contain any characters including spaces, newlines, and special
characters. The entire string after the first '=' is treated as the value:

	MYAPP_-MESSAGE="Hello, World!"
	MYAPP_-JSON_CONFIG={"key": "value", "array": [1,2,3]}
	MYAPP_-MULTILINE="Line 1\nLine 2\nLine 3"
	MYAPP_-PATH="/home/user/config files/app.yaml"

# Error Handling

The package provides detailed error handling for various issues:

## Invalid Prefix Errors

When the prefix doesn't match the required format:

	_, err := env.NewEntriesProvider("invalid-prefix")
	// Returns: ErrInvalidPrefix

## Ambiguous Key Errors

When environment variable names don't follow the expected pattern:

	// These would be ambiguous:
	MYAPP_INVALID_KEY=value1    # Missing dash, doesn't match pattern
	MYAPP_-invalid=value2       # Lowercase not allowed

	// Results in: ErrAmbiguousKey or ErrAmbiguousKeys

# Usage Examples

## Basic Usage

	package main

	import (
		"fmt"
		"os"
		"github.com/byte4ever/dsco/internal/env"
	)

	func main() {
		// Set some environment variables
		os.Setenv("MYAPP_-HOST", "localhost")
		os.Setenv("MYAPP_-PORT", "8080")
		os.Setenv("MYAPP_-VERBOSE", "true")

		// Parse environment variables with prefix
		provider, err := env.NewEntriesProvider("MYAPP")
		if err != nil {
			fmt.Printf("Error parsing environment: %v\n", err)
			return
		}

		// Get parsed values
		values := provider.GetStringValues()
		for key, value := range values {
			fmt.Printf("%s = %s (from %s)\n", key, value.Value, value.Location)
		}
		// Output:
		// host = localhost (from env[MYAPP_-HOST])
		// port = 8080 (from env[MYAPP_-PORT])
		// verbose = true (from env[MYAPP_-VERBOSE])
	}

## Integration with dsco

	import "github.com/byte4ever/dsco"

	type Config struct {
		Host    *string `yaml:"host"`
		Port    *int    `yaml:"port"`
		Verbose *bool   `yaml:"verbose"`
	}

	func main() {
		var config *Config

		// Environment layer will automatically use this package
		_, err := dsco.Fill(
			&config,
			dsco.WithEnvLayer("MYAPP"),  // Uses internal/env
		)

		if err != nil {
			log.Fatal(err)
		}
	}

## Multiple Prefixes

Different prefixes can be used for different configuration sections:

	_, err := dsco.Fill(
		&config,
		dsco.WithEnvLayer("DATABASE"),  # DATABASE_-HOST, DATABASE_-PORT
		dsco.WithEnvLayer("API"),       # API_-KEY, API_-VERSION
		dsco.WithEnvLayer("CACHE"),     # CACHE_-TTL, CACHE_-SIZE
	)

# Location Tracking

Each parsed value includes location information for debugging:

	Location format: "env[PREFIX_SUBKEY]"

Example locations:
  - "env[MYAPP_-HOST]"        # For MYAPP_-HOST=localhost
  - "env[API_-KEY]"           # For API_-KEY=secret
  - "env[DB_-CONNECTION]"     # For DB_-CONNECTION=postgres://...

This location information is used by dsco for error reporting and debugging,
helping users identify exactly where configuration values originated.

# Key Transformation

Environment variable names are transformed to configuration keys:

	MYAPP_-HOST           → "host"
	MYAPP_-MAX-RETRY      → "max-retry"
	MYAPP_-DB_POOL_SIZE   → "db_pool_size"
	MYAPP_-API-V2         → "api-v2"

The transformation process:
1. Remove prefix and separator
2. Remove leading dash
3. Convert to lowercase
4. Preserve hyphens and underscores

# Memory Management

The package is designed for efficient memory usage:

- Pre-allocates maps based on environment variable count
- Efficient regular expression matching with compiled patterns
- Minimal memory overhead per parsed variable
- No persistent state between parsing operations
- Sorts environment variables once for consistent processing

# Thread Safety

The package is thread-safe for concurrent use:

- NewEntriesProvider can be called concurrently from multiple goroutines
- EntriesProvider instances are safe for concurrent read access
- No shared mutable state between provider instances
- Regular expressions are compiled once and reused safely

# Performance Characteristics

The parser is optimized for typical environment usage:

- O(n log n) time complexity due to environment variable sorting
- Efficient regular expression matching with pre-compiled patterns
- Linear scan through environment variables
- Fast prefix matching and key transformation

# Integration Points

This package integrates with other dsco components:

## svalue Package

Uses svalue.Values for consistent value representation:

	type Values map[string]*Value

	type Value struct {
		Location string  // Source location for debugging
		Value    string  // Actual configuration value
	}

## utils Package

Uses utility functions for string formatting:

	utils.FormatStringSequence(ambiguousKeys)  // For error messages

## Error System

Integrates with dsco's error handling through structured error types:

	var (
		ErrInvalidPrefix  = errors.New("invalid prefix")
		ErrAmbiguousKey   = errors.New("is ambiguous")
		ErrAmbiguousKeys  = errors.New("are ambiguous")
	)

## Layer System

Provides the foundation for environment variable configuration layers:

	// Used by WithEnvLayer(prefix)
	// Used by WithStrictEnvLayer(prefix)

# Error Messages

The package provides clear, actionable error messages:

	"INVALID-prefix" : invalid prefix
	"MYAPP_INVALID" is ambiguous
	"MYAPP_BAD1", "MYAPP_BAD2" are ambiguous

Error messages include:
- Exact problematic environment variable names
- Clear indication of the specific problem (invalid prefix, ambiguous format)
- Proper pluralization for multiple issues
- Consistent formatting across all error types

# Environment Variable Processing

The package processes environment variables in a specific order:

1. **Filtering**: Only variables with matching prefix are considered
2. **Sorting**: Variables are sorted alphabetically for consistent processing
3. **Validation**: Each variable is checked against the expected pattern
4. **Transformation**: Valid variables are transformed to configuration keys
5. **Collection**: Ambiguous variables are collected for error reporting

# Best Practices

## Naming Conventions

Use consistent naming patterns:

	# Good: Clear hierarchy
	MYAPP_-DATABASE-HOST=localhost
	MYAPP_-DATABASE-PORT=5432
	MYAPP_-DATABASE-NAME=mydb

	# Good: Logical grouping
	CACHE_-TTL=3600
	CACHE_-SIZE=1000
	CACHE_-ENABLED=true

## Testing

Test with various environment scenarios:

	func TestEnvVars(t *testing.T) {
		os.Setenv("TEST_-HOST", "localhost")
		defer os.Unsetenv("TEST_-HOST")

		provider, err := env.NewEntriesProvider("TEST")
		require.NoError(t, err)

		values := provider.GetStringValues()
		assert.Equal(t, "localhost", values["host"].Value)
		assert.Equal(t, "env[TEST_-HOST]", values["host"].Location)
	}

# Security Considerations

Environment variables are visible to all processes, so:

- Avoid storing sensitive data in environment variables when possible
- Use secure methods for secret management in production
- Be aware that environment variables appear in process listings
- Consider using file-based or secret management systems for sensitive configuration

The package itself does not provide encryption or security features beyond
standard environment variable handling.
*/
package env
