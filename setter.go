package dsco

func V[T any](v T) *T {
	return &v
}
