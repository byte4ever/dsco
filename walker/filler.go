package walker

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/merror"
)

// ErrUninitializedKey represent an error where ....
var ErrUninitializedKey = errors.New("uninitialized key")

type dscoContext struct {
	inputModelRef any
	err           FillerErrors
	layers        Layers

	// ----
	model            ModelInterface
	builders         constraintLayerPolicies
	layerFieldValues []FieldValues
	mustBeUsed       []int
	pathLocations    PathLocations
}

type FillerErrors struct {
	merror.MError
}

var ErrFiller = errors.New("")

func (m FillerErrors) Is(err error) bool {
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
		model, err := NewModel(reflect.TypeOf(c.inputModelRef).Elem())
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
			base, err := builder.GetFieldValues(c.model)
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
		if len(c.mustBeUsed) > 0 {
			for _, idx := range c.mustBeUsed {
				for valUID, e := range c.layerFieldValues[idx] {
					c.err.Add(
						fmt.Errorf(
							"%s %s by %s: %w",
							c.pathLocations[valUID].Path,
							e.location,
							c.pathLocations[valUID].location,
							ErrOverriddenKey,
						),
					)
				}
			}
		}
	}
}

// Fill fills the structure using the layers.
func Fill(
	inputModelRef any,
	layers ...Layer,
) (
	PathLocations,
	error,
) {
	c := newDSCOContext(inputModelRef, layers)

	c.generateModel()
	c.generateBuilders()
	c.generateFieldValues()
	c.fillIt()
	c.checkUnused()

	if c.err.None() {
		return c.pathLocations, nil
	}

	return c.pathLocations, c.err //nolint:wrapcheck
}
