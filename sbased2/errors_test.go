package sbased2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinder_Errors(t *testing.T) {
	t.Parallel()

	type fields struct {
		values stringValues
	}

	tests := []struct {
		name             string
		fields           fields
		wantErrorIs      []error
		wantErrorContain [][]string
	}{
		{
			name: "no error",
			fields: fields{
				values: stringValues{
					"key1": {
						location: "loc-key1",
						state:    used,
					},
					"key2": {
						location: "loc-key2",
						state:    used,
					},
				},
			},
		},
		{
			name:   "no error no keys",
			fields: fields{},
		},
		{
			name: "some errors",
			fields: fields{
				values: stringValues{
					"key1": {
						location: "loc-key1",
						state:    used,
					},
					"key2": {
						location: "loc-key2",
						state:    unbounded,
					},
					"key3": {
						location: "loc-key3",
						state:    used,
					},
					"key4": {
						location: "loc-key4",
						state:    unused,
					},
					"key5": {
						location: "loc-key5",
						state:    used,
					},
				},
			},
			wantErrorIs:      []error{ErrUnboundKey, ErrOverriddenKey},
			wantErrorContain: [][]string{{"loc-key2"}, {"loc-key4"}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				// test invariant
				require.Equalf(
					t,
					len(tt.wantErrorIs),
					len(tt.wantErrorContain),
					"malformed test "+
						"got len(errorIs)=%v and len(wantErrorContain)=%v",
					len(tt.wantErrorIs),
					len(tt.wantErrorContain),
				)

				binder := &Binder{
					values: tt.fields.values,
				}

				errs := binder.Errors()

				require.Len(t, errs, len(tt.wantErrorIs))
				for i, err := range errs {
					require.ErrorIs(t, err, tt.wantErrorIs[i])

					for _, contain := range tt.wantErrorContain[i] {
						require.ErrorContainsf(
							t,
							err,
							contain,
							"error <%v> does not contain %q",
							err,
							contain,
						)
					}
				}
			},
		)
	}
}
