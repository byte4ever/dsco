package walker

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/walker/fvalues"
	"github.com/byte4ever/dsco/walker/ifaces"
	"github.com/byte4ever/dsco/walker/svalues"
)

func TestStringBasedBuilder_Get(t *testing.T) {
	t.Parallel()

	t.Run(
		"none found", func(t *testing.T) {
			t.Parallel()

			sb := StringBasedBuilder{
				values: map[string]*svalues.StringValue{},
			}

			gotFv, err := sb.Get("Some.Path", nil)
			require.Nil(t, gotFv)
			require.NoError(t, err)
		},
	)

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

func TestGetError_Is(t *testing.T) {
	require.NotErrorIs(t, errMocked1, ErrGet)
	require.ErrorIs(t, GetError{}, ErrGet)
}

func TestStringBasedBuilder_GetFieldValuesFrom(t *testing.T) {
	t.Parallel()
	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()

			fvs := fvalues.FieldValues{
				uint(101): &fvalues.FieldValue{
					Location: "loc1",
				},
			}

			sb := &StringBasedBuilder{}

			model := ifaces.NewMockModelInterface(t)
			model.
				On("ApplyOn", sb).
				Return(
					fvs, nil,
				).
				Once()

			gotFvs, gotErr := sb.GetFieldValuesFrom(model)
			require.NoError(t, gotErr)
			require.Equal(t, fvs, gotFvs)
		},
	)

	t.Run(
		"failures", func(t *testing.T) {
			t.Parallel()

			fvs := fvalues.FieldValues{
				uint(101): &fvalues.FieldValue{
					Location: "loc1",
				},
			}

			sb := &StringBasedBuilder{
				values: map[string]*svalues.StringValue{
					"a": {
						Location: "loc-a",
					},
					"b": {
						Location: "loc-b",
					},
					"c": {
						Location: "loc-c",
					},
				},
			}

			model := ifaces.NewMockModelInterface(t)
			model.
				On("ApplyOn", sb).
				Return(
					fvs, errMocked1,
				).
				Once()

			gotFvs, gotErr := sb.GetFieldValuesFrom(model)
			require.Nil(t, gotFvs)

			var e GetError

			require.ErrorAs(t, gotErr, &e)

			require.Equal(t, 4, e.Count())

			require.ErrorIs(t, e.MError[0], errMocked1)

			for idx, expectedLoc := range []string{
				"loc-a",
				"loc-b",
				"loc-c",
			} {
				var ue UnboundedLocationError

				require.ErrorAs(t, e.MError[idx+1], &ue)
				require.Equal(
					t,
					UnboundedLocationError{
						Location: expectedLoc,
					},
					ue,
				)
			}
		},
	)
}

func TestUnboundedLocationErrors_Len(t *testing.T) {
	require.Equal(t, 0, UnboundedLocationErrors{}.Len())
	require.Equal(t, 0, UnboundedLocationErrors(nil).Len())
	require.Equal(
		t, 2, UnboundedLocationErrors{
			UnboundedLocationError{},
			UnboundedLocationError{},
		}.Len(),
	)
}

func TestUnboundedLocationErrors_Swap(t *testing.T) {
	l := UnboundedLocationErrors{
		UnboundedLocationError{
			Location: "A",
		},
		UnboundedLocationError{
			Location: "B",
		},
	}

	expected := UnboundedLocationErrors{
		UnboundedLocationError{
			Location: "B",
		},
		UnboundedLocationError{
			Location: "A",
		},
	}

	l.Swap(0, 1)
	require.Equal(t, expected, l)
}

func TestUnboundedLocationError_Error(t *testing.T) {
	require.Equal(
		t,
		"unbounded location loc-a",
		UnboundedLocationError{
			Location: "loc-a",
		}.Error(),
	)
}

func TestUnboundedLocationError_Is(t *testing.T) {
	require.NotErrorIs(t, UnboundedLocationError{}, errMocked1)
	require.ErrorIs(t, UnboundedLocationError{}, ErrUnboundedLocation)
}

func TestParseError_Is(t *testing.T) {
	require.NotErrorIs(t, ParseError{}, errMocked1)
	require.ErrorIs(t, ParseError{}, ErrParse)
}

func TestParseError_Error(t *testing.T) {
	require.Equal(
		t,
		"parse error on some-path-<int> loc-a",
		ParseError{
			Path:     "some-path",
			Type:     reflect.TypeOf(10),
			Location: "loc-a",
		}.Error(),
	)
}

func TestAliasCollisionError_Is(t *testing.T) {
	require.NotErrorIs(t, AliasCollisionError{}, errMocked1)
	require.ErrorIs(t, AliasCollisionError{}, ErrAliasCollision)
}

func TestAliasCollisionError_Error(t *testing.T) {
	require.Equal(
		t,
		"alias collision-path collides with structure",
		AliasCollisionError{
			Path: "collision-path",
		}.Error(),
	)
}

func TestWithAliases(t *testing.T) {

	mapping := map[string]string{
		"a": "ta",
		"b": "tb",
		"c": "tc",
	}
	r := WithAliases(mapping)
	require.Equal(t, AliasesOption(mapping), r)
}

func TestAliasesOption_apply(t *testing.T) {
	t.Parallel()

	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()

			mapping := map[string]string{
				"a": "ta",
				"b": "tb",
				"c": "tc",
			}

			var io internalOpts

			ao := AliasesOption(mapping)
			require.NoError(t, ao.apply(&io))
			require.Equal(t, mapping, io.aliases)
		},
	)

	t.Run(
		"failure", func(t *testing.T) {
			t.Parallel()

			mapping := map[string]string{}

			var io internalOpts

			ao := AliasesOption(mapping)
			require.ErrorIs(t, ao.apply(&io), ErrNoAliasesProvided)
		},
	)
}
