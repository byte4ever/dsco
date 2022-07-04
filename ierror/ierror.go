package ierror

import (
	"fmt"
)

type IError struct {
	Err   error
	Info  string
	Index int
}

func (e IError) Error() string {
	return fmt.Sprintf(
		"%s #%d: %s",
		e.Info,
		e.Index,
		e.Err.Error(),
	)
}

func (e *IError) Unwrap() error {
	return e.Err
}
