/*
Package svalue provides string value types with location tracking for dsco
configuration sources.

# Overview

The svalue package defines fundamental data structures for representing
string-based configuration values along with their source locations. This
enables precise error reporting and configuration debugging by tracking
where each value originated.

# Core Types

The package provides two primary types:

## Value

Value represents a single string configuration value with its location:

	type Value struct {
		Location string  // Where this value came from
		Value    string  // The actual configuration value
	}

The Location field typically contains source information like:
- "env:MYAPP_HOST" for environment variables
- "cmdline:--port" for command line arguments
- "file:config.yaml:line:15" for file-based sources
- "struct[defaults]:database.host" for struct-based sources

## Values

Values represents a collection of string values keyed by field paths:

	type Values map[string]*Value

Keys are typically dot-separated field paths like:
- "host" for top-level fields
- "database.host" for nested struct fields
- "items[0].name" for array element fields
- "logging.level" for deeply nested configuration

# Usage in Configuration Sources

String value providers implement interfaces that return svalue.Values:

	type CustomProvider struct{}

	func (p *CustomProvider) GetStringValues() svalue.Values {
		return svalue.Values{
			"database.host": &svalue.Value{
				Value:    "production.db.internal",
				Location: "vault:database-credentials:host",
			},
			"database.port": &svalue.Value{
				Value:    "5432",
				Location: "vault:database-credentials:port",
			},
		}
	}

# Location Tracking

Location information is crucial for debugging configuration issues:

	// When an error occurs, you can trace the source
	_, err := dsco.Fill(&config, layers...)
	if err != nil {
		// Error messages include location information:
		// "invalid port value '80bogus' from cmdline:--port"
		// "missing required field 'host' (checked env:MYAPP_HOST,
		// file:config.yaml)"
	}

# Integration with dsco

The svalue types are primarily used internally by dsco's layer system:

1. String-based sources (cmdline, env, files) convert their values to svalue.Values
2. The layer system processes these values during configuration resolution
3. Location information is preserved for error reporting and debugging
4. Type conversion transforms string values to target struct field types

# Value Lifecycle

The typical lifecycle of svalue.Values in dsco:

1. **Collection**: Configuration sources provide svalue.Values
2. **Aggregation**: Multiple sources contribute to the same Values map
3. **Resolution**: Later layers override earlier ones (by key)
4. **Conversion**: String values converted to target types via YAML unmarshaling
5. **Validation**: Location information included in any error messages
6. **Cleanup**: Used values are removed from the map in strict mode

# Testing Coverage

This package maintains 100% test coverage, including:
- Value creation and field access
- Values map operations (get, set, delete, iterate)
- Complex key structures with nested paths
- Location string handling with special characters
- Memory management and pointer safety
- Edge cases with nil values and empty collections

The comprehensive test suite ensures reliable operation across all
configuration scenarios and error conditions.

# Thread Safety

The svalue types are not inherently thread-safe. Concurrent access must
be managed by the calling code. In dsco's architecture, this is handled
by the layer processing pipeline, which operates sequentially during
configuration resolution.

# Memory Considerations

Values in the map are stored as pointers, allowing for efficient memory
usage when the same value appears multiple times. The Location strings
provide debugging information but do consume memory proportional to their
length and specificity.

For optimal memory usage in large configuration scenarios, consider:
- Using shorter, standardized location formats
- Reusing common location strings where possible
- Cleaning up Values maps after processing when appropriate.
*/
package svalue
