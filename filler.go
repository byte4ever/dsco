package dsco

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/internal/fvalue"
	"github.com/byte4ever/dsco/internal/merror"
	model2 "github.com/byte4ever/dsco/internal/model"
	"github.com/byte4ever/dsco/internal/plocation"
)

type dscoContext struct {
	inputModelRef any
	err           FillerErrors
	layers        PoliciesGetter

	// ----
	model            ModelInterface
	builders         constraintLayerPolicies
	layerFieldValues []fvalue.Values
	mustBeUsed       []int
	pathLocations    plocation.Locations
}

type FillerErrors struct {
	merror.MError
}

var ErrFiller = errors.New("")

func (FillerErrors) Is(err error) bool {
	return errors.Is(err, ErrFiller)
}

func newDSCOContext(
	inputModelRef any,
	layers Layers,
) *dscoContext {
	return &dscoContext{
		inputModelRef: inputModelRef,
		layers:        layers,
	}
}

func (c *dscoContext) generateModel() {
	if c.err.None() {
		model, err := model2.NewModel(reflect.TypeOf(c.inputModelRef).Elem())
		if err != nil {
			c.err.Add(err)
			return
		}

		c.model = model
	}
}

func (c *dscoContext) generateBuilders() {
	if c.err.None() {
		var err error
		c.builders, err = c.layers.GetPolicies()

		if err != nil {
			c.err.Add(err)
		}
	}
}

func (c *dscoContext) generateFieldValues() {
	if c.err.None() {
		for idx, builder := range c.builders {
			base, err := builder.GetFieldValuesFrom(c.model)
			if err != nil {
				c.err.Add(
					fmt.Errorf(
						"layer #%d\n %w",
						idx,
						err,
					),
				)

				continue
			}

			if builder.isStrict() {
				c.mustBeUsed = append(c.mustBeUsed, len(c.layerFieldValues))
			}

			c.layerFieldValues = append(c.layerFieldValues, base)
		}
	}
}

func (c *dscoContext) fillIt() {
	if c.err.None() {
		v := reflect.ValueOf(c.inputModelRef).Elem()

		pathLocations, err := c.model.Fill(v, c.layerFieldValues)
		if err != nil {
			c.err.Add(err)
			return
		}

		c.pathLocations = pathLocations
	}
}

func (c *dscoContext) checkUnused() {
	if c.err.None() {
		for _, idx := range c.mustBeUsed {
			for valUID, e := range c.layerFieldValues[idx] {
				c.err.Add(
					OverriddenKeyError{
						Path:             c.pathLocations[valUID].Path,
						Location:         e.Location,
						OverrideLocation: c.pathLocations[valUID].Location,
					},
				)
			}
		}
	}
}

// Fill fills the structure using the layers.
func Fill(
	inputModelRef any,
	layers ...Layer,
) (
	plocation.Locations,
	error,
) {
	fillContext := newDSCOContext(inputModelRef, layers)

	fillContext.generateModel()
	fillContext.generateBuilders()
	fillContext.generateFieldValues()
	fillContext.fillIt()
	fillContext.checkUnused()

	if fillContext.err.None() {
		return fillContext.pathLocations, nil
	}

	return fillContext.pathLocations, fillContext.err //nolint:wrapcheck // ok
}
