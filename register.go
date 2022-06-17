package dsco

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

	Register(R(""))

	Register(&time.Time{})
	Register(R(time.Duration(0)))

}

func longTypeName(_type reflect.Type) string {
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

func Register(value any) {
	t := reflect.TypeOf(value)

	if t.Kind() != reflect.Pointer {
		panic("register requires pointer")
	}

	longName := longTypeName(t)

	if _, dup := nameToType.LoadOrStore(longName, t); dup {
		panic(
			fmt.Sprintf(
				"dsco: %q duplicate type registration",
				longName,
			),
		)
	}
}

func TypeIsRegistered(t reflect.Type) bool {
	if t.Kind() != reflect.Pointer {
		panic("dsco: check if type is registered requires pointers")
	}

	_, found := nameToType.Load(longTypeName(t))

	return found
}
