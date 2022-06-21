package walker

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco"
)

// ErrStructTypeDiffer represent an error where cannot produce a valid value
// base because struct source type and struct type to fillHelper differs.
var ErrStructTypeDiffer = errors.New("struct type differ")

// StructBuilder is a structure layer builder.
type StructBuilder struct {
	base          Base
	inputTypeName string
}

// NewStructBuilder creates a new structure layer builder.
func NewStructBuilder(inputStruct any, id string) (*StructBuilder, error) {
	var maxId int

	base := make(Base)

	wlkr := walker{
		fieldAction: func(order int, path string, value *reflect.Value) error {
			base[order] = assignedValue{
				path:     path,
				location: fmt.Sprintf("struct(%s)[%s]", id, path),
				value:    value,
			}
			return nil
		},
		isGetter: true,
	}

	err := wlkr.walkRec(&maxId, "", reflect.ValueOf(inputStruct))
	if err != nil {
		return nil, err
	}

	return &StructBuilder{
		inputTypeName: dsco.LongTypeName(reflect.TypeOf(inputStruct)),
		base:          base,
	}, nil
}

// GetBaseFor implements BaseGetter interface.
func (s *StructBuilder) GetBaseFor(inputModel any) (Base, []error) {
	modelTypeName := dsco.LongTypeName(
		reflect.TypeOf(inputModel),
	)

	if modelTypeName != s.inputTypeName {
		return nil, []error{
			fmt.Errorf(
				"model(%q) != input(%q): %w",
				modelTypeName,
				s.inputTypeName,
				ErrStructTypeDiffer,
			),
		}
	}

	return s.base, nil
}
