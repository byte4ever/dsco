package sbased2

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
)

func TestBinder_BindAliasCollision(t *testing.T) {
	t.Parallel()

	keyName := "key"
	binder := &Binder{
		internalOpts: internalOpts{
			aliases: map[string]string{
				keyName: "",
			},
		},
		values: nil,
	}

	v := binder.Bind(keyName, nil)
	require.NotNil(t, v)
	require.ErrorIs(t, v.Error, ErrAliasCollision)
	require.ErrorContains(t, v.Error, keyName)
	require.False(t, v.HasValue())
	require.Equal(t, v.Location, "")
}

func TestBinder_BindNoKeyFound(t *testing.T) {
	t.Parallel()

	keyName := "key"
	binder := &Binder{}

	v := binder.Bind(keyName, nil)
	require.NotNil(t, v)
	require.NoError(t, v.Error, ErrAliasCollision)
	require.False(t, v.HasValue())
	require.Equal(t, v.Location, "")
}

func TestBinder_Bind(t *testing.T) {
	t.Parallel()

	type fields struct {
		internalOpts internalOpts
		values       stringValues
	}

	type args struct {
		dstType reflect.Type
		key     string
	}

	tests := []struct {
		args              args
		fields            fields
		wantError         error
		wantValue         reflect.Value
		name              string
		wantLocation      string
		wantErrorContains []string
		wantState         bindState
	}{
		{
			name: "pointer",
			fields: fields{
				values: stringValues{
					"key1": {
						location: "loc-key1",
						value:    "123",
					},
				},
			},
			args: args{
				key:     "key1",
				dstType: reflect.TypeOf(dsco.R(0)),
			},
			wantValue:         reflect.ValueOf(dsco.R(123)),
			wantError:         nil,
			wantErrorContains: nil,
			wantLocation:      "loc-key1",
			wantState:         unused,
		},
		{
			name: "slice",
			fields: fields{
				values: stringValues{
					"key1": {
						location: "loc-key1",
						value:    "[4,3,2,1]",
					},
				},
			},
			args: args{
				key:     "key1",
				dstType: reflect.TypeOf([]int{}),
			},
			wantValue:         reflect.ValueOf([]int{4, 3, 2, 1}),
			wantError:         nil,
			wantErrorContains: nil,
			wantLocation:      "loc-key1",
			wantState:         unused,
		},
		{
			name: "parse error",
			fields: fields{
				values: stringValues{
					"key1": {
						location: "loc-key1",
						value:    "asdf",
					},
				},
			},
			args: args{
				key:     "key1",
				dstType: reflect.TypeOf(dsco.R(0)),
			},
			wantValue:         reflect.Value{},
			wantError:         ErrParse,
			wantErrorContains: []string{"loc-key1"},
			wantLocation:      "loc-key1",
			wantState:         unbounded,
		},
		{
			name: "slice parse error",
			fields: fields{
				values: stringValues{
					"key1": {
						location: "loc-key1",
						value:    "[4,3,2",
					},
				},
			},
			args: args{
				key:     "key1",
				dstType: reflect.TypeOf([]int{}),
			},
			wantValue:         reflect.Value{},
			wantError:         ErrParse,
			wantErrorContains: []string{"loc-key1"},
			wantLocation:      "loc-key1",
			wantState:         unbounded,
		},
		{
			name: "type error",
			fields: fields{
				values: stringValues{
					"key1": {
						location: "loc-key1",
						value:    "123",
					},
				},
			},
			args: args{
				key:     "key1",
				dstType: reflect.TypeOf(4),
			},
			wantValue:         reflect.Value{},
			wantError:         ErrInvalidType,
			wantErrorContains: []string{"loc-key1"},
			wantLocation:      "loc-key1",
			wantState:         unbounded,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				binder := &Binder{
					internalOpts: tt.fields.internalOpts,
					values:       tt.fields.values,
				}

				attempt := binder.Bind(
					tt.args.key, tt.args.dstType,
				)

				if tt.wantError != nil {
					require.ErrorIs(t, attempt.Error, tt.wantError)

					for _, contain := range tt.wantErrorContains {
						require.ErrorContains(t, attempt.Error, contain)
					}

					require.False(t, attempt.HasValue())

					return
				}

				require.True(t, attempt.HasValue())

				// if tt.wantValue.Type().Kind() == reflect.Pointer {
				//
				// }
				require.Equal(
					t, tt.wantValue.Interface(),
					attempt.Value.Interface(),
				)

				require.Equal(t, tt.wantLocation, attempt.Location)
				require.NoError(t, attempt.Error)

				require.Equal(
					t,
					tt.wantState,
					binder.values[tt.args.key].state,
				)

			},
		)
	}
}
