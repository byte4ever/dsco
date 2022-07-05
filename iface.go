package dsco

import (
	"github.com/byte4ever/dsco/svalue"
)

type PoliciesGetter interface {
	GetPolicies() (constraintLayerPolicies, error)
}

// StringValuesProvider defines the behaviour if a string value provider.
type StringValuesProvider interface {
	GetStringValues() svalue.Values
}

type NamedStringValuesProvider interface {
	StringValuesProvider
	GetName() string
}
