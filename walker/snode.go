package walker

import (
	"errors"
	"reflect"

	"github.com/byte4ever/dsco/merror"
)

type StructNode struct {
	Type  reflect.Type
	Index IndexedSubNodes
}

type StructNodeError struct {
	merror.MError
}

var ErrStructNode = errors.New("")

func (e StructNodeError) Is(err error) bool {
	return errors.Is(err, ErrStructNode)
}

func (n StructNode) Fill(
	value reflect.Value, layers []FieldValues,
) (PathLocations, error) {
	var pl PathLocations
	var errs StructNodeError

	v := reflect.New(n.Type.Elem())

	value.Set(v)

	for _, index := range n.Index {
		pln, err := index.Node.Fill(
			value.Elem().FieldByIndex(index.Index),
			layers,
		)
		pl.ReportOther(pln)
		if err != nil {
			errs.Add(err)
		}
	}

	if errs.None() {
		return pl, nil
	}

	return pl, errs
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
