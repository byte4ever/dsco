package dsco

import (
	"reflect"

	"github.com/byte4ever/dsco/internal"
	"github.com/byte4ever/dsco/internal/fvalue"
	"github.com/byte4ever/dsco/internal/plocation"
)

// FieldValuesGetter defines the ability to get field values from a
// configuration model. This interface is implemented by configuration sources
// (layers) to provide their
// configuration values in a standard format that can be merged and processed.
type FieldValuesGetter interface {
	// GetFieldValuesFrom extracts configuration values from the given model.
	// Returns a map of field UIDs to their values, or an error if extraction
	// fails.
	GetFieldValuesFrom(model ModelInterface) (fvalue.Values, error)
}

// ModelInterface represents a configuration structure model that defines how
// configuration data should be processed, validated, and filled into Go
// structs. It provides the metadata and operations needed for dsco's
// configuration processing.
type ModelInterface interface {
	// TypeName returns a human-readable name for the configuration type.
	// Used for error reporting and debugging.
	TypeName() string

	// ApplyOn extracts configuration values using the provided value getter.
	// Returns field values mapped by their unique identifiers.
	ApplyOn(g internal.ValueGetter) (fvalue.Values, error)

	// Expand processes struct expansion using the provided struct expander.
	// This handles complex configuration structures and nested types.
	Expand(g internal.StructExpander) error

	// GetFieldValuesFor extracts field values from the given struct value.
	// Used to convert Go structs into the internal field value representation.
	GetFieldValuesFor(id string, v reflect.Value) fvalue.Values

	// Fill populates the target struct with values from configuration layers.
	// Returns location information for debugging and error reporting.
	Fill(
		inputModelValue reflect.Value,
		layers []fvalue.Values,
	) (plocation.Locations, error)
}
