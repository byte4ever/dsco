package walker

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

// ErrUninitializedKey represent an error where ....
var ErrUninitializedKey = errors.New("uninitialized key")

// FillReport contains all fillHelper location for every key path.
type FillReport []FillReportEntry

// FillErrors represents the list of errors that occur when filling the
// structure.
type FillErrors []error

// FillReportEntry is the fillHelper report for a value.
type FillReportEntry struct {
	path     string // is the key path.
	location string // is the location of the value.
}

func (f FillErrors) Error() string {
	var sb strings.Builder
	for _, err := range f {
		sb.WriteString(err.Error())
		sb.WriteRune('\n')
	}

	return sb.String()
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

// Fill fills the structure using the layers.
func Fill(
	inputModel any,
	layers ...Layer,
) (
	FillReport, error,
) {
	bo := newLayerBuilder()

	for _, layer := range layers {
		err := layer.register(bo)
		if err != nil {
			return nil, err //nolint:wrapcheck // error is clear enough
		}
	}

	return fillHelper(inputModel, bo.builders...)
}

// fillHelper files the structure.
func fillHelper(
	inputModel any,
	baseProvider ...constraintLayerPolicy,
) (
	FillReport, error,
) {
	var (
		fillReport   FillReport
		fillErrors   FillErrors
		maxId        int
		bases        []Base
		mustBeUseIdx []int
	)

	for idx, builder := range baseProvider {
		base, errs := builder.GetBaseFor(inputModel)
		if len(errs) > 0 {
			for _, err := range errs {
				fillErrors = append(
					fillErrors,
					fmt.Errorf("layer #%d: %wlkr", idx, err),
				)
			}

			continue
		}

		if builder.isStrict() {
			mustBeUseIdx = append(mustBeUseIdx, len(bases))
		}

		bases = append(bases, base)
	}

	if len(fillErrors) > 0 {
		return nil, fillErrors
	}

	reportLoc := make(map[int]string)

	wlkr := walker{
		fieldAction: func(id int, path string, value *reflect.Value) error {
			for _, basis := range bases {
				v, location := basis.Get(id)
				if v != nil {
					value.Set(*v)
					reportLoc[id] = location
					fillReport = append(
						fillReport,
						FillReportEntry{
							path:     path,
							location: location,
						},
					)
					return nil
				}
			}

			fillErrors = append(
				fillErrors,
				fmt.Errorf("%s: %wlkr", path, ErrUninitializedKey),
			)

			return nil
		},
	}

	if err := wlkr.walkRec(
		&maxId,
		"",
		reflect.ValueOf(inputModel),
	); err != nil {
		return nil, FillErrors{err}
	}

	if len(mustBeUseIdx) > 0 {
		for _, idx := range mustBeUseIdx {
			//nolint:gocritic // don't care
			for valId, e := range bases[idx] {
				fillErrors = append(
					fillErrors, fmt.Errorf(
						"%s overrided by %s: %wlkr",
						e.location,
						reportLoc[valId],
						ErrOverriddenKey,
					),
				)
			}
		}
	}

	if len(fillErrors) > 0 {
		return nil, fillErrors
	}

	return fillReport, nil
}
