package strukt

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/utils"
)

const (
	ID           = dsco.Origin("struct")
	IDYamlBuffer = dsco.Origin("yaml buffer")
)

type Entry struct {
	Value reflect.Value
}

type Entries map[string]*Entry

type Binder struct {
	entries Entries
	id      dsco.Origin
}

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrTypeMismatch    = errors.New("type mismatch")
)

func (b *Binder) GetPostProcessErrors() []error {
	return nil
}

func (b *Binder) Bind(
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
	origin = b.id
	keyOut = key

	e, found := b.entries[key]

	if !found {
		return
	}

	et := e.Value.Type()

	if et.Kind() != (dstValue).Type().Kind() || et.Elem().Kind() != (dstValue).Type().Elem().Kind() {
		err = fmt.Errorf(
			"cannot bind type %v to type %v: %w",
			et,
			(dstValue).Type(),
			ErrTypeMismatch,
		)

		return
	}

	if set {
		outVal = e.Value
		succeed = true
	}

	return
}

// Provide creates nre env key searcher.
func provide(i interface{}, id dsco.Origin) (*Binder, error) {
	keys := make(Entries)
	res := &Binder{entries: keys}
	v := reflect.ValueOf(i)
	res.id = id

	err := res.buildEntries("", v.Elem())
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Provide creates nre env key searcher.
func Provide(i interface{}) (*Binder, error) {
	return provide(i, ID)
}

type InterfaceProvider interface {
	GetInterface() (interface{}, error)
}

func ProvideFromInterfaceProvider(ip InterfaceProvider) (*Binder, error) {
	k, err := ip.GetInterface()
	if err != nil {
		return nil, err
	}

	return provide(k, IDYamlBuffer)
}

var structToIntercept = map[string]struct{}{
	"*time.Time": {},
}

func (b *Binder) addEntry(key string, value reflect.Value) {
	if !value.IsNil() {
		b.entries[key] = &Entry{
			Value: value,
		}
	}
}

func (b *Binder) buildEntries(
	rootKey string,
	v reflect.Value,
) (err error) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)

		key := utils.GetKeyName(rootKey, fieldType)

		if (fieldType.Type.Kind() != reflect.Ptr) &&
			(fieldType.Type.Kind() != reflect.Slice) {
			return fmt.Errorf("A %s/%v: %w", key, fieldType.Type.String(), ErrUnsupportedType)
		}

		if _, found := structToIntercept[fieldType.Type.String()]; found {
			b.addEntry(key, v.Field(i))
			continue
		}

		e := fieldType.Type.Elem()

		switch e.Kind() {
		case
			reflect.Array, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Func, reflect.Interface,
			reflect.Invalid, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Uintptr, reflect.UnsafePointer:
			return fmt.Errorf("B %s/%v: %w", key, fieldType.Type.String(), ErrUnsupportedType)

		case reflect.Struct:
			if err := b.buildEntries(key, v.Field(i).Elem()); err != nil {
				return err
			}

			continue

		case
			reflect.Int64, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16,
			reflect.Int8, reflect.Uint8, reflect.Int, reflect.Uint, reflect.Float32, reflect.Float64, reflect.Bool,
			reflect.String:
			b.addEntry(key, v.Field(i))

			continue
		}
	}

	return nil
}
