package model

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

var errMocked1 = errors.New("error mocked 1")

type Mocked1Error struct{}

func (Mocked1Error) Error() string {
	return "mocked error #1"
}

func checkAsMockedError1(
	t *testing.T,
	err error,
) {
	t.Helper()

	var me Mocked1Error

	require.ErrorAs(t, err, &me)
}

type Mocked2Error struct{}

func (Mocked2Error) Error() string {
	return "mocked error #1"
}

func checkAsMockedError2(
	t *testing.T,
	err error,
) {
	t.Helper()

	var me Mocked2Error

	require.ErrorAs(t, err, &me)
}

func TestFieldNameCollisionError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"field collision between p1 and p2",
		FieldNameCollisionError{
			Path1: "p1",
			Path2: "p2",
		}.Error(),
	)
}

func TestUnsupportedTypeError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"struct field p1 with unsupported type int",
		UnsupportedTypeError{
			Path: "p1",
			Type: reflect.TypeOf(10),
		}.Error(),
	)
}
