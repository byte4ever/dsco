package sbased2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinder_Use(t *testing.T) {
	t.Parallel()

	type fields struct {
		values stringValues
	}

	type args struct {
		key string
	}

	tests := []struct {
		name             string
		fields           fields
		args             args
		wantErrorIs      error
		wantErrorContain []string
		wantState        bindState
	}{
		{
			name: "success",
			fields: fields{
				values: stringValues{
					"key1": {
						location: "loc-key1",
						state:    unused,
					},
				},
			},
			args: args{
				key: "key1",
			},
			wantState: used,
		},
		{
			name: "invalid key",
			fields: fields{
				values: stringValues{
					"key_oups": {
						location: "loc-key1",
						state:    unused,
					},
				},
			},
			args: args{
				key: "key1",
			},
			wantErrorIs:      ErrKeyNotFound,
			wantErrorContain: []string{"key1"},
			wantState:        unbounded,
		},
		{
			name: "invalid state",
			fields: fields{
				values: stringValues{
					"key1": {
						location: "loc-key1",
						state:    unbounded,
					},
				},
			},
			args: args{
				key: "key1",
			},
			wantErrorIs:      ErrNotUnused,
			wantErrorContain: []string{"loc-key1"},
			wantState:        unbounded,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				binder := &Binder{
					values: tt.fields.values,
				}

				err := binder.Use(tt.args.key)

				if tt.wantErrorIs != nil {
					require.ErrorIs(t, err, tt.wantErrorIs)

					for _, contain := range tt.wantErrorContain {
						require.ErrorContains(t, err, contain)
					}

					return
				}

				require.Equal(
					t,
					tt.wantState,
					binder.values[tt.args.key].state,
				)
			},
		)
	}
}
