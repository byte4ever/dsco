package walker

import (
	"fmt"
	"reflect"
)

type ValueNode struct {
	Type        reflect.Type
	VisiblePath string
	UID         uint
}

func (n *ValueNode) Fill(
	fillReporter FillReporter,
	value reflect.Value,
	layers []FieldValues,
) {
	for _, layer := range layers {
		fieldValue := layer[n.UID]

		if fieldValue != nil {
			delete(layer, n.UID)
			value.Set(fieldValue.value)

			fillReporter.ReportUse(
				n.UID,
				n.VisiblePath,
				fieldValue.location,
			)

			return
		}
	}

	fillReporter.ReportUnused(n.VisiblePath)
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
