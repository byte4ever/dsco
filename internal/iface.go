package internal

import (
	"reflect"

	"github.com/byte4ever/dsco/internal/fvalue"
)

type Getter interface {
	Get(
		path string,
		_type reflect.Type,
	) (
		fieldValue *fvalue.Value,
		err error,
	)
}

type Expander interface {
	Expand(
		path string,
		_type reflect.Type,
	) (err error)
}
