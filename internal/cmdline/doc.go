/*
Package cmdline provides command-line argument parsing functionality for dsco's
configuration system.

# Overview

The cmdline package implements a command-line argument parser that extracts
configuration values from command-line arguments following a strict format.
It serves as one of the configuration layers in dsco's layered configuration
system.

# Command Line Format

Command line arguments must follow the format: --key=value

	--host=localhost
	--port=8080
	--timeout=30s
	--verbose=true

# Key Format Rules

Keys must match the regular expression: ^--([a-z][a-z\d]*(?:[-_][a-z][a-z\d]*)*)=(.+)$

Valid key examples:

	--host=localhost          # Simple lowercase key
	--max-connections=100     # Kebab-case with hyphens
	--db_host=localhost       # Snake_case with underscores
	--api-key-v2=secret       # Mixed alphanumeric

Invalid key examples:

	--Host=value              # Uppercase letters not allowed
	--123key=value            # Cannot start with digits
	--=value                  # Empty key
	host=value                # Missing -- prefix

# Value Format

Values can contain any characters including spaces, special characters, and
unicode. The entire string after the first '=' is treated as the value:

	--message="Hello, World!"
	--path="/home/user/config file.yaml"
	--json={"key": "value", "nested": {"array": [1,2,3]}}

# Error Handling

The package provides comprehensive error handling for common issues:

## Invalid Format Errors

When arguments don't match the required format:

	args := []string{"invalid-arg"}
	_, err := cmdline.NewEntriesProvider(args)
	// Returns: ErrInvalidFormat

## Duplicate Parameter Errors

When the same parameter is specified multiple times:

	args := []string{"--host=first", "--host=second"}
	_, err := cmdline.NewEntriesProvider(args)
	// Returns: ErrDuplicateParam with position information

## Error Details

The ParamError type provides detailed information about parsing failures:

	type ParamError struct {
		Positions []int    // 1-based positions of problematic arguments
		Errs      []error  // Corresponding errors for each position
	}

Example error message:

	"cmdline issue at position #3: --host previous found at position #1: duplicate param"

# Usage Examples

## Basic Usage

	package main

	import (
		"fmt"
		"os"
		"github.com/byte4ever/dsco/internal/cmdline"
	)

	func main() {
		// Parse command line arguments (excluding program name)
		provider, err := cmdline.NewEntriesProvider(os.Args[1:])
		if err != nil {
			fmt.Printf("Error parsing command line: %v\n", err)
			return
		}

		// Get parsed values
		values := provider.GetStringValues()
		for key, value := range values {
			fmt.Printf("%s = %s (from %s)\n", key, value.Value, value.Location)
		}
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

		// Command line layer will automatically use this package
		_, err := dsco.Fill(
			&config,
			dsco.WithCmdlineLayer(),  // Uses internal/cmdline
		)

		if err != nil {
			log.Fatal(err)
		}
	}

# Location Tracking

Each parsed value includes location information for debugging:

	Location format: "cmdline[--key]"

Example locations:
  - "cmdline[--host]"      # For --host=localhost
  - "cmdline[--port]"      # For --port=8080
  - "cmdline[--timeout]"   # For --timeout=30s

This location information is used by dsco for error reporting and debugging,
helping users identify exactly where configuration values originated.

# Memory Management

The package is designed for efficient memory usage:

- Pre-allocates maps and slices based on argument count
- Reuses regular expression compilation via package-level variables
- Minimal memory overhead per parsed argument
- No memory leaks or persistent state between parsing operations

# Thread Safety

The package is thread-safe for concurrent use:

- NewEntriesProvider can be called concurrently from multiple goroutines
- EntriesProvider instances are safe for concurrent read access
- No shared mutable state between provider instances
- Regular expressions are compiled once and reused safely

# Performance Characteristics

The parser is optimized for typical command-line usage:

- O(n) parsing time where n is the number of arguments
- Efficient regular expression matching with pre-compiled patterns
- Minimal memory allocation during parsing
- Fast duplicate detection using map lookups

# Integration Points

This package integrates with other dsco components:

## svalue Package

Uses svalue.Values for consistent value representation across all layers:

	type Values map[string]*Value

	type Value struct {
		Location string  // Source location for debugging
		Value    string  // Actual configuration value
	}

## Error System

Integrates with dsco's error handling through structured error types:

	// ParamError implements error interface
	func (e *ParamError) Error() string

## Layer System

Provides the foundation for command-line configuration layers in dsco:

	// Used by WithCmdlineLayer()
	// Used by WithStrictCmdlineLayer()

# Testing

The package includes comprehensive tests covering:

- Valid argument parsing for various key formats
- Invalid argument detection and error reporting
- Duplicate parameter detection with position tracking
- Edge cases (empty arguments, malformed inputs)
- Error message format verification
- Location string generation
- Memory usage and performance benchmarks

All tests follow dsco's testing standards with parallel execution and
comprehensive error checking using testify framework.

# Error Messages

The package provides clear, actionable error messages:

	"arg "some-thing": invalid format"
	"cmdline issue at position #3: --arg1 previous found at position #1: duplicate param"

Error messages include:
- Exact problematic argument value
- 1-based position numbers for easy identification
- Clear indication of the specific problem
- Consistent formatting across all error types
*/
package cmdline
