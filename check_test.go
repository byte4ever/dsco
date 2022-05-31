package dsco

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/utils/hash"
)

func Test_checkStruct(t *testing.T) {

	t.Run(
		"simple with embedded", func(t *testing.T) {
			type Embedded struct {
				KEY1 *float64
			}

			type LeafType struct {
				*Embedded
				KEY2 *float64
				KEY3 *int
				KEY4 *string
			}

			v := &LeafType{}

			err := checkStruct(v)
			require.NoError(t, err)
		},
	)

	t.Run(
		"simple with sub type", func(t *testing.T) {
			type Sub struct {
				KEY1 *float64
			}

			type Root struct {
				SubKey *Sub
				KEY2   *float64
				KEY3   *int
				KEY4   *string
			}

			v := &Root{}

			err := checkStruct(v)
			require.NoError(t, err)
		},
	)

	t.Run(
		"non empty struct case 1", func(t *testing.T) {
			type Sub struct {
				KEY1 *float64
			}

			type Root struct {
				SubKey *Sub
				KEY2   *float64
				KEY3   *int
				KEY4   *string
			}

			v := &Root{
				SubKey: &Sub{
					KEY1: V(123.31),
				},
			}

			err := checkStruct(v)
			require.ErrorIs(t, err, ErrRequireEmptyStruct)
			require.ErrorContains(t, err, "sub_key")
		},
	)

	t.Run(
		"non empty struct case 2", func(t *testing.T) {
			type Root struct {
				KEY2 *float64
				KEY3 *int
				KEY4 *string
			}

			v := &Root{
				KEY2: V(123123.123),
			}

			err := checkStruct(v)
			fmt.Println(err)
			require.ErrorIs(t, err, ErrRequireEmptyStruct)
		},
	)

	t.Run(
		"non empty struct case 3", func(t *testing.T) {
			type Root struct {
				KEY2 *time.Duration
			}

			v := &Root{
				KEY2: V(10 * time.Second),
			}

			err := checkStruct(v)
			require.ErrorIs(t, err, ErrRequireEmptyStruct)
		},
	)

	t.Run(
		"deep with sub type", func(t *testing.T) {

			type Sub3 struct {
				Key1        *float64
				KeyDuration *time.Duration
				KeyTime     *time.Duration
				KeyHash     *hash.Hash
			}

			type Sub2 struct {
				SubKey *Sub3
				KEY1   *float64
			}

			type Sub1 struct {
				SubKey *Sub2
				KEY1   *float64
			}

			type Root struct {
				SubKey *Sub1
				KEY2   *float64
				KEY3   *int
				KEY4   *string
			}

			v := &Root{}
			err := checkStruct(v)
			require.NoError(t, err)
		},
	)

	t.Run(
		"detect deep recursive", func(t *testing.T) {
			v := &T1Root{}
			err := checkStruct(v)
			require.ErrorIs(t, err, ErrRecursiveStruct)
			require.ErrorContains(t, err, "sub_key_root")
			require.ErrorContains(t, err, "sub_key_root-sub_key1-sub_key2-sub_key3")
			require.ErrorContains(t, err, "*dsco.T1Sub1")
		},
	)

	t.Run(
		"ensure dfs is working properly", func(t *testing.T) {
			v := &T2Root{}
			err := checkStruct(v)
			require.NoError(t, err, ErrRecursiveStruct)
		},
	)

	t.Run(
		"detect deep recursive with root", func(t *testing.T) {
			v := &T3Root{}
			err := checkStruct(v)
			require.ErrorIs(t, err, ErrRecursiveStruct)
			require.ErrorContains(t, err, "main struct")
			require.ErrorContains(t, err, "sub_key_root-sub_key1-sub_key2-cycle_root")
			require.ErrorContains(t, err, "*dsco.T3Root")
		},
	)

	t.Run(
		"detect invalid types ie no pointer", func(t *testing.T) {
			v := &T4Root{}
			err := checkStruct(v)
			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "key4")
			require.ErrorContains(t, err, "int")
		},
	)

	t.Run(
		"detect invalid types", func(t *testing.T) {
			v := &T5Root{}
			err := checkStruct(v)
			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "key4")
			require.ErrorContains(t, err, "map[string]string")
		},
	)

	t.Run(
		"detect invalid types with yaml field name", func(t *testing.T) {
			v := &T6Root{}
			err := checkStruct(v)
			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "renamed")
			require.ErrorContains(t, err, "map[string]string")
		},
	)

}

// ///////////////////////////////////

type T1Sub3 struct {
	SubKey3 *T1Sub1
	KEY1    *float64
}

type T1Sub2 struct {
	SubKey2 *T1Sub3
	KEY1    *float64
}

type T1Sub1 struct {
	SubKey1 *T1Sub2
	KEY1    *float64
}

type T1Root struct {
	SubKeyRoot *T1Sub1
	KEY2       *float64
	KEY3       *int
	KEY4       *string
}

// ///////////////////////////////////////

type T2Sub3 struct {
	KEY1 *float64
}

type T2Sub2 struct {
	SubKey2 *T2Sub3
	KEY1    *float64
}

type T2Sub1 struct {
	SubKey1 *T2Sub2
	SubKey2 *T2Sub2
	KEY1    *float64
}

type T2Root struct {
	SubKeyRoot *T2Sub1
	KEY2       *float64
	KEY3       *int
	KEY4       *string
}

// //////////////////////////////////////////////

type T3Sub3 struct {
	KEY1      *float64
	CycleRoot *T3Root
}

type T3Sub2 struct {
	SubKey2 *T3Sub3
	KEY1    *float64
}

type T3Sub1 struct {
	SubKey1 *T3Sub2
	KEY1    *float64
}

type T3Root struct {
	SubKeyRoot *T3Sub1
	KEY2       *float64
	KEY3       *int
	KEY4       *string
}

// ////////////////////////////////////////

type T4Root struct {
	KEY2 *float64
	KEY3 *int
	KEY4 int
}

// ////////////////////////////////////////

type T5Root struct {
	KEY2 *float64
	KEY3 *int
	KEY4 map[string]string
}

// ////////////////////////////////////////////
type T6Root struct {
	KEY2 *float64
	KEY3 *int
	KEY4 map[string]string `yaml:"renamed"`
}
