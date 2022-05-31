package dsco

import (
	"reflect"
)

type Binder interface {
	Bind(
		key string,
		set bool,
		dstType reflect.Type,
		dstValue *reflect.Value,
	) (
		origin Origin,
		keyOut string,
		succeed bool,
		err error,
	)
	GetPostProcessErrors() []error
}
