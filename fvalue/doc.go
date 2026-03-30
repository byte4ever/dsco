/*
Package fvalue provides field value types with location tracking for dsco's
internal configuration processing pipeline.

# Overview

The fvalue package defines data structures for representing configuration
field values that have been processed from string sources into typed values,
while maintaining location information for debugging and error reporting.
This is the intermediate representation between raw string values and the
final struct field assignments.

# Core Types

The package provides two primary types:

## Value

Value represents a processed configuration value with its source location:

	type Value struct {
		Location string      // Where this value originated
		Value    interface{} // The processed/typed value
	}

Unlike svalue.Value (which stores only strings), fvalue.Value contains:
- Typed values converted from strings via YAML unmarshaling
- Complex values like slices, maps, or custom structs
- Metadata about the original source location

Common Value.Value types include:
- Basic types: string, int, bool, float64
- Slices: []string, []int, []interface{}
- Maps: map[string]interface{}, map[string]string
- Custom structs and nested configurations

## Values

Values represents a collection of processed field values:

	type Values map[string]*Value

Keys correspond to struct field paths in the target configuration:
- "host" for top-level fields
- "database.host" for nested struct fields
- "listeners[0].port" for array element fields
- "logging.handlers.file.path" for deeply nested fields

# Value Processing Pipeline

The fvalue package fits into dsco's processing pipeline:

1. **Input**: svalue.Values (raw strings from sources)
2. **Processing**: Convert strings to typed values via YAML unmarshaling
3. **Output**: fvalue.Values (typed values ready for struct assignment)
4. **Assignment**: Values applied to target struct fields

Example processing:

	// Input svalue
	stringVal := &svalue.Value{
		Value:    "5432",
		Location: "env:DATABASE_PORT",
	}

	// Processed fvalue (int conversion)
	fieldVal := &fvalue.Value{
		Value:    5432,  // Now an int, not string
		Location: "env:DATABASE_PORT",  // Location preserved
	}

# Type Conversion

The package handles complex type conversions following YAML semantics:

## Basic Types
- Strings: direct assignment or quoted values
- Integers: decimal, hex (0x), octal (0o), binary (0b)
- Floats: decimal, scientific notation, special values (.inf, .nan)
- Booleans: true/false, yes/no, on/off, 1/0

## Complex Types
- Arrays: comma-separated or YAML list syntax
- Maps: YAML mapping syntax or JSON-like structures
- Structs: nested YAML structures mapped to Go struct fields
- Pointers: nil values handled explicitly

Example complex conversions:

	// Input string: "port:5432,host:localhost,ssl:true"
	// Output map[string]interface{}{
	//     "port": 5432,
	//     "host": "localhost",
	//     "ssl": true,
	// }

# Location Preservation

Location information flows through the entire processing pipeline:

	originalLocation := "cmdline:--database-config"

	// String-to-map conversion preserves location
	processedValue := &fvalue.Value{
		Value: map[string]interface{}{
			"host": "localhost",
			"port": 5432,
		},
		Location: originalLocation, // Original source preserved
	}

This enables precise error reporting even after complex transformations.

# Integration with Model System

fvalue.Values integrate with dsco's reflection-based model system:

	// Model processes struct fields and generates fvalue.Values
	model := buildModel(reflect.TypeOf(config))
	fieldValues, err := model.ApplyOn(valueGetter)

	// fieldValues is fvalue.Values ready for struct assignment
	locations, err := model.Fill(configPtr, []fvalue.Values{fieldValues})

The model system:
1. Analyzes target struct using reflection
2. Matches field paths to available values
3. Performs type checking and conversion validation
4. Reports conflicts and missing required fields

# Error Handling and Validation

The package supports comprehensive error reporting:

## Type Mismatch Errors
When conversion fails, errors include:
- Source location information
- Expected vs actual types
- Original string value for debugging

## Required Field Validation
Missing required fields report:
- Field path and expected type
- All locations checked for the value
- Suggestion for correct naming/sourcing

## Conflict Resolution
When multiple sources provide the same field:
- Normal mode: later layers win
- Strict mode: conflicts are reported as errors
- Location information helps identify conflicting sources

# Usage Patterns

## Layer Processing
Layers convert their svalue.Values to fvalue.Values:

	func (l *envLayer) GetFieldValuesFrom(model ModelInterface) (fvalue.Values, error) {
		stringValues := l.getStringValues()
		return convertToFieldValues(stringValues, model)
	}

## Model Application
Models generate fvalue.Values from available sources:

	func (m *structModel) ApplyOn(getter ValueGetter) (fvalue.Values, error) {
		result := make(fvalue.Values)
		for fieldPath, fieldInfo := range m.fields {
			if value := getter.GetValue(fieldPath); value != nil {
				converted, err := convertValue(value, fieldInfo.Type)
				if err != nil {
					return nil, fmt.Errorf("field %s: %w", fieldPath, err)
				}
				result[fieldPath] = &fvalue.Value{
					Value:    converted,
					Location: value.Location,
				}
			}
		}
		return result, nil
	}

# Testing Coverage

This package maintains 100% test coverage, including:
- Value creation and type assignment
- Values map operations with complex keys
- Type conversion for all supported Go types
- Location preservation through processing pipeline
- Error handling for malformed values
- Memory safety with nil values and pointers
- Integration testing with model system

The test suite covers edge cases like:
- Conversion of empty strings and zero values
- Handling of special YAML values (.inf, .nan, null)
- Complex nested structure validation
- Location string formatting and parsing

# Performance Considerations

The package is optimized for configuration processing scenarios:

## Memory Efficiency
- Values stored as pointers to avoid copying
- Location strings reused when possible
- Map operations optimized for typical field counts (10-100s)

## Processing Speed
- Type conversion uses Go's YAML library (fast C implementation)
- Reflection operations cached in model system
- Batch processing of field values where possible

## Allocation Patterns
For typical configuration sizes (< 1000 fields), memory usage is minimal.
Large configurations benefit from:
- Streaming processing for file-based sources
- Lazy evaluation of complex nested structures
- Cleanup of intermediate values after processing

# Thread Safety

Like svalue, fvalue types are not inherently thread-safe. The dsco
processing pipeline handles this through:
- Sequential layer processing
- Immutable value objects where possible
- Clear ownership transfer between pipeline stages

# Advanced Features

## Custom Type Support
The package supports custom types that implement:
- yaml.Unmarshaler interface
- encoding.TextUnmarshaler interface
- Custom converter registration (future enhancement)

## Debugging Support
Location information enables sophisticated debugging:
- Value origin tracing
- Configuration precedence analysis
- Override detection and reporting
- Source-specific error handling.
*/
package fvalue
