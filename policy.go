package dsco

type constraintLayerPolicies []constraintLayerPolicy

type constraintLayerPolicy interface {
	FieldValuesGetter
	isStrict() bool
}

type strictLayer struct {
	FieldValuesGetter
}

func (*strictLayer) isStrict() bool {
	return true
}

//nolint:ireturn // this is required
func newStrictLayer(bg FieldValuesGetter) constraintLayerPolicy {
	return &strictLayer{
		FieldValuesGetter: bg,
	}
}

type normalLayer struct {
	FieldValuesGetter
}

func (*normalLayer) isStrict() bool {
	return false
}

//nolint:ireturn // this is required
func newNormalLayer(bg FieldValuesGetter) constraintLayerPolicy {
	return &normalLayer{
		FieldValuesGetter: bg,
	}
}
