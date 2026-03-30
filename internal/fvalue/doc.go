/*
Package fvalue provides field value representation for dsco's configuration
processing system.

# Overview

The fvalue package defines the core data structures used throughout dsco for
representing configuration field values with their associated metadata. It
serves as the fundamental value container that flows through dsco's processing
pipeline, carrying both the actual configuration value and essential debugging
information.

# Core Types

## Value

The Value struct represents a single configuration field value with its metadata:

	type Value struct {
		Value    reflect.Value  // The actual configuration value
		Location string         // Source location for debugging
		Path     string         // Configuration path/key
	}

### Fields

  - **Value**: A reflect.Value containing the actual configuration data. This could be
    any Go type (string, int, bool, struct, slice, etc.) as determined by the
    target configuration structure.

  - **Location**: A human-readable string indicating where this value originated
    (e.g., "cmdline[--host]", "env[MYAPP_-PORT]", "struct[defaults]:timeout").

  - **Path**: The configuration path/key identifying this field within the
    configuration hierarchy (e.g., "host", "database.connection.timeout").

## Values

Values is a map type that associates field identifiers with their values:

	type Values map[uint]*Value

The key is a unique identifier (UID) generated during model scanning that
corresponds to a specific field in the target configuration structure.

# Usage in dsco Pipeline

The fvalue types are used throughout dsco's processing pipeline:

## Layer Processing

Each configuration layer (cmdline, env, struct, etc.) produces fvalue.Values:

	// Command line layer produces values like:
	values := fvalue.Values{
		1: &fvalue.Value{
			Value:    reflect.ValueOf("localhost"),
			Location: "cmdline[--host]",
			Path:     "host",
		},
		2: &fvalue.Value{
			Value:    reflect.ValueOf(8080),
			Location: "cmdline[--port]",
			Path:     "port",
		},
	}

## Model Building

During model construction, fields are assigned unique IDs that serve as map keys:

	// Model assigns UIDs during struct analysis:
	type Model struct {
		fieldMappings map[string]uint  // path -> UID
		// ...
	}

## Value Resolution

Multiple layers provide values that are resolved by precedence:

	// Later layers override earlier ones:
	layer1Values := fvalue.Values{1: &fvalue.Value{...}}  // defaults
	layer2Values := fvalue.Values{1: &fvalue.Value{...}}  // env vars
	layer3Values := fvalue.Values{1: &fvalue.Value{...}}  // cmdline (highest precedence)

## Error Reporting

Location information enables precise error reporting:

	func reportError(fieldUID uint, allLayers []fvalue.Values) {
		for _, layer := range allLayers {
			if value := layer[fieldUID]; value != nil {
				fmt.Printf("Error in field '%s' from %s\n", value.Path, value.Location)
				break
			}
		}
	}

# Value Creation Patterns

## Primitive Values

Creating values for basic Go types:

	// String value
	stringValue := &fvalue.Value{
		Value:    reflect.ValueOf("localhost"),
		Location: "cmdline[--host]",
		Path:     "host",
	}

	// Integer value
	intValue := &fvalue.Value{
		Value:    reflect.ValueOf(8080),
		Location: "env[MYAPP_-PORT]",
		Path:     "port",
	}

	// Boolean value
	boolValue := &fvalue.Value{
		Value:    reflect.ValueOf(true),
		Location: "struct[defaults]:verbose",
		Path:     "verbose",
	}

## Complex Values

Creating values for complex types:

	// Slice value
	sliceValue := &fvalue.Value{
		Value:    reflect.ValueOf([]string{"host1", "host2"}),
		Location: "file[config.yaml]:servers",
		Path:     "servers",
	}

	// Struct value (as pointer)
	structPtr := &MyStruct{Field: "value"}
	structValue := &fvalue.Value{
		Value:    reflect.ValueOf(structPtr),
		Location: "struct[defaults]:nested",
		Path:     "nested",
	}

# Location Conventions

Different configuration sources use consistent location formats:

## Command Line

	"cmdline[--key]"
	"cmdline[--max-connections]"

## Environment Variables

	"env[PREFIX_-KEY]"
	"env[MYAPP_-HOST]"
	"env[DATABASE_-CONNECTION-POOL]"

## Struct Sources

	"struct[id]:path"
	"struct[defaults]:host"
	"struct[production]:database.timeout"

## File Sources

	"file[config.yaml]:path"
	"file[/etc/app/settings.json]:api.key"

# Path Conventions

Configuration paths use dot notation for nested fields:

	"host"                    # Simple field
	"database.host"           # Nested field
	"database.pool.size"      # Deeply nested field
	"servers[0].address"      # Array element (theoretical)

# Memory Management

The fvalue package is designed for efficient memory usage:

## Reflect Value Handling

reflect.Value instances are managed carefully:
- Values are created once per configuration field
- No unnecessary copying of reflect.Value instances
- Values are passed by reference to avoid duplication

## String Allocation

Location and path strings are allocated efficiently:
- Location strings are formatted once when values are created
- Path strings are typically shared across the model
- Minimal string concatenation during normal operation

## Map Efficiency

The Values map type is optimized for typical usage:
- Pre-allocated based on known field count from model analysis
- Efficient uint key lookups (faster than string keys)
- Sparse map usage (only populated fields consume memory)

# Thread Safety

The fvalue types have specific thread safety characteristics:

## Value Struct
- **Immutable after creation**: Once a Value is created, it should not be modified
- **Safe for concurrent reads**: Multiple goroutines can read from a Value safely
- **reflect.Value safety**: The contained reflect.Value follows Go's reflection safety rules

## Values Map
- **Not thread-safe for concurrent writes**: Concurrent modification requires external synchronization
- **Safe for concurrent reads after population**: Read-only access is safe from multiple goroutines
- **Layer isolation**: Each layer creates its own Values map, avoiding conflicts

# Integration with Other Packages

The fvalue package integrates closely with other dsco components:

## reflect Package

	Value.Value field contains reflect.Value for type-safe value handling

## model Package

	Models define field UIDs that serve as Values map keys

## Layer Packages (cmdline, env, etc.)

	Each layer package produces fvalue.Values as output

## Error Reporting

	Location and path information enables detailed error messages

# Debugging Support

The fvalue package provides excellent debugging support:

## Value Inspection

	func debugValue(v *fvalue.Value) {
		fmt.Printf("Path: %s\n", v.Path)
		fmt.Printf("Location: %s\n", v.Location)
		fmt.Printf("Type: %s\n", v.Value.Type())
		fmt.Printf("Value: %v\n", v.Value.Interface())
	}

## Layer Analysis

	func debugLayer(values fvalue.Values) {
		for uid, value := range values {
			fmt.Printf("UID %d: %s = %v (from %s)\n",
				uid, value.Path, value.Value.Interface(), value.Location)
		}
	}

# Best Practices

## Value Creation
- Always set all three fields (Value, Location, Path) when creating Value instances
- Use consistent location formatting for the same source type
- Ensure reflect.Value is valid and represents the correct type

## Values Map Usage
- Pre-allocate Maps with known capacity when possible
- Use uint UIDs consistently as generated by the model package
- Don't modify Values maps after passing to other components

## Error Handling
- Include location information in all error messages
- Use path information to identify problematic configuration fields
- Preserve original reflect.Value for accurate type information

# Performance Characteristics

The package is optimized for dsco's usage patterns:

## Memory Usage
- Minimal overhead per Value (3 fields: reflect.Value + 2 strings)
- Efficient map structure with uint keys
- No hidden allocations or object pools

## Access Patterns
- O(1) lookup time for Values map access by UID
- Fast iteration over Values map for layer processing
- Minimal CPU overhead for value creation and access

## Scalability
- Scales linearly with number of configuration fields
- No performance degradation with deep nesting
- Efficient handling of large configuration structures

This package forms the backbone of dsco's value representation system,
providing the essential types and conventions that enable reliable, debuggable
configuration processing across all layers and components.
*/
package fvalue
