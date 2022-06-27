package walker

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"
)

// ErrUninitializedKey represent an error where ....
var ErrUninitializedKey = errors.New("uninitialized key")

// FillReport contains all fillHelper location for every key path.
type FillReport []FillReportEntry

// FillReportEntry is the fillHelper report for a value.
type FillReportEntry struct {
	path     string // is the key path.
	location string // is the location of the value.
}

// Dump writes fillHelper report.
func (f FillReport) Dump(writer io.Writer) {
	tabWriter := tabwriter.NewWriter(
		writer,
		0,
		0,
		2,
		' ',
		tabwriter.Debug,
	)
	_, _ = fmt.Fprintln(tabWriter, "  path\t  location")
	_, _ = fmt.Fprintln(tabWriter, "  ----\t  --------")

	//nolint:gocritic // don't care it is error processing
	for _, entry := range f {
		_, _ = fmt.Fprintf(
			tabWriter, "  %s\t  %s\n", entry.path, entry.location,
		)
	}

	_ = tabWriter.Flush()
}

type dscoContext struct {
	inputModelRef any
	reporter      FillReporter
	layers        Layers

	// ----
	model            ModelInterface
	builders         constraintLayerPolicies
	layerFieldValues []FieldValues
	mustBeUsed       []int
}

func newDSCOContext(
	inputModelRef any,
	reporter FillReporter,
	layers Layers,
) *dscoContext {
	return &dscoContext{
		inputModelRef: inputModelRef,
		reporter:      reporter,
		layers:        layers,
	}
}

func (c *dscoContext) generateModel() {
	if !c.reporter.Failed() {
		model, errs := NewModel(reflect.TypeOf(c.inputModelRef).Elem())

		if len(errs) > 0 {
			for _, err := range errs {
				c.reporter.ReportError(err)
			}

			return
		}

		c.model = model
	}
}

func (c *dscoContext) generateBuilders() {
	if !c.reporter.Failed() {
		c.builders = c.layers.GetPolicies(c.reporter)
	}
}

func (c *dscoContext) generateFieldValues() {
	if !c.reporter.Failed() {
		for idx, builder := range c.builders {
			base, errs2 := builder.GetFieldValues(c.model)
			if len(errs2) > 0 {
				for _, err2 := range errs2 {
					c.reporter.ReportError(
						fmt.Errorf(
							"layer #%d: %wlkr",
							idx,
							err2,
						),
					)
				}

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
	if !c.reporter.Failed() {
		v := reflect.ValueOf(c.inputModelRef).Elem()

		c.model.Fill(c.reporter, v, c.layerFieldValues)
	}
}

func (c *dscoContext) checkUnused() {
	if !c.reporter.Failed() {
		if len(c.mustBeUsed) > 0 {
			for _, idx := range c.mustBeUsed {
				for valUID, e := range c.layerFieldValues[idx] {
					c.reporter.ReportOverride(valUID, e.location)
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
	FillReport,
	error,
) {
	fillReporter := NewFillReporterImpl()

	c := newDSCOContext(inputModelRef, fillReporter, layers)

	c.generateModel()
	c.generateBuilders()
	c.generateFieldValues()
	c.fillIt()
	c.checkUnused()
	return c.reporter.Result() //nolint:wrapcheck
}
