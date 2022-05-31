package cmdline

import (
	"fmt"
	"regexp"

	"github.com/byte4ever/goconf"
	"github.com/byte4ever/goconf/providers/sbased"
)

const ID = goconf.Origin("cmdline")

var re = regexp.MustCompile(`^--([a-z\d_-]+)=(.+)$`)

// Provider is dummy.
type Provider struct {
	values sbased.StrEntries
}

func (ks *Provider) GetEntries() sbased.StrEntries {
	return ks.values
}

func (ks *Provider) GetOrigin() goconf.Origin {
	return ID
}

func Provide(optionsLine []string) (*Provider, error) {
	lo := len(optionsLine)

	if lo == 0 {
		return &Provider{}, nil
	}

	keys := make(sbased.StrEntries, lo)

	for idx, arg := range optionsLine {
		m := re.FindStringSubmatch(arg)

		if 3 != len(m) {
			return nil, fmt.Errorf("arg #%d - (%s): %w", idx, arg, ErrFormatParam)
		}

		keys[m[1]] = &sbased.StrEntry{
			ExternalKey: fmt.Sprintf("--%s", m[1]),
			Value:       m[2],
		}
	}

	return &Provider{values: keys}, nil
}
