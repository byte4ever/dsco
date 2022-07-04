package cmdline

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var errMocked1 = errors.New("mocked err 1")
var errMocked2 = errors.New("mocked err 2")
var errMocked3 = errors.New("mocked err 3")

func TestParamError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		Positions []int
		Errs      []error
	}

	tests := []struct {
		name   string
		want   string
		fields fields
	}{
		{
			name: "multiple errors",
			fields: fields{
				Positions: []int{1, 21, 301},
				Errs: []error{
					errMocked1,
					errMocked2,
					errMocked3,
				},
			},
			want: "cmdline issues at positions #1, #21 and #301:" +
				" mocked err 1 / mocked err 2 / mocked err 3",
		},
		{
			name: "single error",
			fields: fields{
				Positions: []int{10},
				Errs: []error{
					errMocked1,
				},
			},
			want: "cmdline issue at position #10:" +
				" mocked err 1",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				errParam := &ParamError{
					Positions: tt.fields.Positions,
					Errs:      tt.fields.Errs,
				}

				errorString := errParam.Error()

				require.Equalf(
					t, tt.want, errorString,
					"error string %q does not match expected %q",
					errorString, tt.want,
				)
			},
		)
	}
}

func TestParamError_Error_Precondition(t *testing.T) {
	t.Parallel()

	t.Run(
		"internal length differ",
		func(t *testing.T) {
			t.Parallel()

			errParam := &ParamError{
				Positions: []int{1, 2, 3},
				Errs:      []error{nil},
			}
			require.Panics(
				t, func() {
					_ = errParam.Error()
				},
			)
		},
	)

	t.Run(
		"internal length are 0",
		func(t *testing.T) {
			t.Parallel()

			errParam := &ParamError{}
			require.Panics(
				t, func() {
					_ = errParam.Error()
				},
			)
		},
	)
}
