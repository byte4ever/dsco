package walker

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/walker/fvalues"
	"github.com/byte4ever/dsco/walker/ifaces"
)

// ErrStructTypeDiffer represent an error where cannot produce a valid value
// bases because struct source type and struct type to fillHelper differs.
var ErrStructTypeDiffer = errors.New("struct type differ")

// StructBuilder is a structure layer builder.
type StructBuilder struct {
	value reflect.Value
	id    string
}

func (s *StructBuilder) GetFieldValuesFrom(model ifaces.ModelInterface) (
	fvalues.FieldValues, error,
) {
	modelTName := model.TypeName()
	if ltn := dsco.LongTypeName(s.value.Type()); modelTName != ltn {
		return nil,
			fmt.Errorf(
				"%s != %s: %w",
				modelTName,
				ltn,
				ErrStructTypeDiffer,
			)

	}

	return model.GetFieldValuesFor(
		s.id,
		s.value,
	), nil
}

// NewStructBuilder creates a new structure layer builder.
func NewStructBuilder(inputStruct any, id string) (*StructBuilder, error) {
	if inputStruct == nil {
		return nil, ErrNilInput
	}

	v := reflect.ValueOf(inputStruct)
	vt := v.Type()

	if vt.Kind() != reflect.Pointer ||
		vt.Elem().Kind() != reflect.Struct ||
		v.IsNil() {
		return nil, InvalidInputError{
			Type: vt,
		}
	}

	return &StructBuilder{
		value: reflect.ValueOf(inputStruct),
		id:    id,
	}, nil
}
