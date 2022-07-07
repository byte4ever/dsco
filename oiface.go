package dsco

import (
	"reflect"

	"github.com/byte4ever/dsco/internal"
	"github.com/byte4ever/dsco/internal/fvalue"
	"github.com/byte4ever/dsco/internal/plocation"
)

// FieldValuesGetter defines the ability to get a path/value set (bases).
type FieldValuesGetter interface {
	GetFieldValuesFrom(model ModelInterface) (fvalue.Values, error)
}

type ModelInterface interface {
	TypeName() string
	ApplyOn(g internal.Getter) (fvalue.Values, error)
	GetFieldValuesFor(id string, v reflect.Value) fvalue.Values
	Fill(
		inputModelValue reflect.Value,
		layers []fvalue.Values,
	) (plocation.Locations, error)
}
