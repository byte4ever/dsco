package walker

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco"
)

// ErrStructTypeDiffer represent an error where cannot produce a valid value
// bases because struct source type and struct type to fillHelper differs.
var ErrStructTypeDiffer = errors.New("struct type differ")

// StructBuilder is a structure layer builder.
type StructBuilder struct {
	value reflect.Value
	id    string
}

func (s *StructBuilder) GetFieldValues(model *Model) (FieldValues, []error) {
	ltn := dsco.LongTypeName(s.value.Type())
	if model.typeName != ltn {
		return nil, []error{
			fmt.Errorf(
				"%s != %s: %w",
				model.typeName,
				ltn,
				ErrStructTypeDiffer,
			),
		}
	}

	return model.FeedFieldValues(
		s.id,
		s.value,
	), nil
}

// NewStructBuilder creates a new structure layer builder.
func NewStructBuilder(inputStruct any, id string) (*StructBuilder, error) {
	return &StructBuilder{
		value: reflect.ValueOf(inputStruct),
		id:    id,
	}, nil
}
