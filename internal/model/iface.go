package model

import (
	"reflect"

	"github.com/byte4ever/dsco/internal"
	"github.com/byte4ever/dsco/internal/fvalue"
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
	ApplyOn(g internal.Getter) (fvalue.Values, error)
	Push(o GetOp)
	Count() int
}

type GetOp func(g internal.Getter) (
	uid uint,
	fieldValue *fvalue.Value,
	err error,
)
