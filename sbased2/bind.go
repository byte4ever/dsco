package sbased2

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco"
)

// Bind implements the dscoBinder interface.
func (s *Binder) Bind(
	key string,
	dstType reflect.Type,
	markUsed bool,
) dsco.BoundingAttempt {
	// check for alias collisions
	if _, found := s.aliases[key]; found {
		return dsco.BoundingAttempt{
			Error: ErrAliasCollision,
		}
	}

	entry, found := s.values[key]
	if !found {
		return dsco.BoundingAttempt{}
	}

	entry.bounded = true

	var tp reflect.Value

	switch dstType.Kind() { //nolint:exhaustive // it's expected
	case reflect.Pointer:
		tp = reflect.New(dstType.Elem())
	case reflect.Slice:
		tp = reflect.New(dstType)
	default:
		panic("should never happen")
	}

	if err := yaml.Unmarshal(
		[]byte(entry.value), tp.Interface(),
	); err != nil {
		return dsco.BoundingAttempt{
			Location: entry.location,
			Error: fmt.Errorf(
				"%s: %w",
				entry.location,
				ErrParse,
			),
		}
	}

	entry.used = markUsed

	return dsco.BoundingAttempt{
		Location: entry.location,
		Value:    tp,
	}
}
