package dsco

import (
	"github.com/byte4ever/dsco/internal/svalue"
)

type PoliciesGetter interface {
	GetPolicies() (constraintLayerPolicies, error)
}

// StringValuesProvider defines the behaviour if a string value provider.
type StringValuesProvider interface {
	GetStringValues() svalue.Values
}
