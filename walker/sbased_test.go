package walker

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/walker/svalues"
)

func TestStringBasedBuilder_Get(t *testing.T) {
	t.Parallel()

	t.Run(
		"alias collision", func(t *testing.T) {
			t.Parallel()

			sb := StringBasedBuilder{
				internalOpts: internalOpts{
					aliases: map[string]string{"some-path": "alias"},
				},

				values: map[string]*svalues.StringValue{},
			}

			gotFv, err := sb.Get("Some.Path", nil)
			require.Nil(t, gotFv)
			require.ErrorIs(t, err, ErrAliasCollision)

			var e AliasCollisionError
			require.ErrorAs(t, err, &e)
			require.Equal(
				t, AliasCollisionError{
					Path: "Some.Path",
				}, e,
			)
		},
	)

	t.Run(
		"success pointer", func(t *testing.T) {
			t.Parallel()

			sb := StringBasedBuilder{
				values: map[string]*svalues.StringValue{
					"some-path": {
						Location: "loc1",
						Value:    "123",
					},
				},
			}

			var v int
			pv := &v

			gotFv, err := sb.Get("Some.Path", reflect.TypeOf(pv))
			require.NoError(t, err)
			require.NotNil(t, gotFv)

			require.Equal(t, "loc1", gotFv.Location)
			require.IsType(t, 123, gotFv.Value.Elem().Interface())
		},
	)

	t.Run(
		"success slice", func(t *testing.T) {
			t.Parallel()

			sb := StringBasedBuilder{
				values: map[string]*svalues.StringValue{
					"some-path": {
						Location: "loc1",
						Value:    "[1,2,3,4,5]",
					},
				},
			}

			var v []int
			pv := v

			gotFv, err := sb.Get("Some.Path", reflect.TypeOf(pv))
			require.NoError(t, err)
			require.NotNil(t, gotFv)

			require.Equal(t, "loc1", gotFv.Location)
			require.IsType(
				t,
				[]int{1, 2, 3, 4, 5},
				gotFv.Value.Interface(),
			)
		},
	)

	t.Run(
		"parse error pointer", func(t *testing.T) {
			t.Parallel()

			sb := StringBasedBuilder{
				values: map[string]*svalues.StringValue{
					"some-path": {
						Location: "loc1",
						Value:    "asd",
					},
				},
			}

			var v int
			pv := &v
			vType := reflect.TypeOf(pv)

			gotFv, err := sb.Get("Some.Path", vType)
			require.Nil(t, gotFv)
			require.ErrorIs(t, err, ErrParse)

			var e ParseError
			require.ErrorAs(t, err, &e)

			require.Equal(
				t, ParseError{
					Path:     "Some.Path",
					Type:     vType,
					Location: "loc1",
				}, e,
			)
		},
	)

	t.Run(
		"parse error slice", func(t *testing.T) {
			t.Parallel()

			sb := StringBasedBuilder{
				values: map[string]*svalues.StringValue{
					"some-path": {
						Location: "loc1",
						Value:    "asd",
					},
				},
			}

			var v []int
			pv := v
			vType := reflect.TypeOf(pv)

			gotFv, err := sb.Get("Some.Path", vType)
			require.Nil(t, gotFv)
			require.ErrorIs(t, err, ErrParse)

			var e ParseError
			require.ErrorAs(t, err, &e)

			require.Equal(
				t, ParseError{
					Path:     "Some.Path",
					Type:     vType,
					Location: "loc1",
				}, e,
			)
		},
	)

	t.Run(
		"invalid type", func(t *testing.T) {
			t.Parallel()

			sb := StringBasedBuilder{
				values: map[string]*svalues.StringValue{
					"some-path": {
						Location: "loc1",
						Value:    "asd",
					},
				},
			}

			var v int
			pv := v
			vType := reflect.TypeOf(pv)

			gotFv, err := sb.Get("Some.Path", vType)

			require.Nil(t, gotFv)
			require.ErrorIs(t, err, ErrInvalidType)
		},
	)

}
