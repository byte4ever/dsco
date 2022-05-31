package sbased

type internalOpts struct {
	aliases map[string]string
}

func (o *internalOpts) applyOptions(os []Option) (err error) {
	for _, option := range os {
		if err = option.apply(o); err != nil {
			return
		}
	}

	return
}

type Option interface {
	apply(opts *internalOpts) error
}

type AliasesOption map[string]string

func (a AliasesOption) apply(opts *internalOpts) error {
	opts.aliases = a
	return nil
}

func WithAliases(mapping map[string]string) AliasesOption {
	if lm := len(mapping); lm == 0 {
		return nil
	}

	return mapping
}
