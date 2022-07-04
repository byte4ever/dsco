package dsco

import (
	"github.com/byte4ever/dsco/ifaces"
)

type constraintLayerPolicies []constraintLayerPolicy

type constraintLayerPolicy interface {
	ifaces.FieldValuesGetter
	isStrict() bool
}

type strictLayer struct {
	ifaces.FieldValuesGetter
}

func (*strictLayer) isStrict() bool {
	return true
}

//nolint:ireturn // this is required
func newStrictLayer(bg ifaces.FieldValuesGetter) constraintLayerPolicy {
	return &strictLayer{
		FieldValuesGetter: bg,
	}
}

type normalLayer struct {
	ifaces.FieldValuesGetter
}

func (*normalLayer) isStrict() bool {
	return false
}

//nolint:ireturn // this is required
func newNormalLayer(bg ifaces.FieldValuesGetter) constraintLayerPolicy {
	return &normalLayer{
		FieldValuesGetter: bg,
	}
}
