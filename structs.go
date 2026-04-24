package dsco

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/internal/fvalue"
	"github.com/byte4ever/dsco/registry"
)

// ErrStructTypeDiffer represent an error where cannot produce a valid value
// bases because struct source type and struct type to fillHelper differs.
var ErrStructTypeDiffer = errors.New("struct type differ")

// StructBuilder is a structure layer builder.
type StructBuilder struct {
	value reflect.Value
	id    string
}

func (s *StructBuilder) GetFieldValuesFrom(model ModelInterface) (
	fvalue.Values, error,
) {
	modelTName := model.TypeName() //nolint:ifshort // buggy linter

	if ltn := registry.LongTypeName(s.value.Type()); modelTName != ltn {
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

// ReportInventory implements InventoryReporter by enumerating every
// non-nil field of the source struct and recording its value as a
// FieldProvision. No I/O is performed.
func (s *StructBuilder) ReportInventory(
	mdl ModelInterface,
) (LayerInventory, error) { //nolint:unparam // error required by InventoryReporter interface
	values := mdl.GetFieldValuesFor(s.id, s.value)

	provides := make([]FieldProvision, 0, len(values))
	for _, fv := range values {
		// Dereference pointer values so the report shows the user-visible
		// scalar, not a *T pointer.
		val := fv.Value
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}

		provides = append(provides, FieldProvision{
			FieldUID: fv.Path,
			Value:    val.Interface(),
		})
	}

	return LayerInventory{
		Name:     "struct:" + s.id,
		Provides: provides,
	}, nil
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
