package sbased

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco"
)

// Binder is a string value binder.
type Binder struct {
	internalOpts
	entries  entries
	provider EntriesProvider
}

var _ dsco.Binder = &Binder{}

// Bind implements the dscoBinder interface.
//nolint:revive // refacto soon
func (s *Binder) Bind(
	key string,
	set bool,
	dstValue reflect.Value,
) (
	origin dsco.Origin,
	keyOut string,
	succeed bool,
	outVal reflect.Value,
	err error,
) {
	const errFmt = "%s/%s: %w"

	origin = s.provider.GetOrigin()

	if _, found := s.aliases[key]; found {
		err = fmt.Errorf(errFmt, origin, key, ErrAliasCollision)
		return
	}

	entry, found := s.entries[key]
	if !found {
		return
	}

	entry.bounded = true
	keyOut = entry.ExternalKey

	var tp reflect.Value

	dType := (dstValue).Type()
	switch dType.Kind() { //nolint:exhaustive // it's expected
	case reflect.Pointer:
		tp = reflect.New(dType.Elem())

		if err = yaml.Unmarshal(
			[]byte(entry.Value), tp.Interface(),
		); err != nil {
			err = fmt.Errorf(errFmt, origin, entry.ExternalKey, ErrParse)
			return
		}

		if set {
			entry.used = true
			succeed = true
			outVal = tp
		}

		return

	case reflect.Slice:
		tp = reflect.New(dType)

		if err = yaml.Unmarshal(
			[]byte(entry.Value), tp.Interface(),
		); err != nil {
			err = fmt.Errorf(errFmt, origin, entry.ExternalKey, ErrParse)
			return
		}

		if set {
			entry.used = true
			succeed = true
			outVal = tp.Elem()
		}

		return

	default:
		panic("should never happen")
	}
}

// NewBinder creates a new string based binder.
func NewBinder(
	provider EntriesProvider,
	options ...Option,
) (
	*Binder,
	error,
) {
	internalOptions := internalOpts{}

	if err := internalOptions.applyOptions(options); err != nil {
		return nil, err
	}

	strEntries := provider.GetEntries()

	var es entries

	if le := len(strEntries); le > 0 {
		es = make(entries, le)

		for key, strEntry := range strEntries {
			actualKey, found := internalOptions.aliases[key]
			if !found {
				actualKey = key
			}

			es[actualKey] = &entry{
				Entry:   *strEntry,
				bounded: false,
				used:    false,
			}
		}
	}

	return &Binder{
		internalOpts: internalOptions,
		entries:      es,
		provider:     provider,
	}, nil
}

// GetPostProcessErrors returns all errors encountered during processing of the
// layer.
func (s *Binder) GetPostProcessErrors() []error {
	var errs []error

	const errFormat = "%s/%s: %w"

	origin := s.provider.GetOrigin()

	for _, entry := range s.entries {
		if !entry.bounded {
			errs = append(
				errs, fmt.Errorf(
					errFormat,
					origin,
					entry.ExternalKey,
					ErrUnboundKey,
				),
			)

			continue
		}

		if !entry.used {
			errs = append(
				errs,
				fmt.Errorf(
					errFormat,
					origin,
					entry.ExternalKey,
					ErrOverriddenKey,
				),
			)

			continue
		}
	}

	return errs
}
