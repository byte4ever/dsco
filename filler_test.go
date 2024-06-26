package dsco

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco/internal/fvalue"
	"github.com/byte4ever/dsco/internal/merror"
	"github.com/byte4ever/dsco/internal/plocation"
	"github.com/byte4ever/dsco/ref"
)

func Test_newDSCOContext(t *testing.T) {
	t.Parallel()

	v := 15
	layers := Layers{}

	c := newDSCOContext(
		v,
		layers,
	)

	require.Equal(t, v, c.inputModelRef)
	require.Equal(t, layers, c.layers)
}

func Test_dscoContext_generateModel(t *testing.T) {
	t.Parallel()

	t.Run(
		"skip step",
		func(t *testing.T) {
			t.Parallel()

			c := &dscoContext{
				err: FillerErrors{
					MError: merror.MError{errMocked1},
				},
			}

			c.generateModel()
		},
	)

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			type Root struct {
				X *float64
				Y *float32
			}

			var v *Root

			c := &dscoContext{
				inputModelRef: &v,
			}

			c.generateModel()
			require.NotNil(t, c.model)
			require.True(t, c.err.None())
		},
	)

	t.Run(
		"failure",
		func(t *testing.T) {
			t.Parallel()

			type Root struct {
				Y *float32
				X float64
			}

			var v *Root

			c := &dscoContext{
				inputModelRef: &v,
			}

			c.generateModel()
			require.Nil(t, c.model)
			require.False(t, c.err.None())
		},
	)
}

func Test_dscoContext_generateBuilders(t *testing.T) {
	t.Parallel()

	t.Run(
		"skip step",
		func(t *testing.T) {
			t.Parallel()

			c := &dscoContext{
				err: FillerErrors{
					MError: merror.MError{errMocked1},
				},
			}

			c.generateBuilders()
		},
	)

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			builders := constraintLayerPolicies{
				newMockConstraintLayerPolicy(t),
				newMockConstraintLayerPolicy(t),
				newMockConstraintLayerPolicy(t),
			}

			pg := NewMockPoliciesGetter(t)
			pg.
				On("GetPolicies").
				Return(builders, nil).
				Once()

			c := &dscoContext{
				layers: pg,
			}

			c.generateBuilders()
			require.Equal(t, builders, c.builders)
			require.True(t, c.err.None())
		},
	)

	t.Run(
		"failure",
		func(t *testing.T) {
			t.Parallel()

			pg := NewMockPoliciesGetter(t)
			pg.
				On("GetPolicies").
				Return(nil, errMocked1).
				Once()

			c := &dscoContext{
				layers: pg,
			}

			c.generateBuilders()
			require.Nil(t, c.builders)
			require.False(t, c.err.None())
			require.Len(t, c.err.MError, 1)
			require.ErrorIs(t, c.err.MError[0], errMocked1)
		},
	)
}

func Test_dscoContext_generateFieldValues(t *testing.T) {
	t.Parallel()

	t.Run(
		"skip step",
		func(t *testing.T) {
			t.Parallel()

			c := &dscoContext{
				err: FillerErrors{
					MError: merror.MError{errMocked1},
				},
			}

			c.generateFieldValues()
		},
	)

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			model := NewMockModelInterface(t)

			fvs1 := fvalue.Values{
				200: &fvalue.Value{
					Value:    reflect.Value{},
					Location: "loc1",
				},
			}

			clp1 := newMockConstraintLayerPolicy(t)
			clp1.
				On("GetFieldValuesFrom", model).
				Return(fvs1, nil).
				Once()

			clp1.
				On("isStrict").
				Return(false).
				Once()

			fvs2 := fvalue.Values{
				400: &fvalue.Value{
					Value:    reflect.Value{},
					Location: "loc2",
				},
			}

			clp2 := newMockConstraintLayerPolicy(t)
			clp2.
				On("GetFieldValuesFrom", model).
				Return(fvs2, nil).
				Once()

			clp2.
				On("isStrict").
				Return(true).
				Once()

			builders := constraintLayerPolicies{
				clp1,
				clp2,
			}

			c := &dscoContext{
				builders: builders,
				model:    model,
			}

			c.generateFieldValues()
			require.Equal(
				t,
				[]fvalue.Values{fvs1, fvs2},
				c.layerFieldValues,
			)
			require.Equal(
				t,
				[]int{1},
				c.mustBeUsed,
			)
			require.True(t, c.err.None())
		},
	)

	t.Run(
		"failure",
		func(t *testing.T) {
			t.Parallel()

			model := NewMockModelInterface(t)

			clp1 := newMockConstraintLayerPolicy(t)
			clp1.
				On("GetFieldValuesFrom", model).
				Return(nil, errMocked1).
				Once()

			fvs2 := fvalue.Values{
				400: &fvalue.Value{
					Value:    reflect.Value{},
					Location: "loc2",
				},
			}

			clp2 := newMockConstraintLayerPolicy(t)
			clp2.
				On("GetFieldValuesFrom", model).
				Return(fvs2, nil).
				Once()

			clp2.
				On("isStrict").
				Return(true).
				Once()

			builders := constraintLayerPolicies{
				clp1,
				clp2,
			}

			c := &dscoContext{
				builders: builders,
				model:    model,
			}

			c.generateFieldValues()
			require.Equal(
				t,
				[]fvalue.Values{fvs2},
				c.layerFieldValues,
			)
			require.Equal(
				t,
				[]int{0},
				c.mustBeUsed,
			)
			require.False(t, c.err.None())
			require.Len(t, c.err.MError, 1)
			require.ErrorIs(t, c.err.MError[0], errMocked1)
		},
	)
}

