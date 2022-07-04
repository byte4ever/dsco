package dsco

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/ierror"
	"github.com/byte4ever/dsco/utils"
)

func TestLayers_GetPolicies(t *testing.T) {
	t.Parallel()

	t.Run(
		"success", func(t *testing.T) {
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
		"fail", func(t *testing.T) {

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

			var e LayerErrors

			require.ErrorAs(t, err, &e)

			for idx, et := range []struct {
				ie int
				e  error
			}{
				{0, errMocked1},
				{2, errMocked2},
			} {
				var ie ierror.IError
				require.ErrorAs(t, e.MError[idx], &ie)
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

func TestEnvLayer_register(t *testing.T) {
	t.Parallel()

	for _, x := range []struct {
		layer  Layer
		strict bool
		name   string
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

func TestStructLayer_register(t *testing.T) {
	t.Parallel()

	type Root struct {
		X *float32
		Y *float64
	}

	k := &Root{
		X: utils.R(float32(123.123)),
	}

	for _, x := range []struct {
		layer  Layer
		strict bool
		name   string
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
				os.Args = []string{"cmdName"}

				lb := newLayerBuilder(1)
				err := x.layer.register(lb)

				require.NoError(t, err)
				require.Len(t, lb.builders, 1)
				require.Equal(t, x.strict, lb.builders[0].isStrict())
			},
		)

		t.Run(
			fmt.Sprintf("%s using same id", x.name),
			func(t *testing.T) {
				os.Args = []string{"cmdName"}

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

	type Root struct {
		X *float32
		Y *float64
	}

	for _, x := range []struct {
		layer  Layer
		strict bool
		name   string
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

func TestInvalidInputError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"type int is not a valid pointer on struct",
		InvalidInputError{
			Type: reflect.TypeOf(10),
		}.Error(),
	)
}

func TestInvalidInputError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrInvalidInput)
	require.ErrorIs(t, InvalidInputError{}, ErrInvalidInput)
}

func TestCmdlineAlreadyUsedError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"cmdline already used in position #101",
		CmdlineAlreadyUsedError{
			Index: 101,
		}.Error(),
	)
}

func TestCmdlineAlreadyUsedError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrCmdlineAlreadyUsed)
	require.ErrorIs(t, CmdlineAlreadyUsedError{}, ErrCmdlineAlreadyUsed)
}

func TestDuplicateEnvPrefixError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"layer #101 has same prefix=PREFIX",
		DuplicateEnvPrefixError{
			Index:  101,
			Prefix: "PREFIX",
		}.Error(),
	)
}

func TestDuplicatAeEnvPrefixError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrDuplicateEnvPrefix)
	require.ErrorIs(t, DuplicateEnvPrefixError{}, ErrDuplicateEnvPrefix)
}

func TestDuplicateInputStructError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"struct layer #101 is using same pointer",
		DuplicateInputStructError{
			Index: 101,
		}.Error(),
	)
}

func TestDuplicateInputStructError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrDuplicateInputStruct)
	require.ErrorIs(t, DuplicateInputStructError{}, ErrDuplicateInputStruct)
}

func TestDuplicateStructIDError_Error(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"struct layer #101 is using same id=\"OTHER\"",
		DuplicateStructIDError{
			Index: 101,
			ID:    "OTHER",
		}.Error(),
	)
}

func TestDuplicateStructIDError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrDuplicateStructID)
	require.ErrorIs(t, DuplicateStructIDError{}, ErrDuplicateStructID)
}

func TestLayerErrors_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrLayer)
	require.ErrorIs(t, LayerErrors{}, ErrLayer)
}
