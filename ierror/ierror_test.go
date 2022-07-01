package ierror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var errMocked = errors.New("mocked error")

func TestIError_Error(t *testing.T) {
	require.Equal(
		t,
		"-info- #101: mocked error",
		IError{
			Index: 101,
			Info:  "-info-",
			Err:   errMocked,
		}.Error(),
	)
}

func TestIError_Unwrap(t *testing.T) {
	e := IError{
		Err: errMocked,
	}

	require.Equal(t, errMocked, e.Unwrap())
}
