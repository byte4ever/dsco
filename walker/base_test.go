package walker

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBase_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		id int
	}

	val0 := reflect.ValueOf(0)
	val1 := reflect.ValueOf(1)
	val2 := reflect.ValueOf(2)

	tests := []struct {
		b            Base
		wantValue    *reflect.Value
		wantState    Base
		name         string
		wantLocation string
		args         args
	}{
		{
			name: "success",
			b: Base{
				0: {
					path:     "p0",
					location: "l0",
					value:    &val0,
				},
				1: {
					path:     "p1",
					location: "l1",
					value:    &val1,
				},
				2: {
					path:     "p2",
					location: "l2",
					value:    &val2,
				},
			},
			args: args{
				id: 0,
			},
			wantValue:    &val0,
			wantLocation: "l0",
			wantState: Base{
				1: {
					path:     "p1",
					location: "l1",
					value:    &val1,
				},
				2: {
					path:     "p2",
					location: "l2",
					value:    &val2,
				},
			},
		},
		{
			name: "not found",
			b: Base{
				0: {
					path:     "p0",
					location: "l0",
					value:    &val0,
				},
				1: {
					path:     "p1",
					location: "l1",
					value:    &val1,
				},
				2: {
					path:     "p2",
					location: "l2",
					value:    &val2,
				},
			},
			args: args{
				id: 11,
			},
			wantValue:    nil,
			wantLocation: "",
			wantState: Base{
				0: {
					path:     "p0",
					location: "l0",
					value:    &val0,
				},
				1: {
					path:     "p1",
					location: "l1",
					value:    &val1,
				},
				2: {
					path:     "p2",
					location: "l2",
					value:    &val2,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				value, location := tt.b.Get(tt.args.id)
				require.Equal(t, tt.wantValue, value)
				require.Equal(t, tt.wantLocation, location)
				require.Equal(t, tt.wantState, tt.b)
			},
		)
	}
}
