package walker

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/walker/fvalues"
	"github.com/byte4ever/dsco/walker/ifaces"
)

func TestNewStructBuilder(t *testing.T) {
	id := "id"

	v := 10

	sb, err := NewStructBuilder(v, id)
	require.NoError(t, err)
	require.NotNil(t, sb)
	require.Equal(t, id, sb.id)
	require.Equal(t, 10, sb.value.Interface())
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
			fvsOut := fvalues.FieldValues{
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
