package sbased2

import (
	"fmt"
)

// Option is processing option for string based binder.
type Option interface {
	apply(opts *internalOpts) error
}

type internalOpts struct {
	aliases map[string]string
}

// AliasesOption defines keys aliasing.
type AliasesOption map[string]string

func (o *internalOpts) applyOptions(options []Option) error {
	for i, option := range options {
		if err := option.apply(o); err != nil {
			return fmt.Errorf(
				"when processing option #%d: %w",
				i,
				err,
			)
		}
	}

	return nil
}

func (a AliasesOption) apply(opts *internalOpts) error {
	opts.aliases = a
	return nil
}

// WithAliases returns a keys aliasing option.
func WithAliases(mapping map[string]string) AliasesOption {
	if lm := len(mapping); lm == 0 {
		return nil
	}

	return mapping
}
