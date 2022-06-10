package dsco

// R returns a reference to any value thanks to generics in go1.18.
func R[T any](v T) *T {
	return &v
}
