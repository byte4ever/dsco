package model

import (
	"reflect"

	"github.com/byte4ever/dsco/fvalues"
	"github.com/byte4ever/dsco/ifaces"
	"github.com/byte4ever/dsco/plocation"
)

type Node interface {
	BuildGetList(s *GetList)
	FeedFieldValues(
		srcID string,
		fieldValues fvalues.FieldValues,
		value reflect.Value,
	)
	Fill(
		value reflect.Value,
		layers []fvalues.FieldValues,
	) (plocation.PathLocations, error)
}
type GetListInterface interface {
	ApplyOn(g ifaces.Getter) (fvalues.FieldValues, error)
	Push(o GetOp)
	Count() int
}
