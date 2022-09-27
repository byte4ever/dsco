package model

import (
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/internal"
	"github.com/byte4ever/dsco/internal/fvalue"
	"github.com/byte4ever/dsco/internal/plocation"
	"github.com/byte4ever/dsco/internal/utils"
)

type ValueNode struct {
	Type        reflect.Type
	VisiblePath string
	UID         uint
}

func (n *ValueNode) Fill(
	value reflect.Value, layers []fvalue.Values,
) (plocation.Locations, error) {
	for _, layer := range layers {
		fieldValue := layer[n.UID]

		if fieldValue != nil {
			delete(layer, n.UID)
			value.Set(fieldValue.Value)

			var pl plocation.Locations

			pl.Report(n.UID, n.VisiblePath, fieldValue.Location)

			return pl, nil
		}
	}

	return nil, fmt.Errorf(
		"%s-[%s]: %w",
		n.VisiblePath,
		utils.LongTypeName(value.Type()),
		ErrUninitializedKey,
	)
}

func (n *ValueNode) FeedFieldValues(
	srcID string,
	fieldValues fvalue.Values,
	value reflect.Value,
) {
	if value.IsNil() {
		return
	}

	fieldValues[n.UID] = &fvalue.Value{
		Value:    value,
		Location: fmt.Sprintf("struct[%s]:%s", srcID, n.VisiblePath),
	}
}

func (n *ValueNode) BuildGetList(s *GetList) {
	s.Push(
		func(g internal.Getter) (uint, *fvalue.Value, error) {
			fieldValue, err := g.Get(n.VisiblePath, n.Type)

			return n.UID, fieldValue, err //nolint:wrapcheck // don't wan to wrap
		},
	)
}
