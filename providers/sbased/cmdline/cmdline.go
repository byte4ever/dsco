package cmdline

import (
	"fmt"
	"regexp"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/providers/sbased"
)

const id = dsco.Origin("cmdline")

var re = regexp.MustCompile(`^--([a-z\d_-]+)=(.+)$`)

// EntriesProvider is an entries' provider that extract entries from
// command line.
type EntriesProvider struct {
	values sbased.Entries
}

// GetEntries implements sbased.EntriesProvider interface.
func (ks *EntriesProvider) GetEntries() sbased.Entries {
	return ks.values
}

// GetOrigin implements sbased.EntriesProvider interface.
func (ks *EntriesProvider) GetOrigin() dsco.Origin {
	return id
}

// NewEntriesProvider creates an entries' provider that parses and extract
// parameters from command line.
//
// 		ep, err := NewEntriesProvider(os.Args[1:])
//
// Each command line parameter MUST match  regexp.
func NewEntriesProvider(commandLine []string) (*EntriesProvider, error) {
	lo := len(commandLine)

	if lo == 0 {
		return &EntriesProvider{}, nil
	}

	keys := make(sbased.Entries, lo)

	for idx, arg := range commandLine {
		m := re.FindStringSubmatch(arg)

		if 3 != len(m) { //nolint:gomnd // ok
			return nil, fmt.Errorf("arg #%d - (%s): %w", idx, arg, ErrFormatParam)
		}

		keys[m[1]] = &sbased.Entry{
			ExternalKey: fmt.Sprintf("--%s", m[1]),
			Value:       m[2],
		}
	}

	return &EntriesProvider{values: keys}, nil
}
