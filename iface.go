package dsco

import (
	"github.com/byte4ever/dsco/svalue"
)

// PoliciesGetter defines the ability to retrieve layer policies for
// configuration processing, converting layers into constraint policies.
type PoliciesGetter interface {
	// GetPolicies transforms configuration layers into policy objects
	// that determine how values are processed and validated.
	GetPolicies() (constraintLayerPolicies, error)
}

// StringValuesProvider defines the interface for providers that supply
// string-based configuration values from various sources like environment
// variables, command line arguments, or configuration files.
type StringValuesProvider interface {
	// GetStringValues retrieves all string values provided by this source,
	// returning a map of keys to their corresponding string values and
	// locations.
	GetStringValues() svalue.Values
}

// NamedStringValuesProvider extends StringValuesProvider with identification
// capabilities, allowing the provider to be uniquely named for debugging
// and error reporting purposes.
type NamedStringValuesProvider interface {
	StringValuesProvider
	// GetName returns a human-readable identifier for this provider,
	// used in error messages and debugging output.
	GetName() string
}
