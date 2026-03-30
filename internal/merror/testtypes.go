package merror

// testError is a lightweight error type used only by tests.
// This interface allows us to create mocks for testing error scenarios
// without depending on concrete error implementations.
type testError interface {
	error
}
