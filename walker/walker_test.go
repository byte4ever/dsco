package walker

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
)

func Test_walker_walkSetter(t *testing.T) {
	t.Parallel()

	t.Run(
		"invalid type",
		func(t *testing.T) {
			t.Parallel()

			err := setStruct(
				nil,
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrNilInterface)
		},
	)

	t.Run(
		"not pointer",
		func(t *testing.T) {
			t.Parallel()

			err := setStruct(
				123,
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrExpectPointerOnStruct)
		},
	)

	t.Run(
		"not pointer on struct",
		func(t *testing.T) {
			t.Parallel()

			v := 10
			err := setStruct(
				&v,
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrExpectPointerOnStruct)
		},
	)

	t.Run(
		"field name collision same level",
		func(t *testing.T) {
			t.Parallel()

			type Sub2A struct {
				X *int
			}

			type Sub2B struct {
				X *int
			}

			type Sub1 struct {
				Sub2A
				Sub2B
			}

			type Root struct {
				S *Sub1
			}

			err := setStruct(
				&Root{},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrFieldNameCollision)
			require.ErrorContains(t, err, "S.Sub2A.X")
			require.ErrorContains(t, err, "S.Sub2B.X")
		},
	)

	t.Run(
		"field name collision different level",
		func(t *testing.T) {
			t.Parallel()

			type Sub3 struct {
				X *int
			}

			type Sub2A struct {
				Y *int
				Sub3
			}

			type Sub2B struct {
				X *int
			}

			type Sub1 struct {
				Sub2A
				Sub2B
			}

			type Root struct {
				S *Sub1
			}

			err := setStruct(
				&Root{},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrFieldNameCollision)
			require.ErrorContains(t, err, "S.Sub2A.Sub3.X")
			require.ErrorContains(t, err, "S.Sub2B.X")
		},
	)

	t.Run(
		"invalid field type first level",
		func(t *testing.T) {
			t.Parallel()

			type Root struct {
				KeyA int
			}

			err := setStruct(
				&Root{},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "int")
			require.ErrorContains(t, err, "KeyA")
		},
	)

	t.Run(
		"invalid field type deep level",
		func(t *testing.T) {
			t.Parallel()

			type Sub3 struct {
				Key3A int
			}

			type Sub2 struct {
				S3 *Sub3
			}

			type Sub1 struct {
				S2 *Sub2
			}

			type Root struct {
				S1 *Sub1
			}

			err := setStruct(
				&Root{},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "int")
			require.ErrorContains(t, err, "S1.S2.S3.Key3A")
		},
	)

	t.Run(
		"invalid field type embedded",
		func(t *testing.T) {
			t.Parallel()

			type Sub struct {
				KeyA int
			}

			type Root struct {
				Sub
			}

			err := setStruct(
				&Root{},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "int")
			require.ErrorContains(t, err, "KeyA")
		},
	)

	t.Run(
		"invalid embedded type",
		func(t *testing.T) {
			t.Parallel()

			type Sub struct {
				A *int
			}

			type Root struct {
				*Sub
			}

			err := setStruct(
				&Root{},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrEmbeddedPointer)
			require.ErrorContains(t, err, "Sub")
		},
	)

	t.Run(
		"invalid embedded type deep",
		func(t *testing.T) {
			t.Parallel()

			type Embedded struct {
				X *int
			}

			type Sub2 struct {
				*Embedded
			}

			type Sub1 struct {
				S2 *Sub2
			}

			type Root struct {
				S1 *Sub1
			}

			err := setStruct(
				&Root{},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrEmbeddedPointer)
			require.ErrorContains(t, err, "S1.S2.Embedded")
		},
	)

	t.Run(
		"invalid field type embedded deep level",
		func(t *testing.T) {
			t.Parallel()

			type Sub3 struct {
				Key3A int
			}

			type Sub2 struct {
				Sub3
			}

			type Sub1 struct {
				S2 *Sub2
			}

			type Root struct {
				S1 *Sub1
			}

			err := setStruct(
				&Root{},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "int")
			require.ErrorContains(t, err, "S1.S2.Key3A")
		},
	)

	t.Run(
		"detecting initialized field first level",
		func(t *testing.T) {
			t.Parallel()

			type Root struct {
				Initialized *int
			}

			err := setStruct(
				&Root{
					Initialized: dsco.R(10),
				},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrNotNilValue)
			require.ErrorContains(t, err, "Initialized")
		},
	)

	t.Run(
		"detecting initialized field struct first level",
		func(t *testing.T) {
			t.Parallel()

			type Sub3 struct {
				Initialized *int
			}

			type Sub2 struct {
				S3 *Sub3
			}

			type Sub1 struct {
				S2 *Sub2
			}

			type Root struct {
				InitializedS1 *Sub1
			}
			err := setStruct(
				&Root{
					InitializedS1: &Sub1{
						S2: &Sub2{
							S3: &Sub3{
								Initialized: dsco.R(10),
							},
						},
					},
				},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.ErrorIs(t, err, ErrNotNilValue)
			require.ErrorContains(t, err, "InitializedS1")
		},
	)

	t.Run(
		"manage slice first level",
		func(t *testing.T) {
			t.Parallel()

			type Root struct {
				S []int
			}
			err := setStruct(
				&Root{},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.NoError(t, err)
		},
	)

	t.Run(
		"manage slice deep level",
		func(t *testing.T) {
			t.Parallel()

			type Sub3 struct {
				Slice []int
			}

			type Sub2 struct {
				S3 *Sub3
			}

			type Sub1 struct {
				S2 *Sub2
			}

			type Root struct {
				S1 *Sub1
			}

			err := setStruct(
				&Root{},
				func(order int, path string, value *reflect.Value) error {
					return nil
				},
			)

			require.NoError(t, err)
		},
	)
}
