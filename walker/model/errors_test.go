package model

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

var errMocked1 = errors.New("error mocked 1")
var errMocked2 = errors.New("error mocked 2")

type MockedError1 struct{}

func (m MockedError1) Error() string {
	return "mocked error #1"
}

func checkAsMockedError1(
	t *testing.T,
	err error,
) {
	t.Helper()

	var me MockedError1
	require.ErrorAs(t, err, &me)
}

type MockedError2 struct{}

func (m MockedError2) Error() string {
	return "mocked error #1"
}

func checkAsMockedError2(
	t *testing.T,
	err error,
) {
	t.Helper()

	var me MockedError2
	require.ErrorAs(t, err, &me)
}

func TestFieldNameCollisionError_Error(t *testing.T) {
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
	require.Equal(
		t,
		"struct field p1 with unsupported type int",
		UnsupportedTypeError{
			Path: "p1",
			Type: reflect.TypeOf(10),
		}.Error(),
	)
}
