package cmdline

import (
	"fmt"
	"regexp"

	"github.com/byte4ever/dsco/sbased2"
)

const locationFmt = "cmdline[--%s]"

var re = regexp.MustCompile(
	`^--([a-z][a-z\d]*(?:[-_][a-z][a-z\d]*)*)=(.+)$`,
)

// EntriesProvider is an entries' provider that extract entries from
// command line.
type EntriesProvider struct {
	stringValues sbased2.StringValues
}

// GetStringValues implements sbased2.StringValuesProvider interface.
func (ep *EntriesProvider) GetStringValues() sbased2.StringValues {
	return ep.stringValues
}

// NewEntriesProvider creates an entries' provider that parses and extract
// parameters from command line.
//
// Each command line parameter MUST match regexp '^--([a-z\d_-]+)=(.+)$'.
// ErrParamFormat is returned in such a case.
//
func NewEntriesProvider(commandLine []string) (*EntriesProvider, error) {
	lo := len(commandLine)

	if lo == 0 {
		return &EntriesProvider{}, nil
	}

	dedup := make(map[string]int)

	keys := make(sbased2.StringValues, 0, lo)

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
					"arg %s: %w",
					arg,
					ErrParamFormat,
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
			keys = append(
				keys, &sbased2.StringValue{
					Key:      key,
					Location: fmt.Sprintf(locationFmt, key),
					Value:    groups[2],
				},
			)
		}
	}

	if len(errs) < 1 {
		return &EntriesProvider{stringValues: keys}, nil
	}

	return nil, &ParamError{
		Positions: positions,
		Errs:      errs,
	}
}
