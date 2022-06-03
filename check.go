package dsco

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/byte4ever/dsco/utils"
)

var (
	ErrUnsupportedType    = errors.New("unsupported type")
	ErrRecursiveStruct    = errors.New("recursive struct")
	ErrRequireEmptyStruct = errors.New("require empty struct")
)

func checkStruct(i interface{}) error {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	return checkStructRec(
		map[string]string{t.String(): ""},
		"",
		v.Elem(),
	)
}

func checkStructRec(
	types map[string]string,
	inputKey string,
	v reflect.Value,
) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)

		s := strings.Split(strings.Replace(ft.Tag.Get("yaml"), " ", "", -1), ",")[0]
		if s == "" {
			s = utils.ToSnakeCase(ft.Name)
		}

		key := appendKey(inputKey, s)

		if ft.Type.Kind() != reflect.Ptr && ft.Type.Kind() != reflect.Slice {
			return fmt.Errorf("%s(%s) : %w", key, ft.Type.String(), ErrUnsupportedType)
		}

		e := ft.Type.Elem()

		if e.Kind() == reflect.Struct {
			en := ft.Type.String()

			if pKey, found := types[en]; found {
				return fmt.Errorf(
					"%s cycles with %s for type %s: %w", displayRoot(pKey), displayRoot(key), en, ErrRecursiveStruct,
				)
			}

			types[en] = key

			if f.IsNil() {
				fv := reflect.New(e)
				if err := checkStructRec(types, key, fv.Elem()); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("key %s is defined: %w", key, ErrRequireEmptyStruct)
			}

			delete(types, en)

			continue
		}

		if !f.IsNil() {
			return fmt.Errorf("key %s is defined: %w", key, ErrRequireEmptyStruct)
		}
	}

	return nil
}

func displayRoot(key string) string {
	if key == "" {
		return "main struct"
	}

	return fmt.Sprintf("key %s", key)
}
