package strukt

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/utils"
)

const (
	ID           = dsco.Origin("struct")
	IDYamlFile   = dsco.Origin("yaml file")
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

// ErrConfNotFound is returned when no file is found.
var ErrConfNotFound = errors.New("no configuration file found")

func ProvideFromReader(reader io.Reader, i interface{}) (*Binder, error) {
	k := reflect.New(reflect.TypeOf(i).Elem()).Interface()
	dec := yaml.NewDecoder(
		reader,
	)

	if err := dec.Decode(k); err != nil {
		return nil, fmt.Errorf("while parsing yaml buffer: %w", err)
	}

	return provide(k, IDYamlBuffer)
}

func ProvideFromFile(
	searchPaths []string,
	fileName string,
	i interface{},
) (
	*Binder,
	error,
) {
	input, err := tryToOpen(searchPaths, fileName)
	if err == nil {
		defer func(input *os.File) {
			errClose := input.Close()
			if errClose != nil {
				panic(errClose)
			}
		}(input)

		p, err := ProvideFromReader(input, i)
		if err != nil {
			return p, err
		}

		p.id = IDYamlFile

		return p, nil
	}

	// else define an empty key searcher
	k := reflect.New(reflect.TypeOf(i).Elem()).Interface()

	return provide(k, IDYamlFile)
}

func tryToOpen(paths []string, name string) (*os.File, error) {
	for _, path := range paths {
		fp := filepath.Join(path, name)
		input, err := os.Open(fp)

		if err == nil {
			return input, nil
		}
	}

	return nil, ErrConfNotFound
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

		if v.Field(i).IsNil() {
			continue
		}

		switch fieldType.Type.String() {
		case
			// todo :- lmartin 5/31/22 -: v- should simplify the switch.
			"*hash.Hash",
			"*time.Duration",
			"*time.Time":
			ks.values[key] = &Entry{
				Type:  v.Field(i).Type(),
				Value: v.Field(i),
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
			return fmt.Errorf("B %s/%v: %w", key, e.Name(), ErrUnsupportedType)

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
			ks.values[key] = &Entry{
				Type:  v.Field(i).Type(),
				Value: v.Field(i),
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
