package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/internal/fvalue"
	"github.com/byte4ever/dsco/internal/plocation"
)

func TestValueNode_Fill(t *testing.T) {
	t.Parallel()

	t.Run(
		"", func(t *testing.T) {
			t.Parallel()

			n := &ValueNode{
				Type:        nil,
				VisiblePath: "the.path",
				UID:         50,
			}

			o := 128
			ov := reflect.ValueOf(&o)

			var i *int

			v := reflect.ValueOf(&i).Elem()

			require.True(t, v.CanSet())

			fvs := []fvalue.Values{
				{},
				{
					uint(50): {
						Value:    ov,
						Location: "some-location",
					},
				},
				{},
			}
			ploc, err := n.Fill(
				v,
				fvs,
			)

			require.NoError(t, err)
			require.Contains(
				t, ploc, plocation.Location{
					UID:      50,
					Path:     "the.path",
					Location: "some-location",
				},
			)
			require.NotContains(t, fvs[1], uint(50))
			require.Equal(t, 128, *i)
		},
	)

	t.Run(
		"", func(t *testing.T) {
			t.Parallel()

			n := &ValueNode{
				Type:        nil,
				VisiblePath: "the.path",
				UID:         50,
			}

			o := 128
			ov := reflect.ValueOf(&o)

			var i *int

			v := reflect.ValueOf(&i).Elem()

			require.True(t, v.CanSet())

			ploc, err := n.Fill(
				v,
				[]fvalue.Values{
					{},
					{
						uint(500): {
							Value:    ov,
							Location: "some-location",
						},
					},
					{},
				},
			)

			require.Empty(
				t, ploc,
			)
			require.ErrorIs(t, err, ErrUninitializedKey)
			require.ErrorContains(t, err, "the.path")
			require.Nil(t, i)
		},
	)
}

func TestValueNode_FeedFieldValues(t *testing.T) {
	t.Parallel()

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			n := &ValueNode{
				VisiblePath: "the.path",
				UID:         50,
			}

			fvs := fvalue.Values{}

			i := 15
			pi := &i

			vpi := reflect.ValueOf(pi)

			n.FeedFieldValues(
				"srcID",
				fvs,
				vpi,
			)

			require.Contains(t, fvs, uint(50))

			require.Equal(
				t, fvalue.Value{
					Value:    vpi,
					Location: "struct[srcID]:the.path",
				},
				*fvs[uint(50)],
			)
		},
	)

	t.Run(
		"nil pointer case",
		func(t *testing.T) {
			t.Parallel()

			n := &ValueNode{
				VisiblePath: "the.path",
				UID:         50,
			}

			fvs := fvalue.Values{}

			pi := (*int)(nil)

			vpi := reflect.ValueOf(pi)

			n.FeedFieldValues(
				"srcID",
				fvs,
				vpi,
			)

			require.Empty(t, fvs)
		},
	)
}

func TestValueNode_BuildGetList(t *testing.T) {
	t.Parallel()

	t.Run(
		"",
		func(t *testing.T) {
			t.Parallel()

			someType := reflect.TypeOf(123)
			someValue := reflect.ValueOf(345)
			const path = "the.path"
			const expectedUID = uint(50)
			const location = "some-loc"

			n := &ValueNode{
				Type:        someType,
				VisiblePath: path,
				UID:         expectedUID,
			}

			var gl GetList

			n.BuildGetList(&gl)
			require.Len(t, gl, 1)

			g := NewMockGetter(t)
			g.
				On(
					"Get",
					path,
					someType,
				).
				Return(
					&fvalue.Value{
						Value:    someValue,
						Location: location,
					},
					nil,
				).
				Once()

			uid, fv, err := gl[0](g)
			require.NoError(t, err)
			require.Equal(t, expectedUID, uid)
			require.NotNil(t, fv)
			require.Equal(
				t,
				fvalue.Value{
					Value:    someValue,
					Location: location,
				},
				*fv,
			)
		},
	)
}
