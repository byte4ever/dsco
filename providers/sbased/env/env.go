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

var ErrInvalidPrefix = errors.New("invalid prefix")

type Provider struct {
	entries sbased.StrEntries
	prefix  string
}

var re = regexp.MustCompile(`^([A-Z][A-Z\d]*)-([A-Z][A-Z\d]*(?:[-_][A-Z][A-Z\d]*)*)=(.*)$`)
var rePrefix = regexp.MustCompile(`^[A-Z][A-Z\d]*$`)

func Provide(prefix string) (*Provider, error) {
	// ensure prefix is uppercase
	if !rePrefix.MatchString(prefix) {
		return nil, fmt.Errorf("%q : %w", prefix, ErrInvalidPrefix)
	}

	res := &Provider{
		prefix: prefix,
	}
	env := os.Environ()
	r := make(sbased.StrEntries, len(env))

	for _, s := range env {
		m := re.FindStringSubmatch(s)
		if len(m) == 4 && m[1] == res.prefix {
			r[strings.ToLower(m[2])] = &sbased.StrEntry{
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

func (ks *Provider) GetEntries() sbased.StrEntries {
	return ks.entries
}

const ID = dsco.Origin("env")

func (ks *Provider) GetOrigin() dsco.Origin {
	return ID
}
