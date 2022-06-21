package walker

// TODO :- lmartin 6/21/22 -: don't use internal states

type constraintLayerPolicy interface {
	BaseGetter
	isStrict() bool
}

type constraintLayer struct {
	BaseGetter
	strictMode bool
}

func (p *constraintLayer) isStrict() bool {
	return p.strictMode
}

func strictLayer(bg BaseGetter) *constraintLayer {
	return &constraintLayer{
		BaseGetter: bg,
		strictMode: true,
	}
}

func normalLayer(bg BaseGetter) *constraintLayer {
	return &constraintLayer{
		BaseGetter: bg,
	}
}
