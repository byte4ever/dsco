package walker

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_stackEmbed_pushToStack(t *testing.T) {
	t.Parallel()

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			var s stackEmbed

			type Sub1 struct {
				X *int
			}

			type Sub2 struct {
				X *int
			}

			type Root struct {
				Sub1
				S  *Sub2
				aa *int
				A  *float32
			}

			var k Root

			vt := reflect.TypeOf(k)

			err := s.pushToStack(
				nil,
				0,
				"some-path",
				vt,
			)

			require.NoError(t, err)

			require.Equal(
				t, stackEmbed{
					{
						path:  "some-path",
						index: []int{3},
						field: vt.Field(3),
						depth: 0,
						order: 0,
					},
					{
						path:  "some-path",
						index: []int{2},
						field: vt.Field(2),
						depth: 0,
						order: 0,
					},
					{
						path:  "some-path",
						index: []int{1},
						field: vt.Field(1),
						depth: 0,
						order: 0,
					},
					{
						path:  "some-path",
						index: []int{0},
						field: vt.Field(0),
						depth: 0,
						order: 0,
					},
				}, s,
			)

		},
	)

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			var s stackEmbed

			type Sub1 struct {
				X *int
			}

			type Sub2 struct {
				X *int
			}

			type Root struct {
				Sub1
				S  *Sub2
				aa *int
				A  *float32
			}

			var k Root

			vt := reflect.TypeOf(k)

			err := s.pushToStack(
				[]int{11},
				0,
				"some-path",
				vt,
			)

			require.NoError(t, err)

			for _, embedded := range s {
				fmt.Println(*embedded)
			}

			require.Equal(
				t, stackEmbed{
					{
						path:  "some-path",
						index: []int{11, 3},
						field: vt.Field(3),
						depth: 0,
						order: 0,
					},
					{
						path:  "some-path",
						index: []int{11, 2},
						field: vt.Field(2),
						depth: 0,
						order: 0,
					},
					{
						path:  "some-path",
						index: []int{11, 1},
						field: vt.Field(1),
						depth: 0,
						order: 0,
					},
					{
						path:  "some-path",
						index: []int{11, 0},
						field: vt.Field(0),
						depth: 0,
						order: 0,
					},
				}, s,
			)

		},
	)

	t.Run(
		"error",
		func(t *testing.T) {
			t.Parallel()

			type Root struct {
				aa *int
				A  *float32
			}

			var k *Root

			vt := reflect.TypeOf(k)

			s := stackEmbed{}

			err := s.pushToStack(
				nil,
				0,
				"some-path",
				vt,
			)

			require.ErrorIs(t, err, ErrInvalidEmbedded)
			require.Len(t, s, 0)
		},
	)
}

