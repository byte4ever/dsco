package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/walker/fvalues"
	"github.com/byte4ever/dsco/walker/plocation"
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

			fvs := []fvalues.FieldValues{
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
				t, ploc, plocation.PathLocation{
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
				[]fvalues.FieldValues{
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
