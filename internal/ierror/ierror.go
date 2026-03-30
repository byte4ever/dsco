package ierror

import (
	"fmt"
)

// IError wraps an error with contextual information including an index
// and descriptive info, useful for tracking errors in ordered operations.
type IError struct {
	Err   error  // The underlying error that occurred
	Info  string // Contextual description of the operation
	Index int    // Index position where the error occurred
}

// Error formats the indexed error with context information, providing
// a clear indication of where and what type of error occurred.
func (e IError) Error() string {
	return fmt.Sprintf(
		"%s #%d: %s",
		e.Info,
		e.Index,
		e.Err.Error(),
	)
}

// Unwrap returns the underlying error, enabling error chain traversal
// for error.Is() and error.As() functionality.
func (e *IError) Unwrap() error {
	return e.Err
}
