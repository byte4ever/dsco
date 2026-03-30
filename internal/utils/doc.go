/*
Package utils provides utility functions for dsco's configuration processing
system.

# Overview

The utils package contains common utility functions used throughout dsco for
string processing, formatting, and data transformation. These utilities support
various dsco operations including key name generation, sequence formatting,
and string case conversion.

# Key Generation

## GetKeyName Function

GetKeyName generates configuration keys from struct field information:

	func GetKeyName(prefix string, fieldType reflect.StructField) string

This function is essential for mapping Go struct fields to configuration keys
in various formats (command-line, environment variables, etc.).

### Key Generation Logic

1. **Extract field name**: First attempts to use YAML tag, falls back to field name
2. **Apply case conversion**: Converts Go CamelCase to snake_case
3. **Add prefix**: Prepends prefix with hyphen separator if provided

### Examples

	type Config struct {
		DatabaseHost     string `yaml:"db-host"`
		MaxConnections   int    `yaml:"max_connections"`
		APIKey           string // No tag, uses field name
		TimeoutDuration  time.Duration `yaml:"timeout"`
	}

	// With empty prefix:
	GetKeyName("", field("DatabaseHost"))  // "db-host" (from yaml tag)
	GetKeyName("", field("APIKey"))        // "api-key" (converted from field name)

	// With prefix:
	GetKeyName("app", field("DatabaseHost"))  // "app-db-host"
	GetKeyName("api", field("TimeoutDuration")) // "api-timeout"

# String Case Conversion

## ToSnakeCase Function

ToSnakeCase converts CamelCase strings to snake_case:

	func ToSnakeCase(s string) string

### Conversion Rules

- Inserts underscores before uppercase letters (except the first character)
- Converts all letters to lowercase
- Handles consecutive uppercase letters properly
- Preserves existing underscores and hyphens

### Examples

	ToSnakeCase("DatabaseHost")     // "database_host"
	ToSnakeCase("APIKey")           // "api_key"
	ToSnakeCase("HTTPSProxy")       // "https_proxy"
	ToSnakeCase("XMLHttpRequest")   // "xml_http_request"
	ToSnakeCase("IOTimeout")        // "io_timeout"

# Sequence Formatting

The utils package provides functions for formatting sequences of values in
human-readable formats.

## FormatIndexSequence Function

FormatIndexSequence formats slices of integers with proper English conjunction:

	func FormatIndexSequence(indexes []int) string

### Examples

	FormatIndexSequence([]int{1})           // "#1"
	FormatIndexSequence([]int{1, 2})        // "#1 and #2"
	FormatIndexSequence([]int{1, 2, 3})     // "#1, #2 and #3"
	FormatIndexSequence([]int{1, 2, 3, 4})  // "#1, #2, #3 and #4"

## FormatStringSequence Function

FormatStringSequence formats slices of strings with proper English conjunction:

	func FormatStringSequence(values []string) string

### Examples

	FormatStringSequence([]string{"host"})                    // "host"
	FormatStringSequence([]string{"host", "port"})            // "host and port"
	FormatStringSequence([]string{"host", "port", "timeout"}) // "host, port and timeout"

# Usage in dsco Components

## Environment Variable Processing

Used for generating environment variable keys:

	func generateEnvKey(prefix string, field reflect.StructField) string {
		key := utils.GetKeyName("", field)
		return fmt.Sprintf("%s_-%s", prefix, strings.ToUpper(key))
	}

	// For field `DatabaseHost string yaml:"db-host"`
	// Results in: "MYAPP_-DB-HOST"

## Command Line Processing

Used for generating command-line argument keys:

	func generateCmdlineKey(field reflect.StructField) string {
		return "--" + utils.GetKeyName("", field)
	}

	// For field `MaxRetries int yaml:"max_retries"`
	// Results in: "--max_retries"

## Error Message Formatting

Used for formatting error messages with multiple items:

	func formatAmbiguousKeys(keys []string) string {
		return fmt.Sprintf("%s %s",
			utils.FormatStringSequence(keys),
			getErrorSuffix(len(keys)))
	}

	// For keys ["MYAPP_BAD1", "MYAPP_BAD2", "MYAPP_BAD3"]
	// Results in: "MYAPP_BAD1", "MYAPP_BAD2" and "MYAPP_BAD3" are ambiguous

## Model Field Processing

Used during struct field analysis:

	func analyzeField(field reflect.StructField) FieldInfo {
		return FieldInfo{
			GoName:    field.Name,
			ConfigKey: utils.GetKeyName("", field),
			SnakeName: utils.ToSnakeCase(field.Name),
		}
	}

# YAML Tag Processing

The utils package handles YAML tag extraction and processing:

## Tag Format Support

Supports various YAML tag formats:

	`yaml:"simple"`                    // Simple key name
	`yaml:"kebab-case"`                // Kebab-case key name
	`yaml:"snake_case"`                // Snake_case key name
	`yaml:"key,omitempty"`             // Key with options
	`yaml:"key,omitempty,flow"`        // Key with multiple options
	`yaml:",inline"`                   // Inline embedding (key extracted as empty)

## Tag Parsing Logic

1. **Extract tag value**: Gets the `yaml` struct tag
2. **Clean whitespace**: Removes any spaces from the tag
3. **Split options**: Separates key name from options using comma
4. **Return key name**: Returns the first part (key name) or empty string

# String Processing Utilities

## Internal Helper Functions

The package includes several internal helper functions:

### fieldName Function

Extracts the configuration key from a struct field's YAML tag:

	func fieldName(fieldType reflect.StructField) string

This function:
- Extracts the YAML tag value
- Removes whitespace
- Splits on comma to separate key from options
- Returns the key part (first element)

# Best Practices

## Key Name Generation

Use consistent patterns for key generation:

	// Good: Use GetKeyName for consistency
	cmdlineKey := "--" + utils.GetKeyName("", field)
	envKey := strings.ToUpper(utils.GetKeyName(prefix, field))

	// Less consistent: Manual key generation
	cmdlineKey := "--" + strings.ToLower(field.Name)
	envKey := prefix + "_" + strings.ToUpper(field.Name)

## Case Conversion

Use ToSnakeCase for Go naming convention conversion:

	// Good: Consistent case conversion
	snakeCase := utils.ToSnakeCase(goFieldName)

	// Less reliable: Manual conversion
	snakeCase := strings.ToLower(goFieldName) // Misses CamelCase boundaries

## Error Message Formatting

Use sequence formatting for readable error messages:

	// Good: Proper English formatting
	message := fmt.Sprintf("Keys %s are invalid",
		utils.FormatStringSequence(badKeys))

	// Less readable: Manual formatting
	message := fmt.Sprintf("Keys %s are invalid",
		strings.Join(badKeys, ", "))

# Performance Characteristics

The utils package is optimized for typical dsco usage:

## String Operations

- **ToSnakeCase**: O(n) where n is string length, single pass algorithm
- **GetKeyName**: O(n) string operations, minimal allocations
- **FormatSequence**: O(n) where n is slice length, uses strings.Builder

## Memory Usage

- **Minimal allocations**: Functions minimize string allocations
- **Builder usage**: Uses strings.Builder for efficient concatenation
- **No persistent state**: All functions are stateless

## Regular Expressions

String case conversion uses compiled regular expressions:

- **Compiled once**: Regular expressions are compiled at package level
- **Reused safely**: Thread-safe reuse across function calls
- **Efficient matching**: Optimized patterns for common Go naming patterns

# Thread Safety

All functions in the utils package are thread-safe:

- **No shared state**: All functions are pure and stateless
- **Read-only data**: Regular expressions are read-only after compilation
- **Concurrent safe**: Multiple goroutines can call functions concurrently

# Integration Examples

## Complete Key Generation Pipeline

	func generateConfigKeys(field reflect.StructField) ConfigKeys {
		baseKey := utils.GetKeyName("", field)

		return ConfigKeys{
			YAML:        baseKey,
			CommandLine: "--" + baseKey,
			EnvVar:      "APP_-" + strings.ToUpper(strings.ReplaceAll(baseKey, "-", "_")),
			SnakeCase:   utils.ToSnakeCase(field.Name),
		}
	}

## Error Formatting Pipeline

	func formatValidationErrors(fieldErrors map[string][]string) string {
		var messages []string

		for field, errors := range fieldErrors {
			if len(errors) == 1 {
				messages = append(messages,
					fmt.Sprintf("field '%s': %s", field, errors[0]))
			} else {
				messages = append(messages,
					fmt.Sprintf("field '%s': %s",
						field, utils.FormatStringSequence(errors)))
			}
		}

		return utils.FormatStringSequence(messages)
	}

The utils package provides essential string processing and formatting
capabilities that ensure consistent key generation, readable error messages,
and proper case conversion throughout dsco's configuration processing pipeline.
*/
package utils
