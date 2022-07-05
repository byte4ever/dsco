package dsco

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/ifaces"
	"github.com/byte4ever/dsco/internal/fvalue"
)

func TestNewStructBuilder(t *testing.T) {
	t.Parallel()

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			id := "id"

			type Root struct {
				X *float32
				Y *float32
			}

			v := &Root{}

			sb, err := NewStructBuilder(v, id)
			require.NoError(t, err)
			require.NotNil(t, sb)
			require.Equal(t, id, sb.id)
			require.Equal(t, v, sb.value.Interface())
		},
	)

	t.Run(
		"nil input",
		func(t *testing.T) {
			t.Parallel()

			id := "id"

			sb, err := NewStructBuilder(nil, id)
			require.Nil(t, sb)
			require.ErrorIs(t, err, ErrNilInput)
		},
	)

	t.Run(
		"not a pointer",
		func(t *testing.T) {
			t.Parallel()

			id := "id"

			type Root struct {
				X *float32
				Y *float32
			}

			v := Root{}

			sb, err := NewStructBuilder(v, id)
			require.Nil(t, sb)

			var e InvalidInputError
			require.ErrorAs(t, err, &e)
			require.Equal(t, reflect.TypeOf(v), e.Type)
		},
	)

	t.Run(
		"not a pointer on struct",
		func(t *testing.T) {
			t.Parallel()

			id := "id"

			temp := 123
			v := &temp

			sb, err := NewStructBuilder(v, id)
			require.Nil(t, sb)

			var e InvalidInputError
			require.ErrorAs(t, err, &e)
			require.Equal(t, reflect.TypeOf(v), e.Type)
		},
	)
}

func TestStructBuilder_GetFieldValuesFrom(t *testing.T) {
	t.Parallel()

	t.Run(
		"invalid type name", func(t *testing.T) {
			t.Parallel()

			model := ifaces.NewMockModelInterface(t)
			model.
				On("TypeName").
				Return("other").
				Once()

			sb := &StructBuilder{
				value: reflect.ValueOf(10),
				id:    "",
			}

			gotVfs, err := sb.GetFieldValuesFrom(model)

			require.ErrorIs(t, err, ErrStructTypeDiffer)
			require.Nil(t, gotVfs)
		},
	)

	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()

			id := "id"
			value := reflect.ValueOf(10)
			fvsOut := fvalue.Values{
				uint(200): {
					Location: "loc",
				},
			}

			model := ifaces.NewMockModelInterface(t)
			model.
				On("TypeName").
				Return("int").
				Once()

			model.
				On("GetFieldValuesFor", id, value).
				Return(fvsOut).
				Once()

			sb := &StructBuilder{
				value: value,
				id:    id,
			}

			gotVfs, err := sb.GetFieldValuesFrom(model)

			require.NoError(t, err)
			require.Equal(t, fvsOut, gotVfs)
		},
	)
}