func Test_getVisibleFieldList(t *testing.T) {
	t.Parallel()

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			type TS1 struct{}

			type Sub2 struct {
				Y2         *int
				notVisible *int
				Y1         *TS1
			}

			var vSub2 Sub2
			vSub2Type := reflect.TypeOf(vSub2)

			type Sub1 struct {
				X2 *int
				Sub2
				notVisible *int
				X1         *int
			}

			var vSub1 Sub1
			vSub1Type := reflect.TypeOf(vSub1)

			type Sub0 struct {
				Q *int
				Sub1
				notVisible *int
				W          *int
			}

			var vSub0 Sub0
			vSub0Type := reflect.TypeOf(vSub0)

			type Root struct {
				A *float32
				Sub0
				notVisible *int
				B          *int
			}

			var k Root

			vt := reflect.TypeOf(&k)

			items, errs := getVisibleFieldList(
				"",
				vt,
			)

			require.Nil(t, errs)

			for _, item := range items[:] {
				fmt.Println(*item)
			}

			require.Equal(
				t, elems{
					{
						path:  "",
						index: []int{0},
						field: vt.Elem().Field(0),
						depth: 0,
						order: 0,
					},
					{
						path:  "Sub0",
						index: []int{1, 0},
						field: vSub0Type.Field(0),
						depth: 1,
						order: 1,
					},
					{
						path:  "Sub0.Sub1",
						index: []int{1, 1, 0},
						field: vSub1Type.Field(0),
						depth: 2,
						order: 2,
					},
					{
						path:  "Sub0.Sub1.Sub2",
						index: []int{1, 1, 1, 0},
						field: vSub2Type.Field(0),
						depth: 3,
						order: 3,
					},
					{
						path:  "Sub0.Sub1.Sub2",
						index: []int{1, 1, 1, 2},
						field: vSub2Type.Field(2),
						depth: 3,
						order: 4,
					},
					{
						path:  "Sub0.Sub1",
						index: []int{1, 1, 3},
						field: vSub1Type.Field(3),
						depth: 2,
						order: 5,
					},
					{
						path:  "Sub0",
						index: []int{1, 3},
						field: vSub0Type.Field(3),
						depth: 1,
						order: 6,
					},
					{
						path:  "",
						index: []int{3},
						field: vt.Elem().Field(3),
						depth: 0,
						order: 7,
					},
				}, items,
			)
		},
	)

	t.Run(
		"detecting invalid embedded struct",
		func(t *testing.T) {
			t.Parallel()

			type Sub2 struct {
			}

			type Sub11 struct{}

			type Sub1 struct {
				X *int
				*Sub11
			}

			var vSub1 Sub1
			vSub1Type := reflect.TypeOf(vSub1)

			type Sub0 struct {
			}

			type Root struct {
				*Sub0
				Sub1
				*Sub2
			}

			var k Root

			vt := reflect.TypeOf(&k)

			gotElems, errs := getVisibleFieldList(
				"",
				vt,
			)

			require.Len(t, gotElems, 1)

			require.Equal(
				t,
				elems{
					{
						path:  "Sub1",
						index: []int{1, 0},
						field: vSub1Type.Field(0),
						depth: 1,
						order: 0,
					},
				},
				gotElems,
			)

			require.Len(t, errs, 3)
			require.ErrorIs(t, errs[0], ErrInvalidEmbedded)
			require.ErrorContains(t, errs[0], "Sub0")
			require.ErrorIs(t, errs[1], ErrInvalidEmbedded)
			require.ErrorContains(t, errs[1], "Sub1.Sub11")
			require.ErrorIs(t, errs[2], ErrInvalidEmbedded)
			require.ErrorContains(t, errs[2], "Sub2")
		},
	)

	t.Run(
		"detecting field name collision different depth",
		func(t *testing.T) {
			t.Parallel()

			type Sub11 struct {
				X1 *int
			}

			type Sub1 struct {
				Sub11
				Sub1X *int
			}

			type Root struct {
				Sub1
				X1 *int
			}

			var k Root

			vt := reflect.TypeOf(&k)

			_, errs := getVisibleFieldList(
				"",
				vt,
			)

			require.Len(t, errs, 1)

			require.ErrorIs(t, errs[0], ErrFieldNameCollision)
			require.ErrorContains(t, errs[0], "\"X1\"")
			require.ErrorContains(t, errs[0], "\"Sub1.Sub11.X1\"")
		},
	)

	t.Run(
		"detecting field name collision same depth",
		func(t *testing.T) {
			t.Parallel()

			type Sub1 struct {
				X *int
			}

			type Sub2 struct {
				X *float32
			}

			type Root struct {
				Sub1
				Sub2
			}

			var k Root

			vt := reflect.TypeOf(&k)

			_, errs := getVisibleFieldList(
				"",
				vt,
			)

			require.Len(t, errs, 1)
			require.ErrorIs(t, errs[0], ErrFieldNameCollision)
			require.ErrorContains(t, errs[0], "\"Sub1.X\"")
			require.ErrorContains(t, errs[0], "\"Sub2.X\"")
		},
	)

	t.Run(
		"detecting field name collision deep",
		func(t *testing.T) {
			t.Parallel()

			type Sub11 struct {
				X *int
			}

			type Sub1 struct {
				Sub11
			}

			type Sub22 struct {
				X *float32
			}

			type Sub2 struct {
				Sub22
			}

			type Root struct {
				Sub1
				Sub2
			}

			var k Root

			vt := reflect.TypeOf(&k)

			_, errs := getVisibleFieldList(
				"",
				vt,
			)

			require.Len(t, errs, 1)
			require.ErrorIs(t, errs[0], ErrFieldNameCollision)
			require.ErrorContains(t, errs[0], "\"Sub1.Sub11.X\"")
			require.ErrorContains(t, errs[0], "\"Sub2.Sub22.X\"")
		},
	)
}