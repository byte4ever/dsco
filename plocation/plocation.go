package plocation

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PathLocations contains all fillHelper Location for every key path.
type PathLocations []PathLocation

// PathLocation is the fillHelper report for a value.
type PathLocation struct {
	UID      uint
	Path     string // is the key path.
	Location string // is the Location of the value.
}

// Dump writes fillHelper report.
func (f *PathLocations) Dump(writer io.Writer) {
	tabWriter := tabwriter.NewWriter(
		writer,
		0,
		0,
		2,
		' ',
		tabwriter.Debug,
	)
	_, _ = fmt.Fprintln(tabWriter, "  path\t  Location")
	_, _ = fmt.Fprintln(tabWriter, "  ----\t  --------")

	//nolint:gocritic // don't care it is error processing
	for _, entry := range *f {
		_, _ = fmt.Fprintf(
			tabWriter, "  %s\t  %s\n", entry.Path, entry.Location,
		)
	}

	_ = tabWriter.Flush()
}

func (f *PathLocations) Report(uid uint, path string, location string) {
	*f = append(
		*f, PathLocation{
			UID:      uid,
			Path:     path,
			Location: location,
		},
	)
}

func (f *PathLocations) ReportOther(other PathLocations) {
	*f = append(
		*f, other...,
	)
}
