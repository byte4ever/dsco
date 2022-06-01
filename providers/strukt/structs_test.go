package strukt

import (
	"crypto"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/utils/hash"
)

type mapKeyI map[string]interface{}

func Test(t *testing.T) {
	t.Run(
		"support empty struct in root struct", func(t *testing.T) {
			type Root struct {
				Key1 *float64
				Key2 *int
				Key3 *string
			}

			val1 := dsco.V(123.423)
			val3 := dsco.V("Haha")
			b, err := Provide(
				&Root{
					Key1: val1,
					Key3: val3,
				},
			)

			require.NoError(t, err)
			b.checkValues(
				t,
				mapKeyI{
					"key1": val1,
					"key3": val3,
				},
			)
		},
	)

	t.Run(
		"detect unsupported types in root struct", func(t *testing.T) {
			type Root struct {
				Key1    *float64
				Key2    *int
				Key3    *string
				Invalid int
			}

			_, err := Provide(
				&Root{},
			)

			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "invalid")
			require.ErrorContains(t, err, "int")

		},
	)

	t.Run(
		"support empty struct in sub struct", func(t *testing.T) {
			type Sub struct {
				Key1 *float64
				Key2 *int
				Key3 *string
			}

			type Root struct {
				Sub  *Sub
				Key1 *float64
				Key2 *int
				Key3 *string
			}

			val1 := dsco.V(123.423)
			val3 := dsco.V("Haha")
			b, err := Provide(
				&Root{
					Sub: &Sub{
						Key1: val1,
						Key3: val3,
					},
					Key1: val1,
					Key3: val3,
				},
			)

			require.NoError(t, err)
			b.checkValues(
				t,
				mapKeyI{
					"key1":     val1,
					"key3":     val3,
					"sub-key1": val1,
					"sub-key3": val3,
				},
			)
		},
	)

	t.Run(
		"detect unsupported types in sub struct", func(t *testing.T) {
			type Sub struct {
				Key1    *float64
				Key2    *int
				Key3    *string
				Invalid int
			}

			type Root struct {
				Sub  *Sub
				Key1 *float64
				Key2 *int
				Key3 *string
			}

			_, err := Provide(
				&Root{
					Sub: &Sub{},
				},
			)
			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "sub-invalid")
			require.ErrorContains(t, err, "int")
		},
	)

	t.Run(
		"detect unsupported pointer types in sub struct", func(t *testing.T) {
			type Sub struct {
				Key1    *float64
				Key2    *int
				Key3    *string
				Invalid *map[string]string
			}

			type Root struct {
				Sub  *Sub
				Key1 *float64
				Key2 *int
				Key3 *string
			}

			_, err := Provide(
				&Root{
					Sub: &Sub{},
				},
			)

			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "sub-invalid")
			require.ErrorContains(t, err, "*map[string]string")
		},
	)

	t.Run(
		"detect embedded struct", func(t *testing.T) {

			type Embedded struct {
				KEY1 *float64
			}

			type LeafType struct {
				*Embedded
				KEY2 *float64
				KEY3 *int
				KEY4 *string
			}

			val1 := dsco.V(1.124)
			val2 := dsco.V(123.423)
			val3 := dsco.V(123)
			val4 := dsco.V("Haha")

			v := &LeafType{
				Embedded: &Embedded{
					KEY1: val1,
				},
				KEY2: val2,
				KEY3: val3,
				KEY4: val4,
			}

			b, err := Provide(
				v,
			)

			require.NoError(t, err)
			b.checkValues(
				t,
				mapKeyI{
					"embedded-key1": val1,
					"key2":          val2,
					"key3":          val3,
					"key4":          val4,
				},
			)
		},
	)

	t.Run(
		"support yaml tag in root struct", func(t *testing.T) {

			type LeafType struct {
				KEY2 *float64
				KEY3 *int
				KEY4 *string `yaml:"yaml_key"`
			}

			val2 := dsco.V(123.423)
			val3 := dsco.V(123)
			val4 := dsco.V("Haha")

			v := &LeafType{
				KEY2: val2,
				KEY3: val3,
				KEY4: val4,
			}

			b, err := Provide(
				v,
			)

			require.NoError(t, err)
			b.checkValues(
				t,
				mapKeyI{
					"key2":     val2,
					"key3":     val3,
					"yaml_key": val4,
				},
			)
		},
	)

	t.Run(
		"support yaml tag in sub struct", func(t *testing.T) {

			type SubType struct {
				KEY2 *float64
				KEY3 *int
				KEY4 *string `yaml:"yaml_key"`
			}

			type RootType struct {
				Sub  *SubType
				KEY2 *float64
				KEY3 *int
				KEY4 *string
			}

			val2 := dsco.V(123.423)
			val3 := dsco.V(123)
			val4 := dsco.V("Haha")

			v := &RootType{
				Sub: &SubType{
					KEY2: val2,
					KEY3: val3,
					KEY4: val4,
				},
				KEY2: val2,
				KEY3: val3,
				KEY4: val4,
			}

			b, err := Provide(
				v,
			)

			require.NoError(t, err)
			b.checkValues(
				t,
				mapKeyI{
					"key2":         val2,
					"key3":         val3,
					"key4":         val4,
					"sub-key2":     val2,
					"sub-key3":     val3,
					"sub-yaml_key": val4,
				},
			)
		},
	)

	t.Run(
		"duration hash and date properly handled", func(t *testing.T) {

			type LeafType struct {
				KEY1 *time.Duration
				KEY2 *time.Time
				KEY3 *hash.Hash
			}

			val1 := dsco.V(13 * time.Minute)
			val2 := dsco.V(time.Now())
			val3 := dsco.V(hash.Hash(crypto.SHA256))

			v := &LeafType{
				KEY1: val1,
				KEY2: val2,
				KEY3: val3,
			}

			b, err := Provide(
				v,
			)

			require.NoError(t, err)
			b.checkValues(
				t,
				mapKeyI{
					"key1": val1,
					"key2": val2,
					"key3": val3,
				},
			)
		},
	)
}

func (ks *Binder) checkValues(
	t *testing.T,
	expectedKI mapKeyI,
) {
	t.Helper()

	ki := make(mapKeyI, len(ks.values))

	for k, e := range ks.values {
		require.False(t, e.Value.IsNil())
		ki[k] = e.Value.Interface()
	}

	require.Equal(t, expectedKI, ki)
}

func TestBinder_GetPostProcessErrors(t *testing.T) {
	b := &Binder{}
	require.Nil(t, b.GetPostProcessErrors())
}
