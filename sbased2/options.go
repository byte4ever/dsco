package sbased2

import (
	"errors"
	"fmt"
)

// ErrNoAliasesProvided represent an error where no aliases map was
// provided with option.
var ErrNoAliasesProvided = errors.New("no aliases provided")

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
	if len(a) > 0 {
		opts.aliases = a
		return nil
	}

	return ErrNoAliasesProvided
}

// WithAliases returns a keys aliasing option.
func WithAliases(mapping map[string]string) AliasesOption {
	return mapping
}
