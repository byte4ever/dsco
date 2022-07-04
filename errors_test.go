package dsco

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

var errMocked1 = errors.New("error mocked 1")
var errMocked2 = errors.New("error mocked 2")

func TestInvalidInputError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"type int is not a valid pointer on struct",
		InvalidInputError{
			Type: reflect.TypeOf(10),
		}.Error(),
	)
}

func TestInvalidInputError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrInvalidInput)
	require.ErrorIs(t, InvalidInputError{}, ErrInvalidInput)
}

func TestCmdlineAlreadyUsedError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"cmdline already used in position #101",
		CmdlineAlreadyUsedError{
			Index: 101,
		}.Error(),
	)
}

func TestCmdlineAlreadyUsedError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrCmdlineAlreadyUsed)
	require.ErrorIs(t, CmdlineAlreadyUsedError{}, ErrCmdlineAlreadyUsed)
}

func TestDuplicateEnvPrefixError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"layer #101 has same prefix=PREFIX",
		DuplicateEnvPrefixError{
			Index:  101,
			Prefix: "PREFIX",
		}.Error(),
	)
}

func TestDuplicatAeEnvPrefixError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrDuplicateEnvPrefix)
	require.ErrorIs(t, DuplicateEnvPrefixError{}, ErrDuplicateEnvPrefix)
}

func TestDuplicateInputStructError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"struct layer #101 is using same pointer",
		DuplicateInputStructError{
			Index: 101,
		}.Error(),
	)
}

func TestDuplicateInputStructError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrDuplicateInputStruct)
	require.ErrorIs(t, DuplicateInputStructError{}, ErrDuplicateInputStruct)
}

func TestDuplicateStructIDError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"struct layer #101 is using same id=\"OTHER\"",
		DuplicateStructIDError{
			Index: 101,
			ID:    "OTHER",
		}.Error(),
	)
}

func TestDuplicateStructIDError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrDuplicateStructID)
	require.ErrorIs(t, DuplicateStructIDError{}, ErrDuplicateStructID)
}

func TestLayerErrors_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrLayer)
	require.ErrorIs(t, LayerErrors{}, ErrLayer)
}
