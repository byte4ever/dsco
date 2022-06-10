package dsco

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/utils"
)

// ErrNotPointerOnStruct represents an error where the interface to fill
// is not a pointer to struct.
var ErrNotPointerOnStruct = errors.New("not pointer on struct")

// ErrUnsupportedType represents an error where the interface to fill
// contains some unsupported types.
var ErrUnsupportedType = errors.New("unsupported type")

// ErrRecursiveStruct represents an error where the interface to fill
// is recursive. So it cannot be allocated by the filler.
var ErrRecursiveStruct = errors.New("recursive struct")

// ErrRequireEmptyStruct represents an error where the interface to fill
// is not empty (nil pointers for every field).
var ErrRequireEmptyStruct = errors.New("require empty struct")

func checkStruct(model interface{}) error {
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)

	if modelType.Kind() != reflect.Ptr ||
		modelType.Elem().Kind() != reflect.Struct ||
		modelValue.IsNil() {
		return ErrNotPointerOnStruct
	}

	return checkStructRec(
		map[string]string{modelType.String(): ""},
		"",
		modelValue.Elem(),
	)
}

func checkStructRec(
	types map[string]string,
	inputKey string,
	v reflect.Value,
) error {
	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldTyp := v.Type().Field(i)

		key := utils.GetKeyName(inputKey, fieldTyp)

		if fieldTyp.Type.Kind() != reflect.Ptr &&
			fieldTyp.Type.Kind() != reflect.Slice {
			return fmt.Errorf(
				"%s(%s) : %w",
				key,
				fieldTyp.Type.String(),
				ErrUnsupportedType,
			)
		}

		switch fieldTyp.Type.String() {
		case "*time.Time":
			continue
		}

		fieldRefTyp := fieldTyp.Type.Elem()

		if fieldRefTyp.Kind() == reflect.Struct {
			en := fieldTyp.Type.String()

			if pKey, found := types[en]; found {
				return fmt.Errorf(
					"%s cycles with %s for type %s: %w",
					displayRoot(pKey),
					displayRoot(key),
					en,
					ErrRecursiveStruct,
				)
			}

			types[en] = key

			if !fieldVal.IsNil() {
				return fmt.Errorf(
					"key %s is defined: %w",
					key,
					ErrRequireEmptyStruct,
				)
			}

			fv := reflect.New(fieldRefTyp)
			if err := checkStructRec(types, key, fv.Elem()); err != nil {
				return err
			}

			delete(types, en)

			continue
		}

		if !fieldVal.IsNil() {
			return fmt.Errorf(
				"key %s is defined: %w",
				key,
				ErrRequireEmptyStruct,
			)
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
