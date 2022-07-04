package model

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/fvalues"
	"github.com/byte4ever/dsco/plocation"
	"github.com/byte4ever/dsco/utils"
)

func TestStructNode_Fill(t *testing.T) {
	t.Parallel()

	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()

			fvs := []fvalues.FieldValues{}

			s0 := NewMockNode(t)
			s0.
				On(
					"Fill",
					mock.IsType(reflect.Value{}),
					fvs,
				).
				Return(
					plocation.PathLocations{
						plocation.PathLocation{
							UID:      0,
							Path:     "path.s0.A",
							Location: "loc-s0.A",
						},
						plocation.PathLocation{
							UID:      1,
							Path:     "path.s0.B",
							Location: "loc-s0.B",
						},
					},
					MockedError1{},
				).
				Once()

			s1 := NewMockNode(t)
			s1.
				On(
					"Fill",
					mock.IsType(reflect.Value{}),
					fvs,
				).
				Return(
					plocation.PathLocations{
						plocation.PathLocation{
							UID:      2,
							Path:     "path.s1",
							Location: "loc-s1",
						},
					},
					MockedError2{},
				).
				Once()

			type SType struct {
				A interface{}
				B interface{}
			}

			var i *SType

			v := reflect.ValueOf(&i).Elem()

			require.True(t, v.CanSet())

			n := &StructNode{
				Type: reflect.TypeOf(i),
				Index: IndexedSubNodes{
					&IndexedSubNode{
						Node:  s0,
						Index: []int{0},
					},
					&IndexedSubNode{
						Node:  s1,
						Index: []int{1},
					},
				},
			}

			ploc, err := n.Fill(v, fvs)

			// check errors are properly aggregate.
			checkAsMockedError1(t, err)
			checkAsMockedError2(t, err)

			// path locations are properly aggregate.
			require.Equal(
				t,
				plocation.PathLocations{
					plocation.PathLocation{
						UID:      0,
						Path:     "path.s0.A",
						Location: "loc-s0.A",
					},
					plocation.PathLocation{
						UID:      1,
						Path:     "path.s0.B",
						Location: "loc-s0.B",
					},
					plocation.PathLocation{
						UID:      2,
						Path:     "path.s1",
						Location: "loc-s1",
					},
				},
				ploc,
			)

			// valeur is properly allocated
			require.NotNil(t, i)
		},
	)

	t.Run(
		"returning no error", func(t *testing.T) {
			t.Parallel()

			fvs := []fvalues.FieldValues{}

			s0 := NewMockNode(t)
			s0.
				On(
					"Fill",
					mock.IsType(reflect.Value{}),
					fvs,
				).
				Return(
					plocation.PathLocations{
						plocation.PathLocation{
							UID:      0,
							Path:     "path.s0.A",
							Location: "loc-s0.A",
						},
						plocation.PathLocation{
							UID:      1,
							Path:     "path.s0.B",
							Location: "loc-s0.B",
						},
					},
					nil,
				).
				Once()

			s1 := NewMockNode(t)
			s1.
				On(
					"Fill",
					mock.IsType(reflect.Value{}),
					fvs,
				).
				Return(
					plocation.PathLocations{
						plocation.PathLocation{
							UID:      2,
							Path:     "path.s1",
							Location: "loc-s1",
						},
					},
					nil,
				).
				Once()

			type SType struct {
				A interface{}
				B interface{}
			}

			var i *SType

			v := reflect.ValueOf(&i).Elem()

			require.True(t, v.CanSet())

			n := &StructNode{
				Type: reflect.TypeOf(i),
				Index: IndexedSubNodes{
					&IndexedSubNode{
						Node:  s0,
						Index: []int{0},
					},
					&IndexedSubNode{
						Node:  s1,
						Index: []int{1},
					},
				},
			}

			ploc, err := n.Fill(v, fvs)

			// check errors are properly aggregate.
			require.NoError(t, err)

			// path locations are properly aggregate.
			require.Equal(
				t,
				plocation.PathLocations{
					plocation.PathLocation{
						UID:      0,
						Path:     "path.s0.A",
						Location: "loc-s0.A",
					},
					plocation.PathLocation{
						UID:      1,
						Path:     "path.s0.B",
						Location: "loc-s0.B",
					},
					plocation.PathLocation{
						UID:      2,
						Path:     "path.s1",
						Location: "loc-s1",
					},
				},
				ploc,
			)

			// valeur is properly allocated
			require.NotNil(t, i)
		},
	)
}

