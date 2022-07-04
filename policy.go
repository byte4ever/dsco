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

func (s *strictLayer) isStrict() bool {
	return true
}

func newStrictLayer(bg ifaces.FieldValuesGetter) constraintLayerPolicy {
	return &strictLayer{
		FieldValuesGetter: bg,
	}
}

type normalLayer struct {
	ifaces.FieldValuesGetter
}

func (n *normalLayer) isStrict() bool {
	return false
}

func newNormalLayer(bg ifaces.FieldValuesGetter) constraintLayerPolicy {
	return &normalLayer{
		FieldValuesGetter: bg,
	}
}
