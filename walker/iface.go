package walker

// BaseGetter defines the ability to get a path/value set (base).
type BaseGetter interface {
	GetBaseFor(
		inputModel any,
	) (Base, []error)
}
