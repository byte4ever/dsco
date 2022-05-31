package goconf

func V[T any](v T) *T {
	return &v
}
