package sbased

import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/byte4ever/goconf"
)

type Binder struct {
	internalOpts
	entries  entries
	provider StrEntriesProvider
}

var (
	ErrParse          = errors.New("parse error")
	ErrAliasCollision = errors.New("alias collision")
)

func (s *Binder) Bind(
	key string,
	set bool,
	dstType reflect.Type,
	dstValue *reflect.Value,
) (
	origin goconf.Origin,
	keyOut string,
	succeed bool,
	err error,
) {
	origin = s.provider.GetOrigin()

	if _, found := s.aliases[key]; found {
		err = fmt.Errorf("%s/%s: %w", origin, key, ErrAliasCollision)
		return
	}

	entry, found := s.entries[key]
	if !found {
		return
	}

	entry.bounded = true
	keyOut = entry.ExternalKey

	var tp reflect.Value

	switch dstType.Kind() { //nolint:exhaustive // it's expected
	case reflect.Pointer:
		tp = reflect.New(dstType.Elem())

		if err = yaml.Unmarshal([]byte(entry.Value), tp.Interface()); err != nil {
			err = fmt.Errorf("%s/%s: %w", origin, entry.ExternalKey, ErrParse)
			return
		}

		if set {
			entry.used = true
			succeed = true
			*dstValue = tp
		}

		return

	case reflect.Slice:
		tp = reflect.New(dstType)

		if err = yaml.Unmarshal([]byte(entry.Value), tp.Interface()); err != nil {
			err = fmt.Errorf("%s/%s: %w", origin, entry.ExternalKey, ErrParse)
			return
		}

		if set {
			entry.used = true
			succeed = true
			*dstValue = tp.Elem()
		}

		return

	default:
		panic("should never happen")
	}
}

func Provide(p StrEntriesProvider, options ...Option) (*Binder, error) {
	o := internalOpts{}

	if err := o.applyOptions(options); err != nil {
		return nil, err
	}

	strEntries := p.GetEntries()

	var es entries

	if le := len(strEntries); le > 0 {
		es = make(entries, le)

		for k, v := range strEntries {
			actualKey, found := o.aliases[k]
			if !found {
				actualKey = k
			}

			es[actualKey] = &entry{
				StrEntry: *v,
				bounded:  false,
				used:     false,
			}
		}
	}

	return &Binder{
		internalOpts: o,
		entries:      es,
		provider:     p,
	}, nil
}

var (
	ErrUnboundKey    = errors.New("unbound key")
	ErrOverriddenKey = errors.New("overridden key")
)

func (s *Binder) GetPostProcessErrors() (errs []error) {
	o := s.provider.GetOrigin()

	for _, e := range s.entries {
		if !e.bounded {
			errs = append(errs, fmt.Errorf("%s/%s: %w", o, e.ExternalKey, ErrUnboundKey))
			continue
		}

		if !e.used {
			errs = append(errs, fmt.Errorf("%s/%s: %w", o, e.ExternalKey, ErrOverriddenKey))
			continue
		}
	}

	return
}
