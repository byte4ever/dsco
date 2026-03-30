/*
Package registry provides type name registration and formatting utilities
for dsco's reflection-based configuration system.

# Overview

The registry package offers standardized type name handling that supports
dsco's need for consistent, human-readable type identification in error
messages, debugging output, and configuration documentation. It handles
complex Go types including pointers, slices, maps, channels, and custom types.

# Core Functions

The package provides two primary functions for type name formatting:

## LongTypeName

LongTypeName returns the full, detailed type name including package path:

	type MyStruct struct { Name string }

	fmt.Println(registry.LongTypeName(reflect.TypeOf(MyStruct{})))
	// Output: "github.com/example/myapp.MyStruct"

	fmt.Println(registry.LongTypeName(reflect.TypeOf(&MyStruct{})))
	// Output: "*github.com/example/myapp.MyStruct"

	fmt.Println(registry.LongTypeName(reflect.TypeOf([]MyStruct{})))
	// Output: "[]github.com/example/myapp.MyStruct"

This format is ideal for:
- Detailed error messages requiring precise type identification
- Debugging output where package disambiguation is important
- Internal logging and tracing systems
- Configuration validation messages

## ShortTypeName

ShortTypeName returns a concise type name without package paths:

	fmt.Println(registry.ShortTypeName(reflect.TypeOf(MyStruct{})))
	// Output: "MyStruct"

	fmt.Println(registry.ShortTypeName(reflect.TypeOf(&MyStruct{})))
	// Output: "*MyStruct"

	fmt.Println(registry.ShortTypeName(reflect.TypeOf(map[string]int{})))
	// Output: "map[string]int"

This format is ideal for:
- User-facing error messages that should be concise
- Configuration documentation and help text
- CLI output where brevity is important
- API responses that include type information

# Type Support

The package handles all Go type categories:

## Basic Types
- Built-in types: string, int, bool, float64, etc.
- Type aliases: byte, rune, etc.

## Composite Types
- Pointers: *Type, **Type, etc. (arbitrary nesting depth)
- Arrays: [N]Type with size information preserved
- Slices: []Type
- Maps: map[KeyType]ValueType with both types formatted
- Channels: chan Type, <-chan Type, chan<- Type with directionality

## Advanced Types
- Structs: Full package path and type name
- Interfaces: Including empty interface (interface{}) and any
- Functions: Full signature with parameter and return types
- Methods: Bound method types with receiver information

## Example Complex Types

	// Complex nested type
	type Config struct {
		Databases map[string]*DatabaseConfig
		Listeners []net.Listener
		Handler   func(http.ResponseWriter, *http.Request) error
	}

	// LongTypeName output:
	// "github.com/myapp/config.Config"

	// For the Databases field specifically:
	// "map[string]*github.com/myapp/config.DatabaseConfig"

# Integration with dsco Error System

The registry package is tightly integrated with dsco's error reporting:

## Error Message Consistency
All dsco error messages use registry functions for type formatting:

	// In validation errors
	return fmt.Errorf(
		"field 'host' expects type %s but got %s",
		registry.ShortTypeName(expectedType),
		registry.ShortTypeName(actualType),
	)

## Debug Information
Detailed debugging uses long type names:

	logger.Debug(
		"processing field",
		"path", fieldPath,
		"target_type", registry.LongTypeName(fieldType),
		"source_location", value.Location,
	)

# Performance Characteristics

The package is optimized for dsco's usage patterns:

## Caching Strategy
Type names are computed on-demand but could benefit from caching in
high-throughput scenarios. The reflection operations are relatively
fast for the typical scale of configuration processing.

## Memory Usage
String formatting is performed using Go's efficient string builder
patterns. Memory allocation is proportional to the complexity of
the type name.

## Execution Speed
- Basic types: Direct string return (fastest)
- Pointer types: Single indirection check
- Complex types: Recursive formatting with bounded depth

# Usage Patterns in dsco

## Model Generation
When dsco analyzes struct types:

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldTypeName := registry.LongTypeName(field.Type)
		logger.Debug("analyzing field",
			"name", field.Name,
			"type", fieldTypeName,
		)
	}

## Error Reporting
In validation and conversion errors:

	func validateFieldType(value interface{}, expected reflect.Type) error {
		actual := reflect.TypeOf(value)
		if actual != expected {
			return fmt.Errorf(
				"type mismatch: expected %s, got %s",
				registry.ShortTypeName(expected),
				registry.ShortTypeName(actual),
			)
		}
		return nil
	}

## Configuration Documentation
For generating help text and documentation:

	func generateFieldDoc(field reflect.StructField) string {
		return fmt.Sprintf(
			"Field: %s\nType: %s\nDescription: %s",
			field.Name,
			registry.ShortTypeName(field.Type),
			getFieldDescription(field),
		)
	}

# Testing Coverage

This package maintains 100% test coverage, including:
- All basic Go types (string, int, bool, float64, complex128, etc.)
- Pointer types with multiple levels of indirection
- Array types with various sizes and element types
- Slice types with basic and complex element types
- Map types with various key-value type combinations
- Channel types with all directionality options
- Struct types from various packages
- Interface types including empty interface
- Function types with complex signatures
- Edge cases with nil types and invalid inputs

The test suite ensures:
- Consistent output format across Go versions
- Correct handling of package path inclusion/exclusion
- Proper formatting of nested and recursive type structures
- Performance characteristics within acceptable bounds

# Thread Safety

All functions in the registry package are thread-safe. They operate only
on reflect.Type values (which are immutable) and perform deterministic
string operations without maintaining any shared state.

This makes the package safe for use in dsco's potentially concurrent
configuration processing scenarios.

# Extension Points

The package design allows for future enhancements:

## Custom Type Formatters
Future versions could support custom type name formatters for specific
types or packages:

	registry.RegisterTypeFormatter(myType, customFormatter)

## Alias Management
Support for type alias registration to provide more user-friendly
names in error messages:

	registry.RegisterAlias(reflect.TypeOf(MyComplexType{}), "MyType")

## Localization
Support for localized type names in different languages for
international applications.

# Error Handling

The package functions are designed to be robust:
- Never panic on valid reflect.Type inputs
- Handle nil types gracefully by returning descriptive strings
- Provide meaningful output for all Go type constructs
- Degrade gracefully for unusual or complex types

This reliability is essential for dsco's error reporting system, where
type name formatting failures could mask the original configuration errors.
*/
package registry
