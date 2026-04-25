// Package dsco provides configuration processing using layered sources.
// This file contains the core filling logic that orchestrates the
// configuration process from source layers to final struct population.
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

// dscoContext encapsulates the state and orchestrates the configuration
// filling process across multiple phases: model generation, builder
// creation, field value extraction, struct filling, and validation.
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

// FillerErrors aggregates multiple errors that can occur during the
// configuration filling process across different layers and validation steps.
type FillerErrors struct {
	merror.MError
}

var (
	// ErrFiller is the sentinel error for configuration filling failures.
	ErrFiller = errors.New("")

	// ErrCfgMustBePointer indicates that the value passed to Fill or
	// BuildModel is not a pointer to a struct.
	ErrCfgMustBePointer = errors.New("cfg must be a pointer to a struct")
)

// Is implements error matching for FillerErrors, allowing error.Is checks
// against ErrFiller to detect filling-related errors.
func (FillerErrors) Is(err error) bool {
	return errors.Is(err, ErrFiller)
}

// newDSCOContext creates a new configuration filling context with the
// target struct reference and source layers to process.
func newDSCOContext(
	inputModelRef any,
	layers Layers,
) *dscoContext {
	return &dscoContext{
		inputModelRef: inputModelRef,
		layers:        layers,
	}
}

// generateModel populates c.model from c.inputModelRef.
// c.inputModelRef is **T (Fill is called with &pp where pp is *T), so
// dereference once to obtain *T before delegating to buildModel.
func (c *dscoContext) generateModel() {
	if c.err.None() {
		// Dereference **T → *T so buildModel receives a plain *Struct.
		ptrVal := reflect.ValueOf(c.inputModelRef).Elem().Interface()

		mdl, err := buildModel(ptrVal)
		if err != nil {
			c.err.Add(err)
			return
		}

		c.model = mdl
	}
}

// buildModel constructs the model from a pointer-to-struct configuration.
// inputModelRef must be *T where T is a struct. Used by both Fill (via
// dscoContext.generateModel) and inventory walks.
//
//nolint:iface,ireturn,revive // shared helper; ModelInterface is the contract used by both Fill and inventory
func buildModel(inputModelRef any) (ModelInterface, error) {
	const errCtx = "building model"

	t := reflect.TypeOf(inputModelRef)
	if t == nil || t.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("%s: %w", errCtx, ErrCfgMustBePointer)
	}

	mdl, err := model2.NewModel(t)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCtx, err)
	}

	return mdl, nil
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
