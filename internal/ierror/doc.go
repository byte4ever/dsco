/*
Package ierror provides indexed error handling for dsco's configuration
processing system.

# Overview

The ierror package implements an error type that associates errors with specific
indices, providing precise error reporting for ordered operations like layer
processing, argument parsing, and configuration validation. This enables dsco
to provide detailed error messages that help users identify exactly which
element in a sequence caused a problem.

# Core Type

## IError

IError represents an error that occurred at a specific index within an ordered
operation:

	type IError struct {
		Err   error   // The underlying error that occurred
		Info  string  // Contextual information about the operation
		Index int     // 0-based index where the error occurred
	}

### Fields

  - **Err**: The underlying error that occurred. This could be any error type
    from parsing, validation, or processing operations.

  - **Info**: A descriptive string providing context about what operation was
    being performed (e.g., "layer", "argument", "field processing").

  - **Index**: The 0-based index identifying the specific element in the sequence
    where the error occurred. This corresponds to positions in slices, arrays,
    or other ordered collections.

# Error Interface Implementation

IError implements the standard Go error interface:

	func (e IError) Error() string {
		return fmt.Sprintf(
			"%s #%d: %s",
			e.Info,
			e.Index,
			e.Err.Error(),
		)
	}

## Error Message Format

Error messages follow a consistent format: `"{info} #{index}: {underlying_error}"`

Examples:

	"layer #2: invalid environment prefix: INVALID-prefix"
	"argument #3: --host previous found at position #1: duplicate param"
	"field processing #5: cannot convert 'invalid' to integer"

# Error Unwrapping

IError supports Go 1.13+ error unwrapping:

	func (e *IError) Unwrap() error {
		return e.Err
	}

This enables proper error chain inspection with errors.Is and errors.As:

	var ierr *ierror.IError
	if errors.As(err, &ierr) {
		fmt.Printf("Error at index %d: %v\n", ierr.Index, ierr.Err)
	}

	// Check for specific underlying error types
	if errors.Is(err, someSpecificError) {
		// Handle specific error type
	}

# Usage Patterns

## Layer Processing Errors

When processing configuration layers in sequence:

	func processLayers(layers []Layer) error {
		for i, layer := range layers {
			if err := layer.Process(); err != nil {
				return &ierror.IError{
					Err:   err,
					Info:  "layer",
					Index: i,
				}
			}
		}
		return nil
	}

	// Results in errors like: "layer #2: invalid configuration format"

## Argument Processing Errors

When parsing command-line arguments or similar ordered inputs:

	func parseArguments(args []string) error {
		for i, arg := range args {
			if err := parseArgument(arg); err != nil {
				return &ierror.IError{
					Err:   err,
					Info:  "argument",
					Index: i,
				}
			}
		}
		return nil
	}

	// Results in errors like: "argument #3: invalid format"

## Field Processing Errors

When processing struct fields or configuration values:

	func processFields(fields []Field) error {
		for i, field := range fields {
			if err := validateField(field); err != nil {
				return &ierror.IError{
					Err:   err,
					Info:  "field validation",
					Index: i,
				}
			}
		}
		return nil
	}

	// Results in errors like: "field validation #1: required field missing"

# Integration with dsco Components

## Layer System

The ierror package is heavily used in dsco's layer processing:

	func (layers Layers) GetPolicies() (constraintLayerPolicies, error) {
		var errs LayerErrors

		for index, layer := range layers {
			err := layer.register(bo)
			if err != nil {
				errs.Add(
					ierror.IError{
						Index: index,
						Info:  "layer",
						Err:   err,
					},
				)
			}
		}

		if errs.None() {
			return bo.builders, nil
		}

		return nil, errs
	}

## Command Line Processing

Used for reporting command-line argument parsing errors:

	type ParamError struct {
		Positions []int    // 1-based positions
		Errs      []error  // Corresponding errors
	}

	func (e *ParamError) Error() string {
		// Uses ierror concepts for consistent formatting
		return fmt.Sprintf("cmdline issue at position #%d: %s",
			e.Positions[0], e.Errs[0].Error())
	}

## Model Building

Used during configuration model construction for field processing errors:

	func buildModel(fields []reflect.StructField) error {
		for i, field := range fields {
			if err := analyzeField(field); err != nil {
				return &ierror.IError{
					Err:   err,
					Info:  "model building",
					Index: i,
				}
			}
		}
		return nil
	}

# Error Aggregation

IError works well with dsco's error aggregation systems:

	type MultiError []error

	func (m *MultiError) Add(err error) {
		*m = append(*m, err)
	}

	func processWithAggregation(items []Item) error {
		var errors MultiError

		for i, item := range items {
			if err := processItem(item); err != nil {
				errors.Add(&ierror.IError{
					Err:   err,
					Info:  "item processing",
					Index: i,
				})
			}
		}

		if len(errors) > 0 {
			return errors
		}
		return nil
	}

# Best Practices

## Context Information

Provide clear, specific context in the Info field:

	// Good: Specific operation context
	&ierror.IError{Info: "layer registration", ...}
	&ierror.IError{Info: "command line parsing", ...}
	&ierror.IError{Info: "field validation", ...}

	// Less helpful: Generic context
	&ierror.IError{Info: "processing", ...}
	&ierror.IError{Info: "error", ...}

## Index Consistency

Ensure index values match the user's mental model:

	// For command-line arguments (1-based in user documentation)
	&ierror.IError{
		Index: i,  // But document as position i+1 to users
		...
	}

	// For internal processing (0-based is fine)
	&ierror.IError{
		Index: i,  // Direct slice/array index
		...
	}

## Error Chain Preservation

Always wrap the original error to preserve error chains:

	// Good: Preserves original error
	return &ierror.IError{
		Err:   originalErr,
		Info:  "operation",
		Index: i,
	}

	// Bad: Loses original error information
	return &ierror.IError{
		Err:   fmt.Errorf("something failed"),
		Info:  "operation",
		Index: i,
	}

# Error Handling Examples

## Basic Error Checking

	err := processItems(items)
	if err != nil {
		var ierr *ierror.IError
		if errors.As(err, &ierr) {
			log.Printf("Failed at %s #%d: %v", ierr.Info, ierr.Index, ierr.Err)
			return ierr
		}
		// Handle other error types
	}

## Specific Error Type Checking

	err := processConfiguration(config)
	if err != nil {
		var ierr *ierror.IError
		if errors.As(err, &ierr) {
			// Check the underlying error type
			if errors.Is(ierr.Err, ErrInvalidFormat) {
				log.Printf("Format error in %s at position %d", ierr.Info, ierr.Index)
				return handleFormatError(ierr)
			}
		}
	}

# Thread Safety

IError instances are safe for concurrent use:

- **Immutable after creation**: IError fields should not be modified after creation
- **Safe for concurrent reads**: Multiple goroutines can safely read IError fields
- **Error() method is safe**: The Error() method can be called concurrently

# Memory Considerations

IError is designed for efficient memory usage:

- **Small struct size**: Only three fields (error interface, string, int)
- **No hidden allocations**: Error() method uses fmt.Sprintf, which allocates a new string
- **Error preservation**: Maintains reference to original error without copying

# Performance Characteristics

The IError type is optimized for error reporting scenarios:

- **Fast creation**: Minimal overhead when creating IError instances
- **Efficient formatting**: Error() method uses efficient fmt.Sprintf formatting
- **Low memory footprint**: Minimal memory usage per error instance

This package provides a simple but powerful foundation for precise error
reporting in dsco's configuration processing pipeline, enabling users to
quickly identify and fix configuration problems.
*/
package ierror
