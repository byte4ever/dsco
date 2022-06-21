package walker

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

// FillError represents the list of errors that occur when filling the
// structure.
type FillError []error

func (f FillError) Error() string {
	var sb strings.Builder
	for _, err := range f {
		sb.WriteString(err.Error())
		sb.WriteRune('\n')
	}

	return sb.String()
}

// ErrUninitializedKey represent an error where ....
var ErrUninitializedKey = errors.New("uninitialized key")

// FillReport contains all fill location for every key path.
type FillReport []FillReportEntry

// Dump writes fill report.
func (f FillReport) Dump(writer io.Writer) {
	w := tabwriter.NewWriter(
		writer,
		0,
		0,
		2,
		' ',
		tabwriter.Debug,
	)
	_, _ = fmt.Fprintln(w, "  path\t  location")
	_, _ = fmt.Fprintln(w, "  ----\t  --------")

	for _, entry := range f {
		_, _ = fmt.Fprintf(w, "  %s\t  %s\n", entry.path, entry.location)
	}

	_ = w.Flush()
}

// FillReportEntry is the fill report for a value.
type FillReportEntry struct {
	path     string // is the key path.
	location string // is the location of the value.
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
			return nil, err
		}
	}
	return fill(inputModel, bo.builders...)
}

// fill files the structure.
func fill(
	inputModel any,
	baseProvider ...ConstraintLayerPolicy,
) (
	FillReport, error,
) {
	var (
		fillReport   FillReport
		fillErrors   FillError
		maxId        int
		bases        []Base
		mustBeUseIdx []int
	)

	for idx, builder := range baseProvider {
		b, e := builder.GetBaseFor(inputModel)
		if len(e) > 0 {
			for _, err := range e {
				fillErrors = append(
					fillErrors,
					fmt.Errorf("layer #%d: %w", idx, err),
				)
			}
			continue
		}

		if builder.IsStrict() {
			mustBeUseIdx = append(mustBeUseIdx, len(bases))
		}

		bases = append(bases, b)
	}

	if len(fillErrors) > 0 {
		return nil, fillErrors
	}

	reportLoc := make(map[int]string)

	w := walker{
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
				fmt.Errorf("%s: %w", path, ErrUninitializedKey),
			)

			return nil
		},
	}

	if err := w.walkRec(
		&maxId,
		"",
		reflect.ValueOf(inputModel),
	); err != nil {
		return nil, FillError{err}
	}

	if len(mustBeUseIdx) > 0 {
		for _, idx := range mustBeUseIdx {
			for valId, e := range bases[idx] {
				fillErrors = append(
					fillErrors, fmt.Errorf(
						"%s overrided by %s: %w",
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
