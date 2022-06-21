package walker

import (
	"fmt"
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
		f.errors = append(
			f.errors, fmt.Errorf(
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
	f.errors = append(
		f.errors, fmt.Errorf(
			"%s: %w",
			path,
			ErrUninitializedKey,
		),
	)
}
