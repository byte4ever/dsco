package dsco

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/internal/ierror"
	"github.com/byte4ever/dsco/ref"
	"github.com/byte4ever/dsco/svalue"
)

func TestLayers_GetPolicies(t *testing.T) {
	t.Parallel()

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()
			l1 := NewMockLayer(t)
			l1.
				On(
					"register",
					mock.MatchedBy(
						func(to *layerBuilder) bool {
							return assert.NotNil(t, to)
						},
					),
				).Return(nil).
				Once()

			l2 := NewMockLayer(t)
			l2.
				On(
					"register",
					mock.MatchedBy(
						func(to *layerBuilder) bool {
							return assert.NotNil(t, to)
						},
					),
				).Return(nil).
				Once()

			l3 := NewMockLayer(t)
			l3.
				On(
					"register",
					mock.MatchedBy(
						func(to *layerBuilder) bool {
							return assert.NotNil(t, to)
						},
					),
				).Return(nil).
				Once()

			l4 := NewMockLayer(t)
			l4.
				On(
					"register",
					mock.MatchedBy(
						func(to *layerBuilder) bool {
							return assert.NotNil(t, to)
						},
					),
				).Return(nil).
				Once()

			layers := Layers{l1, l2, l3, l4}

			clp, err := layers.GetPolicies()

			require.NoError(t, err)
			require.NotNil(t, clp)
			require.Len(t, clp, 0)
			require.Equal(t, 4, cap(clp))
		},
	)

	t.Run(
		"fail",
		func(t *testing.T) {
			t.Parallel()

			l1 := NewMockLayer(t)
			l1.
				On(
					"register",
					mock.MatchedBy(
						func(to *layerBuilder) bool {
							return assert.NotNil(t, to)
						},
					),
				).Return(errMocked1).
				Once()

			l2 := NewMockLayer(t)
			l2.
				On(
					"register",
					mock.MatchedBy(
						func(to *layerBuilder) bool {
							return assert.NotNil(t, to)
						},
					),
				).Return(nil).
				Once()

			l3 := NewMockLayer(t)
			l3.
				On(
					"register",
					mock.MatchedBy(
						func(to *layerBuilder) bool {
							return assert.NotNil(t, to)
						},
					),
				).Return(errMocked2).
				Once()

			l4 := NewMockLayer(t)
			l4.
				On(
					"register",
					mock.MatchedBy(
						func(to *layerBuilder) bool {
							return assert.NotNil(t, to)
						},
					),
				).Return(nil).
				Once()

			layers := Layers{l1, l2, l3, l4}

			clp, err := layers.GetPolicies()

			require.Nil(t, clp)

			var expectedErr LayerErrors

			require.ErrorAs(t, err, &expectedErr)

			for idx, et := range []struct {
				e  error
				ie int
			}{
				{errMocked1, 0},
				{errMocked2, 2},
			} {
				var ie ierror.IError
				require.ErrorAs(t, expectedErr.MError[idx], &ie)
				require.Equal(t, et.ie, ie.Index)
				require.Equal(t, et.e, ie.Err)
			}
		},
	)
}

func TestWithStrictStructLayer(t *testing.T) {
	t.Parallel()

	v := 10
	id := "id"

	k := WithStrictStructLayer(v, id)

	require.Equal(t, v, k.input)
	require.Equal(t, id, k.id)
}

func TestWithStructLayer(t *testing.T) {
	t.Parallel()

	v := 10
	id := "id"

	k := WithStructLayer(v, id)

	require.Equal(t, v, k.input)
	require.Equal(t, id, k.id)
}

func TestWithEnvLayer(t *testing.T) {
	t.Parallel()

	prefix := "id"

	k := WithEnvLayer(prefix)

	require.Equal(t, prefix, k.prefix)
}

func TestWithStrictEnvLayer(t *testing.T) {
	t.Parallel()

	prefix := "id"

	k := WithStrictEnvLayer(prefix)

	require.Equal(t, prefix, k.prefix)
}