func Test_dscoContext_fillIt(t *testing.T) {
	t.Parallel()

	t.Run(
		"skip step",
		func(t *testing.T) {
			t.Parallel()

			c := &dscoContext{
				err: FillerErrors{
					MError: merror.MError{errMocked1},
				},
			}

			c.fillIt()
		},
	)

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			type Root struct {
				X *float64
				Y *float32
			}

			var v *Root
			pv := &v

			ve := reflect.ValueOf(pv).Elem()

			base := []fvalue.Values{
				{
					200: &fvalue.Value{
						Value:    reflect.Value{},
						Location: "loc1",
					},
				},
				{
					400: &fvalue.Value{
						Value:    reflect.Value{},
						Location: "loc2",
					},
				},
			}

			ploc := plocation.Locations{
				plocation.Location{
					Path: "p1",
				},
				plocation.Location{
					Path: "p2",
				},
				plocation.Location{
					Path: "p3",
				},
			}

			model := NewMockModelInterface(t)
			model.
				On("Fill", ve, base).
				Return(ploc, nil).
				Once()

			c := &dscoContext{
				inputModelRef:    &v,
				model:            model,
				layerFieldValues: base,
			}

			c.fillIt()
			require.True(t, c.err.None())
			require.Equal(t, ploc, c.pathLocations)
		},
	)

	t.Run(
		"failure",
		func(t *testing.T) {
			t.Parallel()

			type Root struct {
				X *float64
				Y *float32
			}

			var v *Root
			pv := &v

			ve := reflect.ValueOf(pv).Elem()

			base := []fvalue.Values{
				{
					200: &fvalue.Value{
						Value:    reflect.Value{},
						Location: "loc1",
					},
				},
				{
					400: &fvalue.Value{
						Value:    reflect.Value{},
						Location: "loc2",
					},
				},
			}

			model := NewMockModelInterface(t)
			model.
				On("Fill", ve, base).
				Return(nil, errMocked1).
				Once()

			c := &dscoContext{
				inputModelRef:    &v,
				model:            model,
				layerFieldValues: base,
			}

			c.fillIt()
			require.Nil(t, c.pathLocations)
			require.False(t, c.err.None())
			require.Len(t, c.err.MError, 1)
			require.ErrorIs(t, c.err.MError[0], errMocked1)
		},
	)
}

func Test_dscoContext_checkUnused(t *testing.T) {
	t.Parallel()

	t.Run(
		"skip step",
		func(t *testing.T) {
			t.Parallel()

			c := &dscoContext{
				err: FillerErrors{
					MError: merror.MError{errMocked1},
				},
			}

			c.checkUnused()
		},
	)

	t.Run(
		"failure",
		func(t *testing.T) {
			t.Parallel()

			base := []fvalue.Values{
				{
					0: &fvalue.Value{
						Value:    reflect.Value{},
						Location: "loc0",
					},
				},
				{
					1: &fvalue.Value{
						Value:    reflect.Value{},
						Location: "loc1",
					},
				},
			}

			c := &dscoContext{
				pathLocations: plocation.Locations{
					plocation.Location{
						Path:     "path0",
						Location: "foundLoc0",
					},
					plocation.Location{
						Path:     "path1",
						Location: "foundLoc1",
					},
				},
				layerFieldValues: base,
				mustBeUsed:       []int{1},
			}

			c.checkUnused()
			require.False(t, c.err.None())
			require.Len(t, c.err.MError, 1)

			var e OverriddenKeyError

			require.ErrorAs(t, c.err.MError[0], &e)
			require.Equal(
				t, OverriddenKeyError{
					Path:             "path1",
					Location:         "loc1",
					OverrideLocation: "foundLoc1",
				}, e,
			)
		},
	)

	t.Run(
		"succes no strict",
		func(t *testing.T) {
			t.Parallel()

			base := []fvalue.Values{
				{
					0: &fvalue.Value{
						Value:    reflect.Value{},
						Location: "loc0",
					},
				},
				{
					1: &fvalue.Value{
						Value:    reflect.Value{},
						Location: "loc1",
					},
				},
			}

			c := &dscoContext{
				pathLocations: plocation.Locations{
					plocation.Location{
						Path:     "path0",
						Location: "foundLoc0",
					},
					plocation.Location{
						Path:     "path1",
						Location: "foundLoc1",
					},
				},
				layerFieldValues: base,
				mustBeUsed:       nil,
			}

			c.checkUnused()
			require.True(t, c.err.None())
		},
	)

	t.Run(
		"success no remaining values",
		func(t *testing.T) {
			t.Parallel()

			base := []fvalue.Values{
				{
					0: &fvalue.Value{
						Value:    reflect.Value{},
						Location: "loc0",
					},
				},
				{},
			}

			c := &dscoContext{
				pathLocations: plocation.Locations{
					plocation.Location{
						Path:     "path0",
						Location: "foundLoc0",
					},
					plocation.Location{
						Path:     "path1",
						Location: "foundLoc1",
					},
				},
				layerFieldValues: base,
				mustBeUsed:       []int{1},
			}

			c.checkUnused()
			require.True(t, c.err.None())
		},
	)
}

