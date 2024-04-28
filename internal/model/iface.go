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
	BuildExpandList(e *ExpandList)
}

type GetListInterface interface {
	ApplyOn(g internal.ValueGetter) (fvalue.Values, error)
	Push(o GetOp)
	Count() int
}

type ExpandListInterface interface {
	ApplyOn(g internal.StructExpander) error
	Push(o ExpandOp)
	Count() int
}

type GetOp func(g internal.ValueGetter) (
	uid uint,
	fieldValue *fvalue.Value,
	err error,
)

type ExpandOp func(g internal.StructExpander) (
	err error,
)
