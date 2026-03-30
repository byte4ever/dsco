package internal

import (
	"reflect"

	"github.com/byte4ever/dsco/internal/fvalue"
)

// ValueGetter defines the ability to get field values.
type ValueGetter interface {
	Get(path string, fieldType reflect.Type) (*fvalue.Value, error)
}

// StructExpander defines the ability to expand struct definitions.
type StructExpander interface {
	ExpandStruct(path string, structType reflect.Type) error
}

// ModelInterface represents the target configuration structure model.
type ModelInterface interface {
	TypeName() string
}
