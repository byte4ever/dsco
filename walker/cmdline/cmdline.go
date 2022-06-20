package cmdline

import (
	"fmt"
	"regexp"

	"github.com/byte4ever/dsco/walker/svalues"
)

const locationFmt = "cmdline[--%s]"

var re = regexp.MustCompile(
	`^--([a-z][a-z\d]*(?:[-_][a-z][a-z\d]*)*)=(.+)$`,
)

// EntriesProvider is an entries' provider that extract entries from
// command line.
type EntriesProvider struct {
	stringValues svalues.StringValues
}

// GetStringValues implements svalues.StringValuesProvider interface.
func (ep *EntriesProvider) GetStringValues() svalues.StringValues {
	return ep.stringValues
}

// NewEntriesProvider creates an entries' provider that parses and extract
// parameters from command line.
//
// Each command line parameter MUST match regexp '^--([a-z\d_-]+)=(.+)$'.
// ErrInvalidFormat is returned in such a case.
//
func NewEntriesProvider(commandLine []string) (*EntriesProvider, error) {
	lo := len(commandLine)

	if lo == 0 {
		return &EntriesProvider{}, nil
	}

	dedup := make(map[string]int, lo)

	stringValues := make(svalues.StringValues, lo)

	expectedGroups := re.NumSubexp() + 1

	var (
		errs      []error
		positions []int
	)

	for idx, arg := range commandLine {
		groups := re.FindStringSubmatch(arg)

		if len(groups) != expectedGroups {
			errs = append(
				errs, fmt.Errorf(
					"arg %q: %w",
					arg,
					ErrInvalidFormat,
				),
			)

			positions = append(positions, idx+1)

			continue
		}

		key := groups[1]
		if prevPosition, found := dedup[key]; found {
			errs = append(
				errs, fmt.Errorf(
					"--%s previous found at position #%d: %w",
					key,
					prevPosition,
					ErrDuplicateParam,
				),
			)

			positions = append(positions, idx+1)

			continue
		}

		dedup[key] = idx + 1

		if len(errs) < 1 {
			stringValues[key] = &svalues.StringValue{
				Location: fmt.Sprintf(locationFmt, key),
				Value:    groups[2],
			}
		}
	}

	if len(errs) < 1 {
		return &EntriesProvider{stringValues: stringValues}, nil
	}

	return nil, &ParamError{
		Positions: positions,
		Errs:      errs,
	}
}
