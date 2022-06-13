package sbased2

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

type Value struct {
	Succeed  bool
	Location string
	Value    reflect.Value
	Error    error
}

// Bind implements the dscoBinder interface.
//nolint:revive // refacto soon
func (s *Binder) Bind(
	key string,
	dstType reflect.Type,
) Value { // origin dsco.Origin,
	// keyOut string,
	// succeed bool,
	// outVal reflect.Value,
	// err error,

	const errFmt = "%s/%s: %w"

	if _, found := s.aliases[key]; found {
		return Value{
			Succeed:  false,
			Location: key,
			Value:    reflect.Value{},
			Error:    nil,
		}
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
