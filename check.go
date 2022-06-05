package dsco

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/utils"
)

var (
	ErrInvalidDestination = errors.New("requires pointer on struct")
	ErrUnsupportedType    = errors.New("unsupported type")
	ErrRecursiveStruct    = errors.New("recursive struct")
	ErrRequireEmptyStruct = errors.New("require empty struct")
)

func checkStruct(i interface{}) error {
	iType := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	if iType.Kind() != reflect.Ptr || iType.Elem().Kind() != reflect.Struct || v.IsNil() {
		return ErrInvalidDestination
	}

	return checkStructRec(
		map[string]string{iType.String(): ""},
		"",
		v.Elem(),
	)
}

func checkStructRec(
	types map[string]string,
	inputKey string,
	v reflect.Value,
) error {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := v.Type().Field(i)

		key := utils.GetKeyName(inputKey, ft)

		if ft.Type.Kind() != reflect.Ptr && ft.Type.Kind() != reflect.Slice {
			return fmt.Errorf("%s(%s) : %w", key, ft.Type.String(), ErrUnsupportedType)
		}

		switch ft.Type.String() {
		case
			"*time.Time":
			continue
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