func TestWithCmdlineLayer(t *testing.T) {
	t.Parallel()

	o1 := NewMockOption(t)
	o2 := NewMockOption(t)

	k := WithCmdlineLayer(o1, o2)

	require.Equal(t, []Option{o1, o2}, k.options)
}

func TestWithStrictCmdlineLayer(t *testing.T) {
	t.Parallel()

	o1 := NewMockOption(t)
	o2 := NewMockOption(t)

	k := WithStrictCmdlineLayer(o1, o2)

	require.Equal(t, []Option{o1, o2}, k.options)
}

//nolint:paralleltest // using global variable
func TestStrictCmdlineLayer_register(t *testing.T) {
	for _, x := range []struct {
		layer  Layer
		strict bool
	}{
		{
			layer:  &StrictCmdlineLayer{},
			strict: true,
		},
		{
			layer:  &CmdlineLayer{},
			strict: false,
		},
	} {
		x := x

		t.Run(
			"success", func(t *testing.T) {
				os.Args = []string{"cmdName"}

				lb := newLayerBuilder(1)

				err := x.layer.register(lb)

				require.NoError(t, err)
				require.Len(t, lb.builders, 1)
				require.Equal(t, x.strict, lb.builders[0].isStrict())
				require.Contains(t, lb.idDedup, "cmdLine")
			},
		)

		t.Run(
			"using twice", func(t *testing.T) {
				os.Args = []string{"cmdName"}

				lb := newLayerBuilder(1)
				lb.idDedup["cmdLine"] = 101

				err := x.layer.register(lb)

				var e CmdlineAlreadyUsedError
				require.ErrorAs(t, err, &e)
				require.Len(t, lb.builders, 0)
				require.Equal(t, 101, e.Index)
			},
		)

		t.Run(
			"cmdline error", func(t *testing.T) {
				os.Args = []string{"cmdName", "asdasdasd"}

				lb := newLayerBuilder(1)

				require.Error(t, x.layer.register(lb))
			},
		)
	}
}

//nolint:paralleltest // dealing with env variables
func TestEnvLayer_register(t *testing.T) {
	for _, x := range []struct {
		layer  Layer
		name   string
		strict bool
	}{
		{
			name: "strict",
			layer: &StrictEnvLayer{
				prefix: "API",
			},
			strict: true,
		},
		{
			name: "normal",
			layer: &EnvLayer{
				prefix: "API",
			},
			strict: false,
		},
	} {
		x := x

		t.Run(
			fmt.Sprintf("%s success", x.name),
			func(t *testing.T) {
				lb := newLayerBuilder(1)
				err := x.layer.register(lb)

				require.NoError(t, err)
				require.Len(t, lb.builders, 1)
				require.Equal(t, x.strict, lb.builders[0].isStrict())
				require.Contains(t, lb.idDedup, "env(API)")
			},
		)

		t.Run(
			fmt.Sprintf("%s using same prefix", x.name),
			func(t *testing.T) {
				lb := newLayerBuilder(1)
				lb.idDedup["env(API)"] = 101

				err := x.layer.register(lb)

				var e DuplicateEnvPrefixError
				require.ErrorAs(t, err, &e)
				require.Len(t, lb.builders, 0)
				require.Equal(t, 101, e.Index)
				require.Equal(t, "API", e.Prefix)
			},
		)

		t.Run(
			fmt.Sprintf("%s env error", x.name),
			func(t *testing.T) {
				t.Setenv("API-123123-_d--__/", "value")
				lb := newLayerBuilder(1)
				require.Error(t, x.layer.register(lb))
			},
		)
	}
}

