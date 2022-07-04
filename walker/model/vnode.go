package model

import (
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/walker/fvalues"
	"github.com/byte4ever/dsco/walker/ifaces"
	"github.com/byte4ever/dsco/walker/plocation"
)

type ValueNode struct {
	Type        reflect.Type
	VisiblePath string
	UID         uint
}

func (n *ValueNode) Fill(
	value reflect.Value, layers []fvalues.FieldValues,
) (plocation.PathLocations, error) {
	for _, layer := range layers {
		fieldValue := layer[n.UID]

		if fieldValue != nil {
			delete(layer, n.UID)
			value.Set(fieldValue.Value)

			var pl plocation.PathLocations

			pl.Report(n.UID, n.VisiblePath, fieldValue.Location)

			return pl, nil
		}
	}

	return nil, fmt.Errorf(
		"%s: %w",
		n.VisiblePath,
		ErrUninitializedKey,
	)
}

func (n *ValueNode) FeedFieldValues(
	srcID string,
	fieldValues fvalues.FieldValues,
	value reflect.Value,
) {
	if value.IsNil() {
		return
	}

	fieldValues[n.UID] = &fvalues.FieldValue{
		Value:    value,
		Location: fmt.Sprintf("struct[%s]:%s", srcID, n.VisiblePath),
	}
}

func (n *ValueNode) BuildGetList(s *GetList) {
	s.Push(
		func(g ifaces.Getter) (uint, *fvalues.FieldValue, error) {
			fieldValue, err := g.Get(n.VisiblePath, n.Type)

			return n.UID, fieldValue, err //nolint:wrapcheck // don't wan to wrap
		},
	)
}
