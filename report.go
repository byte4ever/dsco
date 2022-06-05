package dsco

import (
	"fmt"
)

type ReportEntry struct {
	Key         string
	ExternalKey string
	Idx         int
	Errors      []error
}

func (re *ReportEntry) isFound() bool {
	return re.Idx != -1
}

var _ reportIface = &Report{}

type Report []ReportEntry

func (r Report) perEntryReport() (errs []error) {
	for _, entry := range r {
		for _, err := range entry.Errors {
			if err != nil {
				errs = append(errs, err)
			}
		}

		if entry.Idx == -1 {
			errs = append(errs, fmt.Errorf("key <%v>: %w", entry.Key, ErrUninitialized))
		}
	}

	return errs
}

func (r *Report) addEntry(e ReportEntry) {
	*r = append(*r, e)
}
