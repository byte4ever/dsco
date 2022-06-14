package sbased2

import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco"
)

// ErrKeyNotFound represent an error where ....
var ErrKeyNotFound = errors.New("key not found")

// ErrNotUnused represent an error where ....
var ErrNotUnused = errors.New("not unused")

// ErrInvalidType represent an error where ....
var ErrInvalidType = errors.New("invalid type")

// Bind implements the dsco.Binder2 interface.
func (s *Binder) Bind(
	key string,
	dstType reflect.Type,
) dsco.BindingAttempt {
	const errFmt = "%s: %w"

	// check for alias collisions
	if _, found := s.aliases[key]; found {
		return dsco.BindingAttempt{
			Error: fmt.Errorf(errFmt, key, ErrAliasCollision),
		}
	}

	entry, found := s.values[key]
	if !found {
		return dsco.BindingAttempt{}
	}

	var tp reflect.Value

	switch dstType.Kind() { //nolint:exhaustive // it's expected
	case reflect.Pointer:
		tp = reflect.New(dstType.Elem())

		if err := yaml.Unmarshal(
			[]byte(entry.value), tp.Interface(),
		); err != nil {
			return dsco.BindingAttempt{
				Location: entry.location,
				Error: fmt.Errorf(
					errFmt,
					entry.location,
					ErrParse,
				),
			}
		}

		entry.state = unused

		return dsco.BindingAttempt{
			Location: entry.location,
			Value:    tp,
		}
	case reflect.Slice:
		tp = reflect.New(dstType)

		if err := yaml.Unmarshal(
			[]byte(entry.value), tp.Interface(),
		); err != nil {
			return dsco.BindingAttempt{
				Location: entry.location,
				Error: fmt.Errorf(
					errFmt,
					entry.location,
					ErrParse,
				),
			}
		}

		entry.state = unused

		return dsco.BindingAttempt{
			Location: entry.location,
			Value:    tp.Elem(),
		}

	default:
		return dsco.BindingAttempt{
			Location: entry.location,
			Error: fmt.Errorf(
				errFmt,
				entry.location,
				ErrInvalidType,
			),
		}
	}
}
