package ifaces

import (
	"reflect"

	"github.com/byte4ever/dsco/fvalue"

	"github.com/byte4ever/dsco/internal/plocation"
)

// FieldValuesGetter defines the ability to get a path/value set (bases).
type FieldValuesGetter interface {
	GetFieldValuesFrom(model ModelInterface) (fvalue.Values, error)
}

type Getter interface {
	Get(
		path string,
		_type reflect.Type,
	) (
		fieldValue *fvalue.Value,
		err error,
	)
}

type ModelInterface interface {
	TypeName() string
	ApplyOn(g Getter) (fvalue.Values, error)
	GetFieldValuesFor(id string, v reflect.Value) fvalue.Values
	Fill(
		inputModelValue reflect.Value,
		layers []fvalue.Values,
	) (plocation.Locations, error)
}
