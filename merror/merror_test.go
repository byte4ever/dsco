package merror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var errMocked1 = errors.New("mocked error 1")
var errMocked2 = errors.New("mocked error 2")
var errMocked3 = errors.New("mocked error 3")

func TestError_Error(t *testing.T) {
	{
		k := MError{
			errMocked1,
			errMocked2,
		}

		expectedString := `mocked error 1
mocked error 2`

		require.Equal(t, expectedString, k.Error())
	}
	{
		k := MError{
			errMocked1,
		}

		expectedString := `mocked error 1`

		require.Equal(t, expectedString, k.Error())
	}
	{
		k := MError{}

		expectedString := ``

		require.Equal(t, expectedString, k.Error())
	}
	{
		k := MError(nil)

		expectedString := ``

		require.Equal(t, expectedString, k.Error())
	}
}

func TestError_Is(t *testing.T) {
	var e MError

	require.ErrorIs(t, e, Err)
}

type dummyError struct{}

func (d dummyError) Error() string {
	panic("implement me")
}

func TestError_As(t *testing.T) {
	te := newMockTestError(t)
	te.
		On("Error").
		Return("mocked1").
		Once()

	k := MError{
		errMocked1,
		te,
		errMocked3,
	}

	var toFind testError

	require.ErrorAs(t, k, &toFind)
	require.ErrorContains(t, toFind, "mocked1")

	var toFindFailure dummyError

	require.False(t, errors.As(k, &toFindFailure))
}

type RootError struct {
	MError
}

var ErrRoot = errors.New("")

func (m RootError) Is(err error) bool {
	return errors.Is(err, ErrRoot)
}

type SubError struct {
	MError
}

var ErrSub = errors.New("")

func (m SubError) Is(err error) bool {
	return errors.Is(err, ErrSub)
}

func TestError_CascadingEffect(t *testing.T) {
	te := newMockTestError(t)
	te.
		On("Error").
		Return("mocked1").
		Once()

	se := SubError{
		MError: MError{errMocked1, te, errMocked2},
	}

	re := RootError{
		MError: MError{
			errMocked3, se, errMocked2,
		},
	}

	toTest := MError{
		errMocked1, re, errMocked2,
	}

	var root RootError

	require.ErrorAs(t, toTest, &root)
	require.Equal(t, re, root)

	var sub SubError

	require.ErrorAs(t, toTest, &sub)
	require.Equal(t, se, sub)

	var toFind testError

	require.ErrorAs(t, toTest, &toFind)
	require.ErrorContains(t, toFind, "mocked1")

	var toFindFailure dummyError

	require.False(t, errors.As(toTest, &toFindFailure))
}

func TestError_Add(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name        string
		args        args
		beforeState MError
		afterState  MError
	}{
		{
			name: "",
			args: args{
				err: errMocked1,
			},
			beforeState: nil,
			afterState:  MError{errMocked1},
		},
		{
			name: "",
			args: args{
				err: errMocked2,
			},
			beforeState: MError{errMocked1},
			afterState:  MError{errMocked1, errMocked2},
		},
	}
	for _, tt := range tests {
		tt.beforeState.Add(tt.args.err)
		require.Equal(t, tt.afterState, tt.beforeState)
	}
}

func TestMError_None(t *testing.T) {
	e := MError{}
	require.True(t, e.None())

	e2 := make(MError, 0)
	require.True(t, e2.None())
}

func TestMError_Count(t *testing.T) {
	var e MError

	require.Equal(t, 0, e.Count())

	e2 := MError{}

	require.Equal(t, 0, e2.Count())
}
