package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/walker/fvalues"
	"github.com/byte4ever/dsco/walker/plocation"
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
