// Package dsco layer policy definitions.
// This file defines constraint policies for configuration layers,
// distinguishing between strict and normal processing modes.
package dsco

// constraintLayerPolicies manages a collection of layer policies that
// determine how configuration values are processed and validated.
type constraintLayerPolicies []constraintLayerPolicy

// constraintLayerPolicy defines the behavior for a configuration layer,
// combining value retrieval capabilities with strictness constraints.
type constraintLayerPolicy interface {
	FieldValuesGetter
	isStrict() bool
}

// strictLayer enforces strict processing where all provided configuration
// values must be consumed during the filling process.
type strictLayer struct {
	FieldValuesGetter
}

// isStrict returns true, indicating this layer requires all values to be used.
func (*strictLayer) isStrict() bool {
	return true
}

// newStrictLayer creates a new strict layer policy wrapping the provided
// field values getter with strict consumption validation.
//
//nolint:ireturn // this is required
func newStrictLayer(bg FieldValuesGetter) constraintLayerPolicy {
	return &strictLayer{
		FieldValuesGetter: bg,
	}
}

// normalLayer allows flexible processing where unused configuration
// values do not trigger validation errors.
type normalLayer struct {
	FieldValuesGetter
}

// isStrict returns false, indicating unused values are permitted.
func (*normalLayer) isStrict() bool {
	return false
}

// newNormalLayer creates a new normal layer policy wrapping the provided
// field values getter with flexible consumption rules.
//
//nolint:ireturn // this is required
func newNormalLayer(bg FieldValuesGetter) constraintLayerPolicy {
	return &normalLayer{
		FieldValuesGetter: bg,
	}
}
