package env

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/providers/sbased"
)

const id = dsco.Origin("env")

// ErrInvalidPrefix represents an error when creating the provider with an
// invalid prefix.
var ErrInvalidPrefix = errors.New("invalid prefix")

// EntriesProvider is an entries' provider that extract entries from
// environment variables.
type EntriesProvider struct {
	entries sbased.Entries
	prefix  string
}

var (
	re       = regexp.MustCompile(`^([A-Z][A-Z\d]*)-([A-Z][A-Z\d]*(?:[-_][A-Z][A-Z\d]*)*)=(.*)$`)
	rePrefix = regexp.MustCompile(`^[A-Z][A-Z\d]*$`)
)

// NewEntriesProvider creates an entries provider based on environment variable scanning.
// It's sensitive to a prefix that *MUST* match this regexp '^[A-Z][A-Z\d]*$'.
func NewEntriesProvider(prefix string) (*EntriesProvider, error) {
	// ensure prefix is uppercase
	if !rePrefix.MatchString(prefix) {
		return nil, fmt.Errorf("%q : %w", prefix, ErrInvalidPrefix)
	}

	res := &EntriesProvider{
		prefix: prefix,
	}
	env := os.Environ()
	r := make(sbased.Entries, len(env))

	for _, s := range env {
		m := re.FindStringSubmatch(s)
		if len(m) == 4 && m[1] == res.prefix {
			r[strings.ToLower(m[2])] = &sbased.Entry{
				ExternalKey: fmt.Sprintf("%s-%s", res.prefix, m[2]),
				Value:       m[3],
			}
		}
	}

	if len(r) > 0 {
		res.entries = r
	}

	return res, nil
}

func (ks *EntriesProvider) GetEntries() sbased.Entries {
	return ks.entries
}

func (ks *EntriesProvider) GetOrigin() dsco.Origin {
	return id
}
