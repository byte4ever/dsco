package walker

type ConstraintLayerPolicy interface {
	BaseGetter
	IsStrict() bool
}

type ConstraintLayer struct {
	BaseGetter
	isStrict bool
}

func (p *ConstraintLayer) IsStrict() bool {
	return p.isStrict
}

func StrictLayer(bg BaseGetter) *ConstraintLayer {
	return &ConstraintLayer{
		BaseGetter: bg,
		isStrict:   true,
	}
}

func NormalLayer(bg BaseGetter) *ConstraintLayer {
	return &ConstraintLayer{
		BaseGetter: bg,
	}
}
