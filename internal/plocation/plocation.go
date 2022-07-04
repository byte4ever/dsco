package plocation

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Locations contains all fillHelper Location for every key path.
type Locations []Location

// Location is the fillHelper report for a value.
type Location struct {
	Path     string
	Location string
	UID      uint
}

// Dump writes filling report in writer.
func (f *Locations) Dump(writer io.Writer) {
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

// Report adds a new fill report entry.
func (f *Locations) Report(uid uint, path string, location string) {
	*f = append(
		*f, Location{
			UID:      uid,
			Path:     path,
			Location: location,
		},
	)
}

// Append other locations.
func (f *Locations) Append(other Locations) {
	*f = append(
		*f, other...,
	)
}
