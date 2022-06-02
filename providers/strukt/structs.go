package strukt

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/utils"
)

const (
	ID           = dsco.Origin("struct")
	IDYamlBuffer = dsco.Origin("yaml buffer")
)

type Entry struct {
	Type  reflect.Type
	Value reflect.Value
}

type Binder struct {
	values map[string]*Entry
	id     dsco.Origin
}

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrTypeMismatch    = errors.New("type mismatch")
)

func (ks *Binder) GetPostProcessErrors() []error {
	return nil
}

func (ks *Binder) Bind(
	key string,
	set bool,
	dstType reflect.Type,
	dstValue *reflect.Value,
) (
	origin dsco.Origin,
	keyOut string,
	succeed bool,
	err error,
) {
	origin = ks.id
	keyOut = key

	e, found := ks.values[key]

	if !found {
		return
	}

	if e.Type.Kind() != dstType.Kind() || e.Type.Elem().Kind() != dstType.Elem().Kind() {
		err = fmt.Errorf(
			"cannot bind type %v to type %v: %w",
			e.Type,
			dstType,
			ErrTypeMismatch,
		)

		return
	}

	if set {
		*dstValue = e.Value
		succeed = true
	}

	return
}

// Provide creates nre env key searcher.
func provide(i interface{}, id dsco.Origin) (*Binder, error) {
	keys := make(map[string]*Entry)
	res := &Binder{values: keys}
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	res.id = id

	err := res.scanStructure("", t.Elem(), v.Elem())
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Provide creates nre env key searcher.
func Provide(i interface{}) (*Binder, error) {
	return provide(i, ID)
}

type ReadCloseProvider interface {
	ReadClose(perform func(r io.Reader) error) error
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

func (ks *Binder) scanStructure(
	rootKey string,
	t reflect.Type,
	v reflect.Value,
) (err error) {
	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)

		name := strings.Split(
			strings.ReplaceAll(
				fieldType.Tag.Get("yaml"),
				" ",
				"",
			),
			",",
		)[0]

		var s string

		if name != "" {
			s = name
		} else {
			s = utils.ToSnakeCase(fieldType.Name)
		}

		key := appendKey(rootKey, s)

		if (fieldType.Type.Kind() != reflect.Ptr) &&
			(fieldType.Type.Kind() != reflect.Slice) {
			return fmt.Errorf("A %s/%v: %w", key, fieldType.Type.String(), ErrUnsupportedType)
		}

		switch fieldType.Type.String() {
		case
			"*time.Time":
			if !v.Field(i).IsNil() {
				ks.values[key] = &Entry{
					Type:  v.Field(i).Type(),
					Value: v.Field(i),
				}
			}

			continue
		}

		e := fieldType.Type.Elem()
		switch e.Kind() {
		case
			reflect.Array,
			reflect.Chan,
			reflect.Complex128,
			reflect.Complex64,
			reflect.Func,
			reflect.Interface,
			reflect.Invalid,
			reflect.Map,
			reflect.Ptr,
			reflect.Slice,
			reflect.Uintptr,
			reflect.UnsafePointer:
			return fmt.Errorf("B %s/%v: %w", key, fieldType.Type.String(), ErrUnsupportedType)

		case reflect.Struct:
			if err := ks.scanStructure(key, e, v.Field(i).Elem()); err != nil {
				return err
			}

		case
			reflect.Int64,
			reflect.Uint64,
			reflect.Int32,
			reflect.Uint32,
			reflect.Int16,
			reflect.Uint16,
			reflect.Int8,
			reflect.Uint8,
			reflect.Int,
			reflect.Uint,
			reflect.Float32,
			reflect.Float64,
			reflect.Bool,
			reflect.String:
			if !v.Field(i).IsNil() {
				ks.values[key] = &Entry{
					Type:  v.Field(i).Type(),
					Value: v.Field(i),
				}
			}
		}
	}

	return nil
}

func appendKey(a, b string) string {
	if a == "" {
		return b
	}

	return a + "-" + b
}