func TestStructNode_FeedFieldValues(t *testing.T) {
	t.Parallel()

	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()

			fvs := fvalues.FieldValues{}

			type SType struct {
				A *int
				B *float32
			}

			i := &SType{
				A: utils.R(123),
				B: utils.R(float32(123.123)),
			}

			v := reflect.ValueOf(&i).Elem()

			srcID := "srcID"

			s0 := NewMockNode(t)
			s0.
				On(
					"FeedFieldValues",
					srcID,
					fvs,
					mock.MatchedBy(
						func(v reflect.Value) bool {
							i := v.Interface()
							require.IsType(t, i, (*int)(nil))
							vi, ok := i.(*int)
							require.True(t, ok)
							return assert.Equal(t, 123, *vi)
						},
					),
				).
				Return().
				Once()

			s1 := NewMockNode(t)
			s1.
				On(
					"FeedFieldValues",
					srcID,
					fvs,
					mock.MatchedBy(
						func(v reflect.Value) bool {
							i := v.Interface()
							require.IsType(t, i, (*float32)(nil))
							vi, ok := i.(*float32)
							require.True(t, ok)
							return assert.Equal(t, float32(123.123), *vi)
						},
					),
				).
				Return().
				Once()

			n := StructNode{
				Type: nil,
				Index: IndexedSubNodes{
					&IndexedSubNode{
						Node:  s0,
						Index: []int{0},
					},
					&IndexedSubNode{
						Node:  s1,
						Index: []int{1},
					},
				},
			}

			n.FeedFieldValues(srcID, fvs, v)
		},
	)

	t.Run(
		"nil value case", func(t *testing.T) {
			t.Parallel()

			fvs := fvalues.FieldValues{}

			type SType struct {
				A *int
				B *float32
			}

			var i *SType

			v := reflect.ValueOf(&i).Elem()

			srcID := "srcID"

			s0 := NewMockNode(t)
			s1 := NewMockNode(t)

			n := StructNode{
				Type: nil,
				Index: IndexedSubNodes{
					&IndexedSubNode{
						Node:  s0,
						Index: []int{0},
					},
					&IndexedSubNode{
						Node:  s1,
						Index: []int{1},
					},
				},
			}

			n.FeedFieldValues(srcID, fvs, v)
		},
	)
}

func TestIndexedSubNodes_GetIndexes(t *testing.T) {
	sn := IndexedSubNodes{
		&IndexedSubNode{
			Node:  nil,
			Index: []int{1},
		},
		&IndexedSubNode{
			Node:  nil,
			Index: []int{1, 2},
		},
		&IndexedSubNode{
			Node:  nil,
			Index: []int{1, 2, 3},
		},
	}

	require.Equal(
		t, [][]int{
			{1},
			{1, 2},
			{1, 2, 3},
		},
		sn.GetIndexes(),
	)
}

func TestStructNode_BuildGetList(t *testing.T) {
	var gl GetList

	s0 := NewMockNode(t)
	s0.On("BuildGetList", &gl).Return().Once()

	s1 := NewMockNode(t)
	s1.On("BuildGetList", &gl).Return().Once()

	s2 := NewMockNode(t)
	s2.On("BuildGetList", &gl).Return().Once()

	n := &StructNode{
		Type: nil,
		Index: IndexedSubNodes{
			&IndexedSubNode{
				Node: s0,
			},
			&IndexedSubNode{
				Node: s1,
			},
			&IndexedSubNode{
				Node: s2,
			},
		},
	}

	n.BuildGetList(&gl)
}

func TestStructNode_PushSubNodes(t *testing.T) {
	s0 := NewMockNode(t)

	n := &StructNode{
		Index: nil,
	}

	fieldIndex := []int{1, 2, 3}
	n.PushSubNodes(fieldIndex, s0)

	require.Len(t, n.Index, 1)
	require.Equal(
		t,
		IndexedSubNode{
			Node:  s0,
			Index: fieldIndex,
		},
		*(n.Index[0]),
	)
}

func TestStructNodeError_Is(t *testing.T) {
	{
		e := MockedError1{}

		require.False(t, errors.Is(e, ErrStructNode))
	}
	{
		e := StructNodeError{}
		require.True(t, errors.Is(e, ErrStructNode))
	}
}
