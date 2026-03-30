/*
Package merror provides multiple error aggregation functionality for dsco's
configuration processing system.

# Overview

The merror package implements an error aggregation type that can collect and
manage multiple errors from different sources or operations. This is essential
for dsco's validation and processing pipeline, where multiple configuration
issues may need to be reported simultaneously rather than stopping at the
first error.

# Core Type

## MError

MError represents a collection of multiple errors:

	type MError []error

MError is a slice of errors that implements the standard error interface,
allowing it to be used anywhere a regular error is expected while providing
the ability to aggregate multiple error conditions.

# Error Interface Implementation

## Error() Method

MError implements the error interface by concatenating all contained errors:

	func (m MError) Error() string {
		if len(m) == 0 {
			return ""
		}

		var sb strings.Builder
		sb.WriteString(m[0].Error())

		for _, err := range m[1:] {
			sb.WriteRune('\n')
			sb.WriteString(err.Error())
		}

		return sb.String()
	}

### Error Message Format

Multiple errors are presented as newline-separated messages:

	error1: first problem
	error2: second problem
	error3: third problem

This format makes it easy for users to identify and address multiple
configuration issues in a single run.

# Error Unwrapping Support

## Is() Method

MError supports Go 1.13+ error inspection with errors.Is:

	func (MError) Is(err error) bool {
		return errors.Is(err, Err)
	}

This allows checking if an MError represents the general multiple error condition:

	if errors.Is(err, merror.Err) {
		// This is a multiple error condition
	}

## As() Method

MError supports error type assertion with errors.As:

	func (m MError) As(errAs any) bool {
		for _, err := range m {
			if errors.As(err, errAs) {
				return true
			}
		}
		return false
	}

This enables checking for specific error types within the error collection:

	var validationErr *ValidationError
	if errors.As(multiErr, &validationErr) {
		// At least one validation error exists in the collection
	}

# Management Methods

## Add() Method

Add appends a new error to the collection:

	func (m *MError) Add(err error) {
		*m = append(*m, err)
	}

Usage example:

	var errors merror.MError
	if err := validateField1(); err != nil {
		errors.Add(err)
	}
	if err := validateField2(); err != nil {
		errors.Add(err)
	}

## None() Method

None returns true if no errors have been collected:

	func (m *MError) None() bool {
		return len(*m) == 0
	}

Usage for conditional error returns:

	var errors merror.MError
	// ... collect errors ...

	if errors.None() {
		return result, nil
	}
	return nil, errors

## Count() Method

Count returns the number of errors in the collection:

	func (m *MError) Count() int {
		return len(*m)
	}

Useful for logging and metrics:

	log.Printf("Configuration validation found %d errors", errors.Count())

# Usage Patterns

## Validation Error Aggregation

Collect multiple validation errors before returning:

	func validateConfig(config *Config) error {
		var errors merror.MError

		if config.Host == nil {
			errors.Add(fmt.Errorf("host is required"))
		}

		if config.Port != nil && *config.Port <= 0 {
			errors.Add(fmt.Errorf("port must be positive"))
		}

		if config.Timeout != nil && *config.Timeout < 0 {
			errors.Add(fmt.Errorf("timeout cannot be negative"))
		}

		if errors.None() {
			return nil
		}
		return errors
	}

## Layer Processing Error Collection

Aggregate errors from multiple configuration layers:

	func processLayers(layers []Layer) error {
		var errors merror.MError

		for i, layer := range layers {
			if err := layer.Process(); err != nil {
				errors.Add(&ierror.IError{
					Index: i,
					Info:  "layer processing",
					Err:   err,
				})
			}
		}

		if errors.None() {
			return nil
		}
		return errors
	}

## Field Processing Error Aggregation

Collect errors during struct field processing:

	func processFields(fields []Field) error {
		var errors merror.MError

		for _, field := range fields {
			if err := validateFieldType(field); err != nil {
				errors.Add(fmt.Errorf("field %s: %w", field.Name, err))
			}

			if err := validateFieldTags(field); err != nil {
				errors.Add(fmt.Errorf("field %s tags: %w", field.Name, err))
			}
		}

		if errors.None() {
			return nil
		}
		return errors
	}

# Integration with dsco Components

## Model Building

Used during configuration model construction:

	type ModelError struct {
		merror.MError
	}

	func buildModel(structType reflect.Type) (*Model, error) {
		var errs ModelError

		fields, fieldErrs := analyzeFields(structType)
		for _, err := range fieldErrs {
			errs.Add(err)
		}

		if !errs.None() {
			return nil, errs
		}

		return &Model{fields: fields}, nil
	}

## Layer Registration

Used for collecting layer registration errors:

	type LayerErrors struct {
		merror.MError
	}

	func registerLayers(layers []Layer) error {
		var errs LayerErrors

		for i, layer := range layers {
			if err := layer.Register(); err != nil {
				errs.Add(&ierror.IError{
					Index: i,
					Info:  "layer",
					Err:   err,
				})
			}
		}

		if errs.None() {
			return nil
		}
		return errs
	}

## Struct Node Processing

Used in struct node error collection:

	type StructNodeError struct {
		merror.MError
	}

	func (n StructNode) Fill(value reflect.Value, layers []fvalue.Values) error {
		var errs StructNodeError

		for _, index := range n.Index {
			if err := index.Node.Fill(subValue, layers); err != nil {
				errs.Add(err)
			}
		}

		if errs.None() {
			return nil
		}
		return errs
	}

# Error Type Embedding

MError is often embedded in specific error types for domain-specific error handling:

	type ValidationErrors struct {
		merror.MError
	}

	type ProcessingErrors struct {
		merror.MError
	}

	type LayerErrors struct {
		merror.MError
	}

This pattern allows:
- Type-specific error handling with errors.As
- Maintaining MError functionality
- Domain-specific error behavior when needed

# Best Practices

## Early Error Collection

Collect errors early rather than stopping at first error:

	// Good: Collect all validation issues
	func validate(config *Config) error {
		var errors merror.MError
		// ... add multiple validation errors ...
		if errors.None() {
			return nil
		}
		return errors
	}

	// Less helpful: Stop at first error
	func validate(config *Config) error {
		if config.Host == nil {
			return fmt.Errorf("host required")
		}
		// ... other validations never reached ...
	}

## Meaningful Error Messages

Provide context in individual error messages:

	// Good: Clear, specific messages
	errors.Add(fmt.Errorf("field 'host': value cannot be empty"))
	errors.Add(fmt.Errorf("field 'port': value %d out of range 1-65535", port))

	// Less helpful: Generic messages
	errors.Add(fmt.Errorf("invalid value"))
	errors.Add(fmt.Errorf("validation failed"))

## Conditional Return Pattern

Use the None() method for clean conditional returns:

	var errors merror.MError
	// ... collect errors ...

	if errors.None() {
		return result, nil  // Success case
	}
	return nil, errors     // Error case

# Error Handling Examples

## Basic Usage

	err := validateConfiguration(config)
	if err != nil {
		var merrs merror.MError
		if errors.As(err, &merrs) {
			fmt.Printf("Found %d configuration errors:\n", merrs.Count())
			fmt.Println(merrs.Error())
		} else {
			fmt.Printf("Single error: %v\n", err)
		}
	}

## Specific Error Type Checking

	err := processLayers(layers)
	if err != nil {
		var layerErrs LayerErrors
		if errors.As(err, &layerErrs) {
			log.Printf("Layer processing failed with %d errors", layerErrs.Count())
			for i, subErr := range layerErrs.MError {
				log.Printf("  Error %d: %v", i+1, subErr)
			}
		}
	}

# Thread Safety

MError has specific thread safety characteristics:

- **Not safe for concurrent writes**: Adding errors concurrently requires external synchronization
- **Safe for concurrent reads**: Once populated, MError can be read concurrently
- **Error() method is safe**: String generation can be called concurrently

# Memory Management

MError is designed for efficient memory usage:

- **Slice-based storage**: Efficient append operations with Go's slice growth
- **Reference-based**: Stores error interfaces, not error values
- **Minimal overhead**: Only slice header overhead plus error references

# Performance Characteristics

- **Fast append**: O(1) amortized for Add() operations
- **Linear error formatting**: O(n) where n is number of errors for Error()
- **Efficient type checking**: O(n) for As() method to check all contained errors

The merror package provides essential error aggregation capabilities that enable
dsco to provide comprehensive error reporting, helping users identify and fix
multiple configuration issues efficiently.
*/
package merror
