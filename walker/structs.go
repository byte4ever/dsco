package walker

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco"
)

type StructBuilder struct {
	inputTypeName string
	base          Base
}

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

// ErrStructTypeDiffer represent an error where cannot produce a valid value
// base because struct source type and struct type to fill differs.
var ErrStructTypeDiffer = errors.New("struct type differ")

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
