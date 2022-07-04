package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/fvalues"
	"github.com/byte4ever/dsco/internal/plocation"
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
				aa *int //nolint:unused // don't care
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
				aa *int //nolint:unused // don't care
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
				aa *int //nolint:unused // don't care
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
				notVisible *int //nolint:unused // don't care
				Y1         *TS1
			}

			var vSub2 Sub2
			vSub2Type := reflect.TypeOf(vSub2)

			type Sub1 struct {
				X2 *int
				Sub2
				notVisible *int //nolint:unused // don't care
				X1         *int
			}

			var vSub1 Sub1
			vSub1Type := reflect.TypeOf(vSub1)

			type Sub0 struct {
				Q *int
				Sub1
				notVisible *int //nolint:unused // don't care
				W          *int
			}

			var vSub0 Sub0
			vSub0Type := reflect.TypeOf(vSub0)

			type Root struct {
				A *float32
				Sub0
				notVisible *int //nolint:unused // don't care
				B          *int
			}

			var k Root

			vt := reflect.TypeOf(&k)

			items, errs := getVisibleFieldList(
				"",
				vt,
			)

			require.Nil(t, errs)

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

			var e FieldNameCollisionError
			require.ErrorAs(t, errs[0], &e)
			require.Equal(
				t, FieldNameCollisionError{
					Path1: "X1",
					Path2: "Sub1.Sub11.X1",
				}, e,
			)
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

			var e FieldNameCollisionError
			require.ErrorAs(t, errs[0], &e)
			require.Equal(
				t,
				FieldNameCollisionError{
					Path1: "Sub2.X",
					Path2: "Sub1.X",
				},
				e,
			)
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
			var e FieldNameCollisionError
			require.ErrorAs(t, errs[0], &e)
			require.Equal(
				t, FieldNameCollisionError{
					Path1: "Sub2.Sub22.X",
					Path2: "Sub1.Sub11.X",
				}, e,
			)
		},
	)
}

func TestModel_TypeName(t *testing.T) {
	t.Parallel()

	m := &Model{
		typeName: "type-name",
	}
	require.Equal(t, "type-name", m.TypeName())
}

func TestNewModel(t *testing.T) {
	t.Parallel()

	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()

			type Root struct {
				X *float64
				Y *float64
			}

			m, err := NewModel(reflect.TypeOf(&Root{}))
			require.NoError(t, err)
			require.Equal(t, uint(2), m.fieldCount)
			require.NotNil(t, m.accelerator)
			require.NotNil(t, m.getList)
			require.Equal(t, 2, m.getList.Count())
		},
	)

	t.Run(
		"errors", func(t *testing.T) {
			t.Parallel()

			type Root struct {
				Y *float64
				X float64
			}

			m, err := NewModel(reflect.TypeOf(&Root{}))

			var e ModelError
			require.ErrorAs(t, err, &e)
			require.False(t, e.None())
			require.Nil(t, m)
		},
	)
}

func TestModel_ApplyOn(t *testing.T) {
	t.Parallel()

	fvs := make(fvalues.FieldValues, 10)

	getter := NewMockGetter(t)

	getList := NewMockGetListInterface(t)
	getList.On("ApplyOn", getter).Return(fvs, errMocked1).Once()

	m := &Model{
		getList: getList,
	}

	gotFvs, err := m.ApplyOn(getter)

	require.Equal(t, fvs, gotFvs)
	require.Equal(t, errMocked1, err)
}

func TestModel_Fill(t *testing.T) {
	t.Parallel()

	v := reflect.ValueOf(1)
	layers := []fvalues.FieldValues{nil, nil, nil}
	ploc := plocation.PathLocations{plocation.PathLocation{}}

	accelerator := NewMockNode(t)
	accelerator.
		On("Fill", v, layers).
		Return(ploc, errMocked1).
		Once()

	m := &Model{
		accelerator: accelerator,
	}

	gotPLoc, err := m.Fill(v, layers)

	require.Equal(t, ploc, gotPLoc)
	require.Equal(t, errMocked1, err)
}

func TestModelError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrModel)
	require.ErrorIs(t, ModelError{}, ErrModel)
}

func TestModel_GetFieldValuesFor(t *testing.T) {
	t.Parallel()

	id := "id"
	v := reflect.ValueOf(1)

	accelerator := NewMockNode(t)
	accelerator.
		On(
			"FeedFieldValues",
			id,
			mock.MatchedBy(
				func(ifvs fvalues.FieldValues) bool {
					return assert.Empty(t, ifvs)
				},
			),
			v,
		).
		Return().
		Once()

	m := &Model{
		fieldCount:  101,
		accelerator: accelerator,
	}

	gotFvs := m.GetFieldValuesFor(id, v)

	require.Empty(t, gotFvs)
}
