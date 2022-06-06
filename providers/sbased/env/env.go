package env

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/providers/sbased"
)

var _ sbased.EntriesProvider = &EntriesProvider{}

const id = dsco.Origin("env")

// EntriesProvider is an entries' provider that extract entries from
// environment variables.
type EntriesProvider struct {
	entries sbased.Entries
	prefix  string
}

var (
	reSubKey = regexp.MustCompile(`^-[A-Z][A-Z\d]*(?:[-_][A-Z][A-Z\d]*)*$`)
	rePrefix = regexp.MustCompile(`^[A-Z][A-Z\d]*$`)
)

func getRePrefixed(prefix string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s([^=]+)=(.*)$", prefix))
}

// NewEntriesProvider creates an entries provider based on environment variable scanning.
// It's sensitive to a prefix that *MUST* match this regexp '^[A-Z][A-Z\d]*$'.
func NewEntriesProvider(prefix string) (*EntriesProvider, []error) {
	// ensure prefix is uppercase
	if !rePrefix.MatchString(prefix) {
		return nil, []error{fmt.Errorf("%q : %w", prefix, ErrInvalidPrefix)}
	}

	res := &EntriesProvider{
		prefix: prefix,
	}
	env := os.Environ()
	r, errs := extractEntries(env, prefix)

	if len(errs) > 0 {
		return nil, errs
	}

	if len(r) > 0 {
		res.entries = r
	}

	return res, nil
}

func extractEntries(env []string, prefix string) (sbased.Entries, []error) {
	var errs []error

	r := make(sbased.Entries, len(env))

	sort.Strings(env)

	rePrefixed := getRePrefixed(prefix)
	for _, s := range env {
		m := rePrefixed.FindStringSubmatch(s)

		if len(m) == 3 { //nolint:gomnd // this is expected
			if reSubKey.MatchString(m[1]) {
				r[strings.ToLower(m[1][1:])] = &sbased.Entry{
					ExternalKey: fmt.Sprintf("%s%s", prefix, m[1]),
					Value:       m[2],
				}
			} else {
				errs = append(errs, fmt.Errorf("env var %s%s: %w", prefix, m[1], ErrInvalidKeyFormat))
			}
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return r, nil
}

// GetEntries implements sbased.EntriesProvider interface.
func (ks *EntriesProvider) GetEntries() sbased.Entries {
	return ks.entries
}

// GetOrigin implements sbased.EntriesProvider interface.
func (ks *EntriesProvider) GetOrigin() dsco.Origin {
	return id
}