//nolint:paralleltest // dealing with global variables
func TestStructLayer_register(t *testing.T) {
	type Root struct {
		X *float32
		Y *float64
	}

	k := &Root{
		X: ref.R(float32(123.123)),
	}

	for _, x := range []struct {
		layer  Layer
		name   string
		strict bool
	}{
		{
			name: "strict",
			layer: &StrictStructLayer{
				input: k,
				id:    "default",
			},
			strict: true,
		},
		{
			name: "normal",
			layer: &StructLayer{
				input: k,
				id:    "default",
			},
			strict: false,
		},
	} {
		x := x

		t.Run(
			fmt.Sprintf("%s success", x.name),
			func(t *testing.T) {
				lb := newLayerBuilder(1)
				err := x.layer.register(lb)

				require.NoError(t, err)
				require.Len(t, lb.builders, 1)
				require.Equal(t, x.strict, lb.builders[0].isStrict())
				require.Contains(
					t,
					lb.idDedup,
					fmt.Sprintf(
						"structPtr(%d)",
						reflect.ValueOf(k).Pointer(),
					),
				)
				require.Contains(
					t,
					lb.idDedup,
					"structId(default)",
				)
			},
		)

		t.Run(
			fmt.Sprintf("%s using same id", x.name),
			func(t *testing.T) {
				lb := newLayerBuilder(1)
				lb.idDedup["structId(default)"] = 101

				err := x.layer.register(lb)

				var e DuplicateStructIDError
				require.ErrorAs(t, err, &e)
				require.Len(t, lb.builders, 0)
				require.Equal(t, 101, e.Index)
				require.Equal(t, "default", e.ID)
			},
		)

		t.Run(
			fmt.Sprintf("%s using same ptr", x.name),
			func(t *testing.T) {
				os.Args = []string{"cmdName"}

				lb := newLayerBuilder(1)
				lb.idDedup["structId(default)"] = 101
				lb.idDedup[fmt.Sprintf(
					"structPtr(%d)",
					reflect.ValueOf(k).Pointer(),
				)] = 101

				err := x.layer.register(lb)

				var e DuplicateInputStructError
				require.ErrorAs(t, err, &e)
				require.Len(t, lb.builders, 0)
				require.Equal(t, 101, e.Index)
			},
		)
	}
}

func TestStructLayer_register2(t *testing.T) {
	t.Parallel()

	for _, x := range []struct { //nolint:paralleltest //  linter is buggy
		layer  Layer
		name   string
		strict bool
	}{
		{
			name:   "strict",
			layer:  &StrictStructLayer{},
			strict: true,
		},
		{
			name:   "normal",
			layer:  &StructLayer{},
			strict: false,
		},
	} {
		x := x

		t.Run(
			fmt.Sprintf("%s invalid type", x.name),
			func(t *testing.T) {
				t.Parallel()

				lb := newLayerBuilder(1)
				err := x.layer.register(lb)
				require.Error(t, err)
				require.Len(t, lb.builders, 0)
			},
		)
	}
}

func TestWithStringValueProvider(t *testing.T) {
	options := []Option{NewMockOption(t), NewMockOption(t)}
	p := NewMockNamedStringValuesProvider(t)

	l := WithStringValueProvider(p, options...)

	require.Equal(t, options, l.options)
	require.Equal(t, p, l.provider)
}

func TestWithStrictStringValueProvider(t *testing.T) {
	options := []Option{NewMockOption(t), NewMockOption(t)}
	p := NewMockNamedStringValuesProvider(t)

	l := WithStrictStringValueProvider(p, options...)

	require.Equal(t, options, l.options)
	require.Equal(t, p, l.provider)
}

