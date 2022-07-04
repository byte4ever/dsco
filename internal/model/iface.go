package model

import (
	"reflect"

	"github.com/byte4ever/dsco/fvalue"
	"github.com/byte4ever/dsco/ifaces"
	"github.com/byte4ever/dsco/internal/plocation"
)

type Node interface {
	BuildGetList(s *GetList)
	FeedFieldValues(
		srcID string,
		fieldValues fvalue.Values,
		value reflect.Value,
	)
	Fill(
		value reflect.Value,
		layers []fvalue.Values,
	) (plocation.Locations, error)
}
type GetListInterface interface {
	ApplyOn(g ifaces.Getter) (fvalue.Values, error)
	Push(o GetOp)
	Count() int
}
