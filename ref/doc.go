/*
Package ref provides utility functions for creating references to values.

# Overview

The ref package contains helper functions to create pointers to values,
which is essential for dsco's pointer-based configuration system. This
eliminates the need for intermediate variables when creating configuration
structures with explicit pointer fields.

# Primary Function

The main utility function is R, which returns a pointer to any value:

	// Instead of this verbose approach:
	host := "localhost"
	port := 8080
	config := &Config{
		Host: &host,
		Port: &port,
	}

	// Use the convenient R function:
	config := &Config{
		Host: ref.R("localhost"),
		Port: ref.R(8080),
	}

# Generic Support

The R function uses Go generics to work with any type, providing type safety
while maintaining convenience:

	// Works with all basic types
	stringPtr := ref.R("hello")        // *string
	intPtr := ref.R(42)                // *int
	boolPtr := ref.R(true)             // *bool
	durationPtr := ref.R(time.Second)  // *time.Duration

	// Works with complex types
	slicePtr := ref.R([]string{"a", "b"})    // *[]string
	mapPtr := ref.R(map[string]int{"x": 1})  // *map[string]int

	// Works with custom structs
	type Custom struct{ Value int }
	customPtr := ref.R(Custom{Value: 123})   // *Custom

# Zero Values

The function correctly handles zero values, creating pointers to them:

	zeroInt := ref.R(0)          // *int pointing to 0
	emptyString := ref.R("")     // *string pointing to ""
	falseBool := ref.R(false)    // *bool pointing to false

This is important for dsco's explicit configuration model, where zero values
must be distinguishable from nil (unconfigured) values.

# Testing Coverage

This package maintains 100% test coverage, including:
- All Go basic types (string, int, bool, float64)
- Zero value handling for each type
- Complex types (structs, slices, maps)
- Generic type parameter validation
- Memory safety verification

The tests ensure the function works correctly across all supported Go types
and maintains the expected pointer relationships.

# Usage in dsco

This package is primarily used within dsco configuration structs:

	type ServiceConfig struct {
		Database *DatabaseConfig `yaml:"database"`
		Server   *ServerConfig   `yaml:"server"`
	}

	type DatabaseConfig struct {
		Host     *string `yaml:"host"`
		Port     *int    `yaml:"port"`
		Timeout  *int    `yaml:"timeout"`
	}

	// Clean configuration creation
	defaults := &ServiceConfig{
		Database: &DatabaseConfig{
			Host:    ref.R("localhost"),
			Port:    ref.R(5432),
			Timeout: ref.R(30),
		},
		Server: &ServerConfig{
			Port:    ref.R(8080),
			Timeout: ref.R(60),
		},
	}

The ref.R function is aliased as dsco.R for convenience in the main package.
*/
package ref