func TestStrictStringProviderLayer_register(t *testing.T) {
	t.Parallel()

	svs := svalue.Values{
		"v1": &svalue.Value{
			Location: "l1",
			Value:    "v1",
		},
	}

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			p := NewMockNamedStringValuesProvider(t)
			p.
				On("GetName").
				Return("name").
				Once()
			p.
				On("GetStringValues").
				Return(svs, nil).
				Once()

			x := &StrictStringProviderLayer{
				StringProviderLayer: StringProviderLayer{
					provider: p,
				},
			}

			lb := newLayerBuilder(1)
			err := x.register(lb)

			require.NoError(t, err)
			require.Len(t, lb.builders, 1)
			require.True(t, lb.builders[0].isStrict())
			require.Contains(t, lb.idDedup, "stringProvider(name)")
		},
	)

	t.Run(
		"duplicate id",
		func(t *testing.T) {
			t.Parallel()

			p := NewMockNamedStringValuesProvider(t)
			p.
				On("GetName").
				Return("name").
				Once()

			x := &StrictStringProviderLayer{
				StringProviderLayer: StringProviderLayer{
					provider: p,
				},
			}

			lb := newLayerBuilder(1)
			lb.idDedup["stringProvider(name)"] = 101
			err := x.register(lb)

			require.Len(t, lb.builders, 0)

			var e DuplicateStringProviderError
			require.ErrorAs(t, err, &e)
			require.Equal(
				t,
				DuplicateStringProviderError{
					Index: 101,
					ID:    "name",
				},
				e,
			)
			require.Len(t, lb.idDedup, 1)
		},
	)

	t.Run(
		"option error",
		func(t *testing.T) {
			t.Parallel()

			p := NewMockNamedStringValuesProvider(t)
			p.
				On("GetName").
				Return("name").
				Once()

			o := NewMockOption(t)
			o.On(
				"apply",
				mock.MatchedBy(
					func(opts *internalOpts) bool {
						return assert.NotNil(t, opts)
					},
				),
			).Return(errMocked1)

			x := &StrictStringProviderLayer{
				StringProviderLayer: StringProviderLayer{
					provider: p,
					options:  []Option{o},
				},
			}

			lb := newLayerBuilder(1)
			err := x.register(lb)

			var e ierror.IError
			require.ErrorAs(t, err, &e)
			require.ErrorIs(t, e.Err, errMocked1)
			require.Equal(t, 0, e.Index)
			require.Len(t, lb.builders, 0)
			require.Len(t, lb.idDedup, 1)
		},
	)
}

func TestStringProviderLayer_register(t *testing.T) {
	t.Parallel()

	svs := svalue.Values{
		"v1": &svalue.Value{
			Location: "l1",
			Value:    "v1",
		},
	}

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			p := NewMockNamedStringValuesProvider(t)
			p.
				On("GetName").
				Return("name").
				Once()
			p.
				On("GetStringValues").
				Return(svs, nil).
				Once()

			x := &StringProviderLayer{
				provider: p,
			}

			lb := newLayerBuilder(1)
			err := x.register(lb)

			require.NoError(t, err)
			require.Len(t, lb.builders, 1)
			require.False(t, lb.builders[0].isStrict())
			require.Contains(t, lb.idDedup, "stringProvider(name)")
		},
	)

	t.Run(
		"duplicate id",
		func(t *testing.T) {
			t.Parallel()

			p := NewMockNamedStringValuesProvider(t)
			p.
				On("GetName").
				Return("name").
				Once()

			x := &StringProviderLayer{
				provider: p,
			}

			lb := newLayerBuilder(1)
			lb.idDedup["stringProvider(name)"] = 101
			err := x.register(lb)

			require.Len(t, lb.builders, 0)

			var e DuplicateStringProviderError
			require.ErrorAs(t, err, &e)
			require.Equal(
				t,
				DuplicateStringProviderError{
					Index: 101,
					ID:    "name",
				},
				e,
			)
			require.Len(t, lb.idDedup, 1)
		},
	)

	t.Run(
		"option error",
		func(t *testing.T) {
			t.Parallel()

			p := NewMockNamedStringValuesProvider(t)
			p.
				On("GetName").
				Return("name").
				Once()

			o := NewMockOption(t)
			o.On(
				"apply",
				mock.MatchedBy(
					func(opts *internalOpts) bool {
						return assert.NotNil(t, opts)
					},
				),
			).Return(errMocked1)

			x := &StringProviderLayer{
				provider: p,
				options:  []Option{o},
			}

			lb := newLayerBuilder(1)
			err := x.register(lb)

			var e ierror.IError
			require.ErrorAs(t, err, &e)
			require.ErrorIs(t, e.Err, errMocked1)
			require.Equal(t, 0, e.Index)
			require.Len(t, lb.builders, 0)
			require.Len(t, lb.idDedup, 1)
		},
	)
}

func TestDuplicateStringProviderError_Error(t *testing.T) {
	require.Equal(
		t,
		"string provider layer #101 is using same id=\"<name>\"",
		DuplicateStringProviderError{
			Index: 101,
			ID:    "<name>",
		}.Error(),
	)
}
