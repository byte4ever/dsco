package walker

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PathLocations contains all fillHelper location for every key path.
type PathLocations []PathLocation

// PathLocation is the fillHelper report for a value.
type PathLocation struct {
	UID      uint
	Path     string // is the key path.
	location string // is the location of the value.
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
	_, _ = fmt.Fprintln(tabWriter, "  path\t  location")
	_, _ = fmt.Fprintln(tabWriter, "  ----\t  --------")

	//nolint:gocritic // don't care it is error processing
	for _, entry := range *f {
		_, _ = fmt.Fprintf(
			tabWriter, "  %s\t  %s\n", entry.Path, entry.location,
		)
	}

	_ = tabWriter.Flush()
}

func (f *PathLocations) Report(uid uint, path string, location string) {
	*f = append(
		*f, PathLocation{
			UID:      uid,
			Path:     path,
			location: location,
		},
	)
}

func (f *PathLocations) ReportOther(other PathLocations) {
	*f = append(
		*f, other...,
	)
}
