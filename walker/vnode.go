package walker

import (
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/walker/plocation"
)

type ValueNode struct {
	Type        reflect.Type
	VisiblePath string
	UID         uint
}

func (n *ValueNode) Fill(
	value reflect.Value, layers []FieldValues,
) (plocation.PathLocations, error) {
	for _, layer := range layers {
		fieldValue := layer[n.UID]

		if fieldValue != nil {
			delete(layer, n.UID)
			value.Set(fieldValue.value)

			var pl plocation.PathLocations
			pl.Report(n.UID, n.VisiblePath, fieldValue.location)
			return pl, nil
		}
	}

	return nil, fmt.Errorf("%w", ErrUninitializedKey)
}

func (n *ValueNode) FeedFieldValues(
	srcID string,
	fieldValues FieldValues,
	value reflect.Value,
) {
	if value.IsNil() {
		return
	}

	fieldValues[n.UID] = &FieldValue{
		value:    value,
		location: fmt.Sprintf("struct[%s]:%s", srcID, n.VisiblePath),
	}
}

func (n *ValueNode) BuildGetList(s *GetList) {
	s.Push(
		func(g Getter) (uint, *FieldValue, error) {
			fieldValue, err := g.Get(n.VisiblePath, n.Type)

			return n.UID, fieldValue, err //nolint:wrapcheck // don't wan to wrap
		},
	)
}
