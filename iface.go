package dsco

import (
	"github.com/byte4ever/dsco/internal/svalues"
)

type PoliciesGetter interface {
	GetPolicies() (constraintLayerPolicies, error)
}

// StringValuesProvider defines the behaviour if a string value provider.
type StringValuesProvider interface {
	GetStringValues() svalues.StringValues
}
