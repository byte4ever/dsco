package walker

import (
	"github.com/byte4ever/dsco/walker/ifaces"
)

// TODO :- lmartin 6/21/22 -: don't use internal states

type constraintLayerPolicies []constraintLayerPolicy

type constraintLayerPolicy interface {
	ifaces.FieldValuesGetter
	isStrict() bool
}

type constraintLayer struct {
	ifaces.FieldValuesGetter
	strictMode bool
}

func (p *constraintLayer) isStrict() bool {
	return p.strictMode
}

func strictLayer(bg ifaces.FieldValuesGetter) *constraintLayer {
	return &constraintLayer{
		FieldValuesGetter: bg,
		strictMode:        true,
	}
}

func normalLayer(bg ifaces.FieldValuesGetter) *constraintLayer {
	return &constraintLayer{
		FieldValuesGetter: bg,
	}
}
