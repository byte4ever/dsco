package model

import (
	"errors"
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
