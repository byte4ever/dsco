package strukt

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
)

type mapKeyI map[string]interface{}

var errMocked = errors.New("mocked error")

type T1Root struct {
	A *int
	B *float64
}

func TestProvide(t *testing.T) {
	t.Parallel()

	t.Run(
		"support empty struct in root struct",
		func(t *testing.T) {
			t.Parallel()

			type Root struct {
				Key1 *float64
				Key2 *int
				Key3 *string
			}

			val1 := dsco.R(123.423)
			val3 := dsco.R("Haha")
			binder, err := NewBinder(
				&Root{
					Key1: val1,
					Key3: val3,
				},
			)

			require.NoError(t, err)
			binder.checkValues(
				t,
				mapKeyI{
					"key1": val1,
					"key3": val3,
				},
			)
		},
	)

	t.Run(
		"detect unsupported types in root struct",
		func(t *testing.T) {
			t.Parallel()

			type Root struct {
				Key1    *float64
				Key2    *int
				Key3    *string
				Invalid int
			}

			_, err := NewBinder(
				&Root{},
			)

			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "invalid")
			require.ErrorContains(t, err, "int")

		},
	)

	t.Run(
		"support empty struct in sub struct",
		func(t *testing.T) {
			t.Parallel()

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

			val1 := dsco.R(123.423)
			val3 := dsco.R("Haha")
			binder, err := NewBinder(
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
			binder.checkValues(
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
		"detect unsupported types in sub struct",
		func(t *testing.T) {
			t.Parallel()

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

			_, err := NewBinder(
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
		"detect unsupported pointer types in sub struct",
		func(t *testing.T) {
			t.Parallel()

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

			_, err := NewBinder(
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
		"detect embedded struct",
		func(t *testing.T) {
			t.Parallel()

			type Embedded struct {
				KEY1 *float64
			}

			type LeafType struct {
				*Embedded
				KEY2 *float64
				KEY3 *int
				KEY4 *string
			}

			val1 := dsco.R(1.124)
			val2 := dsco.R(123.423)
			val3 := dsco.R(123)
			val4 := dsco.R("Haha")

			binder, err := NewBinder(
				&LeafType{
					Embedded: &Embedded{
						KEY1: val1,
					},
					KEY2: val2,
					KEY3: val3,
					KEY4: val4,
				},
			)

			require.NoError(t, err)
			binder.checkValues(
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
		"support yaml tag in root struct",
		func(t *testing.T) {
			t.Parallel()

			type LeafType struct {
				KEY2 *float64
				KEY3 *int
				KEY4 *string `yaml:"yaml_key"`
			}

			val2 := dsco.R(123.423)
			val3 := dsco.R(123)
			val4 := dsco.R("Haha")

			binder, err := NewBinder(
				&LeafType{
					KEY2: val2,
					KEY3: val3,
					KEY4: val4,
				},
			)

			require.NoError(t, err)
			binder.checkValues(
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
		"support yaml tag in sub struct",
		func(t *testing.T) {
			t.Parallel()

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

			val2 := dsco.R(123.423)
			val3 := dsco.R(123)
			val4 := dsco.R("Haha")

			binder, err := NewBinder(
				&RootType{
					Sub: &SubType{
						KEY2: val2,
						KEY3: val3,
						KEY4: val4,
					},
					KEY2: val2,
					KEY3: val3,
					KEY4: val4,
				},
			)

			require.NoError(t, err)
			binder.checkValues(
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
		"time type support",
		func(t *testing.T) {
			t.Parallel()

			type LeafType struct {
				KEY1 *time.Duration
				KEY2 *time.Time
				KEY3 *uint32
			}

			val1 := dsco.R(13 * time.Minute)
			val2 := dsco.R(time.Now())
			val3 := dsco.R(uint32(34))

			binder, err := NewBinder(
				&LeafType{
					KEY1: val1,
					KEY2: val2,
					KEY3: val3,
				},
			)

			require.NoError(t, err)
			binder.checkValues(
				t,
				mapKeyI{
					"key1": val1,
					"key2": val2,
					"key3": val3,
				},
			)
		},
	)

	t.Run(
		"slice type support",
		func(t *testing.T) {
			t.Parallel()

			type LeafType struct {
				KEY1 []int
			}

			val1 := []int{1, 2, 3}
			v := &LeafType{
				KEY1: val1,
			}

			b, err := NewBinder(
				v,
			)

			require.NoError(t, err)
			b.checkValues(
				t,
				mapKeyI{
					"key1": val1,
				},
			)
		},
	)
}

func (b *Binder) checkValues(
	t *testing.T,
	expectedKI mapKeyI,
) {
	t.Helper()

	ki := make(mapKeyI, len(b.entries))

	for k, e := range b.entries {
		require.False(t, e.Value.IsNil())
		ki[k] = e.Value.Interface()
	}

	require.Equal(t, expectedKI, ki)
}

func TestBinder_GetPostProcessErrors(t *testing.T) {
	t.Parallel()

	b := &Binder{}
	require.Nil(t, b.GetPostProcessErrors())
}

func TestProvideFromInterfaceProvider(t *testing.T) {
	t.Parallel()

	t.Run(
		"interface provider failure",
		func(t *testing.T) {
			t.Parallel()

			mip := NewMockInterfaceProvider(t)
			mip.On("GetInterface").Once().Return(nil, errMocked)

			provider, err := ProvideFromInterfaceProvider(mip)
			require.ErrorIs(t, err, errMocked)
			require.Nil(t, provider)
		},
	)

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			mip := NewMockInterfaceProvider(t)
			valA := dsco.R(123)
			valB := dsco.R(999.999)
			mip.On("GetInterface").
				Once().
				Return(
					&T1Root{
						A: valA,
						B: valB,
					}, nil,
				)

			provider, err := ProvideFromInterfaceProvider(mip)
			require.NoError(t, err)
			require.NotNil(t, provider)
			provider.checkValues(
				t, mapKeyI{
					"a": valA,
					"b": valB,
				},
			)
		},
	)
}

func TestBinder_Bind(t *testing.T) {
	t.Parallel()

	key := "a"

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			v := dsco.R(123.321)
			vValue := reflect.ValueOf(v)

			var k *float64
			vTargetValue := reflect.ValueOf(k)

			b := &Binder{
				entries: entries{
					key: &entry{
						Value: vValue,
					},
				},
			}

			o, keyOut, succeed, outVal, err := b.Bind(key, true, vTargetValue)
			require.NoError(t, err)
			require.Equal(t, id, o)
			require.True(t, succeed)
			require.Equal(t, key, keyOut)
			k = outVal.Interface().(*float64)
			require.Equal(t, v, k)
		},
	)

	t.Run(
		"key not found",
		func(t *testing.T) {
			t.Parallel()

			v := dsco.R(123.321)
			vValue := reflect.ValueOf(v)

			var k *float64
			vTargetValue := reflect.ValueOf(k)

			b := &Binder{
				entries: entries{
					key: &entry{
						Value: vValue,
					},
				},
			}

			invalidKey := "not_existing"
			o, keyOut, succeed, outVal, err := b.Bind(invalidKey, true, vTargetValue)
			require.NoError(t, err)
			require.Equal(t, id, o)
			require.False(t, succeed)
			require.Equal(t, invalidKey, keyOut)
			require.Equal(t, reflect.Value{}, outVal)
		},
	)

	t.Run(
		"type mismatch",
		func(t *testing.T) {
			t.Parallel()

			v := dsco.R(123.321)
			vValue := reflect.ValueOf(v)

			var k *int
			vTargetValue := reflect.ValueOf(k)

			b := &Binder{
				entries: entries{
					key: &entry{
						Value: vValue,
					},
				},
			}

			o, keyOut, succeed, outVal, err := b.Bind(key, true, vTargetValue)
			require.ErrorIs(t, err, ErrTypeMismatch)
			require.ErrorContains(t, err, "*float64 to type *int")
			require.Equal(t, id, o)
			require.False(t, succeed)
			require.Equal(t, key, keyOut)
			require.Equal(t, reflect.Value{}, outVal)
		},
	)
}