//nolint:paralleltest // using global variables
func TestFill(t *testing.T) {
	t.Run(
		"success",
		func(t *testing.T) {

			type Sub struct {
				FirstName    *string
				LastName     *string
				TrainingTime *time.Duration
				T            *time.Time
				B            *bool
			}

			type Root struct {
				A    *int
				B    *float64
				T    *time.Time
				Z    *Sub
				NaNa *int
				L    []string
			}

			os.Args = []string{
				"appName",
				"--a=1234",
				"--z-b=yes",
				"--z-first_name=Laura",
			}

			// t.Setenv("TST-A", "123")
			t.Setenv("TST-B", "123.1234")
			// t.Setenv("API-Z-FIRST_NAME", "Laurent")

			var pp *Root
			fillReport, err := Fill(
				&pp,
				WithEnvLayer("API"),
				WithEnvLayer("TST"),
				WithStrictCmdlineLayer(),
				WithStructLayer(
					&Root{
						B: ref.R(0.0),
						Z: &Sub{
							FirstName: ref.R("Rose"),
							LastName:  ref.R("Dupont"),
							B:         ref.R(false),
						},
					}, "dflt1",
				),
				WithStructLayer(
					&Root{
						A: ref.R(120),
						B: ref.R(2333.32),
						T: ref.R(time.Now().UTC()),
						Z: &Sub{
							FirstName:    ref.R("Lola"),
							LastName:     ref.R("MARTIN"),
							TrainingTime: ref.R(800 * time.Second),
							T:            ref.R(time.Now().UTC()),
							B:            ref.R(true),
						},
						NaNa: ref.R(2331),
						L:    []string{"A", "B", "C"},
					}, "dflt2",
				),
			)

			require.NoError(t, err)

			bb, err := yaml.Marshal(pp)
			require.NoError(t, err)

			t.Log(string(bb))

			fillReport.Dump(os.Stdout)
		},
	)

	t.Run(
		"failure",
		func(t *testing.T) {

			type Sub struct {
				FirstName    *string
				TrainingTime *time.Duration
				T            *time.Time
				B            *bool
				LastName     string
			}

			type Root struct {
				A    *int
				B    *float64
				T    *time.Time
				Z    *Sub
				NaNa *int
				L    []string
			}

			os.Args = []string{
				"appName",
				"--a=1234",
				"--z-b=yes",
				"--z-first_name=Laura",
			}

			// t.Setenv("TST-A", "123")
			t.Setenv("TST-B", "123.1234")
			// t.Setenv("API-Z-FIRST_NAME", "Laurent")

			var pp *Root
			_, err := Fill(
				&pp,
				WithEnvLayer("API"),
				WithEnvLayer("TST"),
				WithStrictCmdlineLayer(),
				WithStructLayer(
					&Root{
						B: ref.R(0.0),
						Z: &Sub{
							FirstName: ref.R("Rose"),
							B:         ref.R(false),
						},
					}, "dflt1",
				),
				WithStructLayer(
					&Root{
						A: ref.R(120),
						B: ref.R(2333.32),
						T: ref.R(time.Now().UTC()),
						Z: &Sub{
							FirstName:    ref.R("Lola"),
							TrainingTime: ref.R(800 * time.Second),
							T:            ref.R(time.Now().UTC()),
							B:            ref.R(true),
						},
						NaNa: ref.R(2331),
						L:    []string{"A", "B", "C"},
					}, "dflt2",
				),
			)

			require.Error(t, err)
		},
	)
}

func TestFillerErrors_Is(t *testing.T) {
	t.Parallel()

	require.ErrorIs(t, FillerErrors{}, ErrFiller)
}
