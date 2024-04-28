package dsco

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/internal/fvalue"
	"github.com/byte4ever/dsco/internal/ierror"
	"github.com/byte4ever/dsco/internal/model"
	"github.com/byte4ever/dsco/svalue"
)

func TestStringBasedBuilder_Get(t *testing.T) {
	t.Parallel()

	t.Run(
		"none found", func(t *testing.T) {
			t.Parallel()

			sb := StringBasedBuilder{
				values: map[string]*svalue.Value{},
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

				values: map[string]*svalue.Value{},
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
				values: map[string]*svalue.Value{
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
				values: map[string]*svalue.Value{
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
		"success complex slice", func(t *testing.T) {
			t.Parallel()

			// S3MappingItem defines mapping from a S3 bucket.
			type S3MappingItem struct {
				// The ID of the S3 mapping.
				ID *string `yaml:"id,omitempty"`

				// Bucket name on S3.
				BucketName *string `yaml:"bucket_name,omitempty"`

				// prefix to apply to the object key.
				Prefix *string `yaml:"prefix,omitempty"`
			}

			// S3MappingItems defines a list of single S3 mapping items.
			type S3MappingItems []*S3MappingItem

			type Conf struct {
				S3 S3MappingItems `yaml:"s3,omitempty"`
			}

			sb := StringBasedBuilder{
				values: map[string]*svalue.Value{
					"some-path": {
						Location: "loc1",
						Value: `---
s3:
  - id: bolos
    bucket_name: my_bucket
    prefix: "p"
`,
					},
				},
			}

			var v *Conf
			pv := v

			gotFv, err := sb.Get("Some.Path", reflect.TypeOf(pv))
			require.NoError(t, err)
			require.NotNil(t, gotFv)

			require.Equal(t, "loc1", gotFv.Location)
			require.Equal(
				t,
				&Conf{
					S3MappingItems{
						{
							ID:         R("bolos"),
							BucketName: R("my_bucket"),
							Prefix:     R("p"),
						},
					},
				},
				gotFv.Value.Interface(),
			)

			require.NoError(t, err)
		},
	)

	t.Run(
		"success complex slice", func(t *testing.T) {
			t.Parallel()

			sb := StringBasedBuilder{
				values: map[string]*svalue.Value{
					"some-path": {
						Location: "loc1",
						Value:    "[{a: 123, b:123.32}]",
					},
				},
			}

			type ComplexItem struct {
				A *int     `yaml:"a,omitempty"`
				B *float64 `yaml:"b,omitempty"`
			}

			var v []*ComplexItem
			pv := v

			gotFv, err := sb.Get("Some.Path", reflect.TypeOf(pv))
			require.NoError(t, err)
			require.NotNil(t, gotFv)

			require.Equal(t, "loc1", gotFv.Location)
			require.IsType(
				t,
				[]*ComplexItem{
					{
						A: R(102),
						B: R(3.14),
					},
				},
				gotFv.Value.Interface(),
			)
		},
	)

	t.Run(
		"parse error pointer", func(t *testing.T) {
			t.Parallel()

			sb := StringBasedBuilder{
				values: map[string]*svalue.Value{
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
				values: map[string]*svalue.Value{
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
				values: map[string]*svalue.Value{
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

	t.Run(
		"success expand", func(t *testing.T) {
			t.Parallel()

			sb := StringBasedBuilder{
				expandedValues: map[string]*fvalue.Value{
					"Some.Path": {
						Value:    reflect.ValueOf(R(123)),
						Location: "loc1",
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
}

func TestGetError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrGet)
	require.ErrorIs(t, GetError{}, ErrGet)
}

func TestStringBasedBuilder_GetFieldValuesFrom(t *testing.T) {
	t.Parallel()
	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()

			fvs := fvalue.Values{
				uint(101): &fvalue.Value{
					Location: "loc1",
				},
			}

			sb := &StringBasedBuilder{}

			model := NewMockModelInterface(t)
			model.
				On("ApplyOn", sb).
				Return(
					fvs, nil,
				).
				Once()

			model.
				On("Expand", sb).
				Return(
					nil,
				).
				Once()

			gotFvs, gotErr := sb.GetFieldValuesFrom(model)
			require.NoError(t, gotErr)
			require.Equal(t, fvs, gotFvs)
		},
	)

	t.Run(
		"failures from apply on", func(t *testing.T) {
			t.Parallel()

			fvs := fvalue.Values{
				uint(101): &fvalue.Value{
					Location: "loc1",
				},
			}

			sb := &StringBasedBuilder{
				values: map[string]*svalue.Value{
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

			model := NewMockModelInterface(t)
			model.
				EXPECT().
				Expand(sb).
				Return(
					nil,
				).
				Once()

			model.
				EXPECT().
				ApplyOn(sb).
				Return(
					fvs,
					errMocked1,
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

	t.Run(
		"failures from expand", func(t *testing.T) {
			t.Parallel()

			fvs := fvalue.Values{
				uint(101): &fvalue.Value{
					Location: "loc1",
				},
			}

			sb := &StringBasedBuilder{
				values: map[string]*svalue.Value{
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
				expandedValues: map[string]*fvalue.Value{
					"exp-a": {
						Location: "exp-loc-a",
					},
					"exp-b": {
						Location: "exp-loc-b",
					},
					"exp-c": {
						Location: "exp-loc-c",
					},
				},
			}

			model := NewMockModelInterface(t)
			model.
				EXPECT().
				Expand(sb).
				Return(
					errMocked1,
				).
				Once()

			model.
				EXPECT().
				ApplyOn(sb).
				Return(
					fvs,
					nil,
				).
				Once()

			gotFvs, gotErr := sb.GetFieldValuesFrom(model)
			require.Nil(t, gotFvs)

			var e GetError

			require.ErrorAs(t, gotErr, &e)

			require.Equal(t, 7, e.Count())

			require.ErrorIs(t, e.MError[0], errMocked1)

			for idx, expectedLoc := range []string{
				"exp-loc-a",
				"exp-loc-b",
				"exp-loc-c",
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

			for idx, expectedLoc := range []string{
				"loc-a",
				"loc-b",
				"loc-c",
			} {
				var ue UnboundedLocationError

				require.ErrorAs(t, e.MError[idx+1+3], &ue)
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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

	require.Equal(
		t,
		"unbounded location loc-a",
		UnboundedLocationError{
			Location: "loc-a",
		}.Error(),
	)
}

func TestUnboundedLocationError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, UnboundedLocationError{}, errMocked1)
	require.ErrorIs(t, UnboundedLocationError{}, ErrUnboundedLocation)
}

func TestParseError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, ParseError{}, errMocked1)
	require.ErrorIs(t, ParseError{}, ErrParse)
}

func TestParseError_Error(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	require.NotErrorIs(t, AliasCollisionError{}, errMocked1)
	require.ErrorIs(t, AliasCollisionError{}, ErrAliasCollision)
}

func TestAliasCollisionError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"alias collision-path collides with structure",
		AliasCollisionError{
			Path: "collision-path",
		}.Error(),
	)
}

func TestWithAliases(t *testing.T) {
	t.Parallel()

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

func TestNewStringBasedBuilder(t *testing.T) {
	t.Parallel()

	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()
			p := NewMockStringValuesProvider(t)

			sv := svalue.Values{
				"a": &svalue.Value{
					Location: "",
					Value:    "",
				},
			}

			p.On("GetStringValues").Return(sv).Once()

			o := NewMockOption(t)
			o.
				On(
					"apply",
					mock.MatchedBy(
						func(o *internalOpts) bool {
							return assert.NotNil(t, o)
						},
					),
				).
				Return(nil).
				Once()

			b, err := NewStringBasedBuilder(p, o)

			require.NoError(t, err)
			require.Equal(t, sv, b.values)
		},
	)

	t.Run(
		"nil provider", func(t *testing.T) {
			t.Parallel()

			b, err := NewStringBasedBuilder(nil)

			require.ErrorIs(t, err, ErrNilProvider)
			require.Nil(t, b)
		},
	)

	t.Run(
		"option error", func(t *testing.T) {
			t.Parallel()
			p := NewMockStringValuesProvider(t)

			o1 := NewMockOption(t)
			o1.
				On(
					"apply",
					mock.MatchedBy(
						func(o *internalOpts) bool {
							return assert.NotNil(t, o)
						},
					),
				).
				Return(nil).
				Once()

			o2 := NewMockOption(t)
			o2.
				On(
					"apply",
					mock.MatchedBy(
						func(o *internalOpts) bool {
							return assert.NotNil(t, o)
						},
					),
				).
				Return(errMocked1).
				Once()

			b, err := NewStringBasedBuilder(p, o1, o2)
			require.Nil(t, b)

			var e ierror.IError
			require.ErrorAs(t, err, &e)

			require.Equal(
				t, ierror.IError{
					Index: 1,
					Info:  "when processing option",
					Err:   errMocked1,
				}, e,
			)
		},
	)

	t.Run(
		"with alias option", func(t *testing.T) {
			t.Parallel()
			p := NewMockStringValuesProvider(t)
			p.EXPECT().GetStringValues().Return(
				svalue.Values{
					"aliaskey": {
						Location: "l1",
						Value:    "v1",
					},
					"key": {
						Location: "l2",
						Value:    "v2",
					},
				},
			)

			b, err := NewStringBasedBuilder(p, WithAliases(map[string]string{
				"aliaskey": "actual-key",
			}))

			require.NoError(t, err)

			require.Equal(t, b.values, svalue.Values{
				"actual-key": {
					Location: "l1",
					Value:    "v1",
				},
				"key": {
					Location: "l2",
					Value:    "v2",
				},
			})
		},
	)
}

func TestOverriddenKeyError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"for path <path> <location> is override by <overrideLocation>",
		OverriddenKeyError{
			Path:             "<path>",
			Location:         "<location>",
			OverrideLocation: "<overrideLocation>",
		}.Error(),
	)
}

func TestOverriddenKeyError_Is(t *testing.T) {
	t.Parallel()

	require.ErrorIs(t, OverriddenKeyError{}, ErrOverriddenKey)
	require.NotErrorIs(t, OverriddenKeyError{}, errMocked1)
}

func TestStringBasedBuilder_ExpandStruct(t *testing.T) {
	t.Parallel()

	type SomeStruct struct {
		A *int     `yaml:"a"`
		B *float64 `yaml:"b"`
		C *struct {
			X *int     `yaml:"x"`
			Y *float64 `yaml:"y"`
		} `yaml:"c"`
	}

	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()

			var p *SomeStruct

			sb := StringBasedBuilder{
				values: map[string]*svalue.Value{
					"p-t": {
						Location: "",
						Value: `---
a: 123
b: 123.412
c:
  x: 2332
  y: 123.234`,
					},
				},
				expandedValues: make(map[string]*fvalue.Value),
			}

			err := sb.ExpandStruct("P.T", reflect.TypeOf(p))
			require.NoError(t, err)

			require.Contains(t, sb.expandedValues, "P.T.A")
			require.Equal(t, R(123), sb.expandedValues["P.T.A"].Value.Interface())
			require.Contains(t, sb.expandedValues, "P.T.B")
			require.Equal(t, R(123.412), sb.expandedValues["P.T.B"].Value.Interface())
			require.Contains(t, sb.expandedValues, "P.T.C.X")
			require.Equal(t, R(2332), sb.expandedValues["P.T.C.X"].Value.Interface())
			require.Contains(t, sb.expandedValues, "P.T.C.Y")
			require.Equal(t, R(123.234), sb.expandedValues["P.T.C.Y"].Value.Interface())
		},
	)

	t.Run(
		"partial success", func(t *testing.T) {
			t.Parallel()

			var p *SomeStruct

			sb := StringBasedBuilder{
				values: map[string]*svalue.Value{
					"p-t": {
						Location: "",
						Value: `---
a: 123
c:
  x: 2332
  y: 123.234`,
					},
				},
				expandedValues: make(map[string]*fvalue.Value),
			}

			err := sb.ExpandStruct("P.T", reflect.TypeOf(p))
			require.NoError(t, err)

			require.Contains(t, sb.expandedValues, "P.T.A")
			require.Equal(t, R(123), sb.expandedValues["P.T.A"].Value.Interface())
			require.Contains(t, sb.expandedValues, "P.T.C.X")
			require.Equal(t, R(2332), sb.expandedValues["P.T.C.X"].Value.Interface())
			require.Contains(t, sb.expandedValues, "P.T.C.Y")
			require.Equal(t, R(123.234), sb.expandedValues["P.T.C.Y"].Value.Interface())
		},
	)

	t.Run(
		"no string value", func(t *testing.T) {
			t.Parallel()

			var p *SomeStruct

			sb := StringBasedBuilder{
				values:         map[string]*svalue.Value{},
				expandedValues: make(map[string]*fvalue.Value),
			}

			err := sb.ExpandStruct("P.T", reflect.TypeOf(p))
			require.NoError(t, err)
			require.Empty(t, sb.values)
			require.Empty(t, sb.expandedValues)
		},
	)

	t.Run(
		"parse error", func(t *testing.T) {
			t.Parallel()

			var p *SomeStruct

			sb := StringBasedBuilder{
				values: map[string]*svalue.Value{
					"p-t": {
						Location: "",
						Value: `---
a: 123
b: 123.412
c:
  x: bad-int
  y: 123.234`,
					},
				},
				expandedValues: make(map[string]*fvalue.Value),
			}

			err := sb.ExpandStruct("P.T", reflect.TypeOf(p))
			require.ErrorIs(t, err, &ParseError{})
			var parseErr *ParseError
			errors.As(err, &parseErr)
			require.Equal(t, "P.T", parseErr.Path)
			require.Equal(t, reflect.TypeOf(p), parseErr.Type)
			require.Empty(t, parseErr.Location)
		},
	)

	t.Run(
		"alias collision", func(t *testing.T) {
			t.Parallel()

			var p *SomeStruct

			sb := StringBasedBuilder{
				internalOpts: internalOpts{
					aliases: map[string]string{
						"alias": "aliased-key",
					},
				},
			}

			err := sb.ExpandStruct("Alias", reflect.TypeOf(p))

			var asErr *AliasCollisionError
			require.ErrorAs(t, err, &asErr)
			require.Equal(t, "Alias", asErr.Path)
		},
	)

	t.Run(
		"model failure", func(t *testing.T) {
			t.Parallel()

			type SomeBadStruct struct {
				A int      `yaml:"a"`
				B *float64 `yaml:"b"`
				C *struct {
					X *int     `yaml:"x"`
					Y *float64 `yaml:"y"`
				} `yaml:"c"`
			}

			var p *SomeBadStruct

			sb := StringBasedBuilder{
				values: map[string]*svalue.Value{
					"p-t": {
						Location: "",
						Value: `---
a: 123
b: 123.412
c:
  x: 2332
  y: 123.234`,
					},
				},
				expandedValues: make(map[string]*fvalue.Value),
			}

			err := sb.ExpandStruct("P.T", reflect.TypeOf(p))

			var asErr model.UnsupportedTypeError
			require.ErrorAs(t, err, &asErr)
		},
	)
}
