package internal

import (
	"reflect"

	"github.com/byte4ever/dsco/internal/fvalue"
)

type ValueGetter interface {
	Get(
		path string,
		_type reflect.Type,
	) (
		fieldValue *fvalue.Value,
		err error,
	)
}

type StructExpander interface {
	ExpandStruct(
		path string,
		_type reflect.Type,
	) (err error)
}
