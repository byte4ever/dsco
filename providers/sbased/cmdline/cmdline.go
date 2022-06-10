package cmdline

import (
	"fmt"
	"regexp"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/providers/sbased"
)

// EntriesProvider is an entries' provider that extract entries from
// command line.
type EntriesProvider struct {
	values sbased.Entries
}

const id = dsco.Origin("cmdline")

var re = regexp.MustCompile(
	`^--([a-z][a-z\d]*(?:[-_][a-z][a-z\d]*)*)=(.+)$`,
)

// GetEntries implements sbased.EntriesProvider interface.
func (ep *EntriesProvider) GetEntries() sbased.Entries {
	return ep.values
}

// GetOrigin implements sbased.EntriesProvider interface.
func (*EntriesProvider) GetOrigin() dsco.Origin {
	return id
}

// NewEntriesProvider creates an entries' provider that parses and extract
// parameters from command line.
//
// Each command line parameter MUST match regexp '^--([a-z\d_-]+)=(.+)$'.
// ErrParamFormat is returned in such a case.
//
// Creation will fail if some duplicated options are present.
// ErrDuplicateParam is returned in such case.
//
func NewEntriesProvider(commandLine []string) (*EntriesProvider, error) {
	lo := len(commandLine)

	if lo == 0 {
		return &EntriesProvider{}, nil
	}

	keys := make(sbased.Entries, lo)

	expectedGroups := re.NumSubexp() + 1
	for idx, arg := range commandLine {
		m := re.FindStringSubmatch(arg)

		if len(m) != expectedGroups {
			return nil, fmt.Errorf(
				"arg #%d - (%s): %w",
				idx,
				arg,
				ErrParamFormat,
			)
		}

		_, found := keys[m[1]]
		if found {
			return nil, fmt.Errorf("--%s: %w", m[1], ErrDuplicateParam)
		}

		keys[m[1]] = &sbased.Entry{
			ExternalKey: fmt.Sprintf("--%s", m[1]),
			Value:       m[2],
		}
	}

	return &EntriesProvider{values: keys}, nil
}
