package env

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/byte4ever/dsco/svalues"
	"github.com/byte4ever/dsco/utils"
)

// EntriesProvider is an entries' provider that extract entries from
// environment variables.
type EntriesProvider struct {
	stringValues svalues.StringValues
}

const (
	reSubKeyExp = `^-[A-Z][A-Z\d]*(?:[-_][A-Z][A-Z\d]*)*$`
	rePrefixExp = `^[A-Z][A-Z\d]*$`
)

var (
	reSubKey = regexp.MustCompile(reSubKeyExp)
	rePrefix = regexp.MustCompile(rePrefixExp)
)

func getRePrefixed(prefix string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s([^=]+)=(.*)$", prefix))
}

func newProvider(
	prefix string,
	environ []string,
) (
	*EntriesProvider,
	error,
) {
	// ensure prefix is uppercase
	if !rePrefix.MatchString(prefix) {
		return nil, fmt.Errorf("%q : %w", prefix, ErrInvalidPrefix)
	}

	stringValues, err := extractStringValues(environ, prefix)

	if err != nil {
		return nil, err
	}

	res := &EntriesProvider{}

	if len(stringValues) > 0 {
		res.stringValues = stringValues
	}

	return res, nil
}

// NewEntriesProvider creates an entries provider based on environment variable
// scanning. It's sensitive to a prefix that *MUST* match this regexp
// '^[A-Z][A-Z\d]*$'.
func NewEntriesProvider(prefix string) (*EntriesProvider, error) {
	return newProvider(prefix, os.Environ())
}

func extractStringValues(env []string, prefix string) (
	svalues.StringValues, error,
) {
	var ambiguousKeys []string

	stringValues := make(svalues.StringValues, len(env))

	sort.Strings(env)

	rePrefixed := getRePrefixed(prefix)
	for _, s := range env {
		groups := rePrefixed.FindStringSubmatch(s)

		if len(groups) == rePrefixed.NumSubexp()+1 {
			if reSubKey.MatchString(groups[1]) {
				stringValues[strings.ToLower(groups[1][1:])] =
					&svalues.StringValue{
						Location: fmt.Sprintf(
							"env[%s%s]",
							prefix,
							groups[1],
						),
						Value: groups[2],
					}

				continue
			}

			ambiguousKeys = append(
				ambiguousKeys, fmt.Sprintf(
					"%s%s",
					prefix,
					groups[1],
				),
			)
		}
	}

	const ambiguousErrFmt = "%s %w"

	if len(ambiguousKeys) > 0 {
		if len(ambiguousKeys) == 1 {
			return nil, fmt.Errorf(
				ambiguousErrFmt,
				utils.FormatStringSequence(ambiguousKeys),
				ErrAmbiguousKey,
			)
		}

		return nil, fmt.Errorf(
			ambiguousErrFmt,
			utils.FormatStringSequence(ambiguousKeys),
			ErrAmbiguousKeys,
		)
	}

	return stringValues, nil
}

// GetStringValues implements sbased2.StringValuesProvider interface.
func (e *EntriesProvider) GetStringValues() svalues.StringValues {
	return e.stringValues
}
