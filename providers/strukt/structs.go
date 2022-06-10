package strukt

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/utils"
)

const id = dsco.Origin("struct")

// TODO :- lmartin 6/10/22 -: need to add a global type interceptor

//nolint:gochecknoglobals // need to add a global interceptor
var structToIntercept = map[string]struct{}{"*time.Time": {}}

type entry struct {
	Value reflect.Value
}

type entries map[string]*entry

// Binder is a binder for struct.
type Binder struct {
	entries entries
}

// InterfaceProvider defines the ability to provide an interface.
type InterfaceProvider interface {
	GetInterface() (interface{}, error)
}

// ErrUnsupportedType represents an error when the struct contains a field with
// an unsupported type.
var ErrUnsupportedType = errors.New("unsupported type")

// ErrTypeMismatch represents am error when binding fail because type differ for
// the same key.
var ErrTypeMismatch = errors.New("type mismatch")

// GetPostProcessErrors implements dsco.Binder interface.
func (*Binder) GetPostProcessErrors() []error {
	return nil
}

// Bind implements dsco.Binder interface.
//nolint:nakedret, revive // need refactoring here
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
	origin = id
	keyOut = key

	entry, found := b.entries[key]

	if !found {
		return
	}

	entryValTyp := entry.Value.Type()

	if entryValTyp.Kind() != (dstValue).Type().Kind() ||
		entryValTyp.Elem().Kind() != (dstValue).Type().Elem().Kind() {
		err = fmt.Errorf(
			"cannot bind type %v to type %v: %w",
			entryValTyp,
			(dstValue).Type(),
			ErrTypeMismatch,
		)

		return
	}

	if set {
		outVal = entry.Value
		succeed = true
	}

	return
}

// NewBinder creates nre env key searcher.
func NewBinder(i interface{}) (*Binder, error) {
	keys := make(entries)
	res := &Binder{entries: keys}
	v := reflect.ValueOf(i)

	if err := res.buildEntries("", v.Elem()); err != nil {
		return nil, err
	}

	return res, nil
}

// ProvideFromInterfaceProvider creates a binder using an interface provider.
func ProvideFromInterfaceProvider(ip InterfaceProvider) (*Binder, error) {
	k, err := ip.GetInterface()
	if err != nil {
		return nil, fmt.Errorf("when getting interface: %w", err)
	}

	return NewBinder(k)
}

func (b *Binder) addEntry(key string, value reflect.Value) {
	if !value.IsNil() {
		b.entries[key] = &entry{
			Value: value,
		}
	}
}

//nolint:gocognit // is going to be refactored
func (b *Binder) buildEntries(
	rootKey string,
	value reflect.Value,
) (err error) {
	// TODO :- lmartin 6/10/22 -: use structure checker from dsco at
	//  creation time.
	const errFormat = "%s/%value: %w"

	valueTyp := value.Type()

	for i := 0; i < value.NumField(); i++ {
		fieldType := valueTyp.Field(i)

		key := utils.GetKeyName(rootKey, fieldType)

		if (fieldType.Type.Kind() != reflect.Ptr) &&
			(fieldType.Type.Kind() != reflect.Slice) {
			return fmt.Errorf(
				errFormat,
				key,
				fieldType.Type.String(),
				ErrUnsupportedType,
			)
		}

		if _, found := structToIntercept[fieldType.Type.String()]; found {
			b.addEntry(key, value.Field(i))
			continue
		}

		e := fieldType.Type.Elem()

		switch e.Kind() {
		case
			reflect.Array, reflect.Chan, reflect.Complex128, reflect.Complex64,
			reflect.Func, reflect.Interface, reflect.Invalid, reflect.Map,
			reflect.Ptr, reflect.Slice, reflect.Uintptr, reflect.UnsafePointer:
			return fmt.Errorf(
				errFormat,
				key,
				fieldType.Type.String(),
				ErrUnsupportedType,
			)

		case reflect.Struct:
			if err := b.buildEntries(key, value.Field(i).Elem()); err != nil {
				return err
			}

			continue

		case
			reflect.Int64, reflect.Uint64, reflect.Int32, reflect.Uint32,
			reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8,
			reflect.Int, reflect.Uint, reflect.Float32, reflect.Float64,
			reflect.Bool, reflect.String:
			b.addEntry(key, value.Field(i))

			continue
		}
	}

	return nil
}
