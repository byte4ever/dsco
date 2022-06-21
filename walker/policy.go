package walker

// TODO :- lmartin 6/21/22 -: don't use internal states

type constraintLayerPolicies []constraintLayerPolicy

type constraintLayerPolicy interface {
	FieldValuesGetter
	isStrict() bool
}

type constraintLayer struct {
	FieldValuesGetter
	strictMode bool
}

func (p *constraintLayer) isStrict() bool {
	return p.strictMode
}

func strictLayer(bg FieldValuesGetter) *constraintLayer {
	return &constraintLayer{
		FieldValuesGetter: bg,
		strictMode:        true,
	}
}

func normalLayer(bg FieldValuesGetter) *constraintLayer {
	return &constraintLayer{
		FieldValuesGetter: bg,
	}
}
