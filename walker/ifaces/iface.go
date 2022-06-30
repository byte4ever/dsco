package ifaces

import (
	"reflect"

	"github.com/byte4ever/dsco/walker/fvalues"
	"github.com/byte4ever/dsco/walker/plocation"
)

// FieldValuesGetter defines the ability to get a path/value set (bases).
type FieldValuesGetter interface {
	GetFieldValuesFrom(model ModelInterface) (fvalues.FieldValues, error)
}

type Getter interface {
	Get(
		path string,
		_type reflect.Type,
	) (
		fieldValue *fvalues.FieldValue,
		err error,
	)
}

type ModelInterface interface {
	TypeName() string
	ApplyOn(g Getter) (fvalues.FieldValues, error)
	GetFieldValuesFor(id string, v reflect.Value) fvalues.FieldValues
	Fill(
		inputModelValue reflect.Value,
		layers []fvalues.FieldValues,
	) (plocation.PathLocations, error)
}
