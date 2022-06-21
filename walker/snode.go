package walker

import (
	"reflect"
)

type StructNode struct {
	Type  reflect.Type
	Index IndexedSubNodes
}

func (n *StructNode) Fill(
	fillReporter FillReporter,
	value reflect.Value,
	layers []FieldValues,
) {
	v := reflect.New(n.Type.Elem())

	value.Set(v)

	for _, index := range n.Index {
		index.Node.Fill(
			fillReporter,
			value.Elem().FieldByIndex(index.Index),
			layers,
		)
	}
}

func (n *StructNode) FeedFieldValues(
	srcID string,
	fieldValues FieldValues,
	value reflect.Value,
) {
	if value.IsNil() {
		return
	}

	for _, index := range n.Index {
		index.Node.FeedFieldValues(
			srcID, fieldValues,
			value.Elem().FieldByIndex(index.Index),
		)
	}
}

type IndexedSubNodes []*IndexedSubNode

type IndexedSubNode struct {
	Node  Node
	Index []int
}

func (i IndexedSubNodes) GetIndexes() [][]int {
	var ri [][]int

	for _, node := range i {
		ri = append(ri, node.Index)
	}

	return ri
}

func (n *StructNode) PushSubNodes(index []int, scanned Node) {
	n.Index = append(
		n.Index,
		&IndexedSubNode{
			Index: index, Node: scanned,
		},
	)
}

func (n *StructNode) BuildGetList(s *GetList) {
	for _, index := range n.Index {
		index.Node.BuildGetList(s)
	}
}
