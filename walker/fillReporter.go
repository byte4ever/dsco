package walker

import (
	"fmt"
	"strings"
)

type FillReporterImpl struct {
	report    FillReport
	errors    FillErrors
	locations map[uint]string
}

func (f *FillReporterImpl) Failed() bool {
	return len(f.errors) > 0
}

func (f *FillReporterImpl) ReportError(err error) {
	f.errors = append(f.errors, err)
}

func (f *FillReporterImpl) ReportOverride(
	uid uint,
	location string,
) {
	if prevLocation, found := f.locations[uid]; found {
		f.ReportError(
			fmt.Errorf(
				"%s overrided by %s: %w",
				location,
				prevLocation,
				ErrOverriddenKey,
			),
		)
	}
}

func (f *FillReporterImpl) Result() (FillReport, error) {
	if len(f.errors) > 0 {
		return nil, f.errors
	}

	return f.report, nil
}

func NewFillReporterImpl() *FillReporterImpl {
	return &FillReporterImpl{
		locations: make(map[uint]string),
	}
}

func (f *FillReporterImpl) ReportUse(
	uid uint,
	path string,
	location string,
) {
	f.locations[uid] = location
	f.report = append(
		f.report, FillReportEntry{
			path:     path,
			location: location,
		},
	)
}

func (f *FillReporterImpl) ReportUnused(path string) {
	f.ReportError(
		fmt.Errorf(
			"%s: %w",
			path,
			ErrUninitializedKey,
		),
	)
}

// FillErrors represents the list of errors that occur when filling the
// structure.
type FillErrors []error

func (f FillErrors) Error() string {
	var sb strings.Builder
	for _, err := range f {
		sb.WriteString(err.Error())
		sb.WriteRune('\n')
	}

	return sb.String()
}
