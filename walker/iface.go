package walker

type BaseGetter interface {
	GetBaseFor(
		inputModel any,
	) (Base, []error)
}
