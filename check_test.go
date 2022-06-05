package dsco

import (
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
				KEY5 *time.Time
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
		"detect invalid destination nil case", func(t *testing.T) {
			var v *T4Root
			err := checkStruct(v)
			require.ErrorIs(t, err, ErrInvalidDestination)
		},
	)

	t.Run(
		"detect invalid destination not pointer", func(t *testing.T) {
			var v int
			err := checkStruct(v)
			require.ErrorIs(t, err, ErrInvalidDestination)
		},
	)

	t.Run(
		"detect invalid destination not pointer on struct", func(t *testing.T) {
			var v *int
			err := checkStruct(v)
			require.ErrorIs(t, err, ErrInvalidDestination)
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

// ///////////////////////////////////////

// //////////////////////////////////////////////

// ////////////////////////////////////////

// ////////////////////////////////////////

// ////////////////////////////////////////////
