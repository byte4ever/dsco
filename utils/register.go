package utils

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

var nameToType sync.Map //nolint:gochecknoglobals // required for registration

//nolint:gochecknoinits // required at loading time
func init() {
	registerDefaultTypes()
}

func registerDefaultTypes() {
	Register(R(0))
	Register(R(int8(0)))
	Register(R(int16(0)))
	Register(R(int32(0)))
	Register(R(int64(0)))

	Register(R(uint(0)))
	Register(R(uint8(0)))
	Register(R(uint16(0)))
	Register(R(uint32(0)))
	Register(R(uint64(0)))

	Register(R(float32(0)))
	Register(R(float64(0)))

	Register(R(true))

	Register(R(""))

	Register(&time.Time{})
	Register(R(time.Duration(0)))
}

// LongTypeName returns long name for a type.
func LongTypeName(_type reflect.Type) string {
	var sb strings.Builder

	tp := _type

	if tp.Kind() == reflect.Ptr {
		sb.WriteRune('*')

		tp = _type.Elem()
	}

	if pkg := tp.PkgPath(); pkg != "" {
		sb.WriteString(pkg)
		sb.WriteRune('/')
	}

	sb.WriteString(tp.String())

	return sb.String()
}

// Register registers the type of the value.
func Register(value any) {
	valueType := reflect.TypeOf(value)

	if valueType.Kind() != reflect.Pointer {
		panic("register requires pointer")
	}

	longName := LongTypeName(valueType)

	if _, dup := nameToType.LoadOrStore(longName, valueType); dup {
		panic(
			fmt.Sprintf(
				"dsco: %q duplicate type registration",
				longName,
			),
		)
	}
}

// TypeIsRegistered returns true when type t is registered.
func TypeIsRegistered(t reflect.Type) bool {
	_, found := nameToType.Load(LongTypeName(t))

	return found
}
