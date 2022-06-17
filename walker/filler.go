package walker

import (
	"fmt"
	"reflect"
)

type fillerValue struct {
	path  string
	value *reflect.Value
}

type fillerValues []*fillerValue

func Fill(inputModel any) error {
	var values fillerValues

	w := walker{
		walkFunc: func(path string, value *reflect.Value) error {
			values = append(
				values, &fillerValue{
					path:  path,
					value: value,
				},
			)
			return nil
		},
	}

	if err := w.walk("", reflect.ValueOf(inputModel)); err != nil {
		return err
	}

	for _, value := range values {
		fmt.Println(value.path, value.value.Type().String())
	}

	return nil
}
