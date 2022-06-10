package dsco

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrUninitialized is returned when a field cannot be filled with a value.
var ErrUninitialized = errors.New("uninitialized")

// ReportEntry is a report entry.
type ReportEntry struct {
	Value       reflect.Value
	Key         string
	ExternalKey string
	Idx         int
	Errors      []error
}

// Report is a binding report.
type Report []ReportEntry

var _ reportInterface = &Report{}

func (re *ReportEntry) isFound() bool {
	return re.Idx != -1
}

func (r Report) perEntryReport() (errs []error) {
	for _, entry := range r { //nolint:gocritic // will be refactored soon
		for _, err := range entry.Errors {
			if err != nil {
				errs = append(errs, err)
			}
		}

		if entry.Idx == -1 {
			errs = append(
				errs,
				fmt.Errorf("key <%v>: %w", entry.Key, ErrUninitialized),
			)
		}
	}

	return errs
}

func (r *Report) addEntry(e ReportEntry) {
	*r = append(*r, e)
}
