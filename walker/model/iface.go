package model

import (
	"reflect"

	"github.com/byte4ever/dsco/walker/fvalues"
	"github.com/byte4ever/dsco/walker/plocation"
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
