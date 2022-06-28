package walker

import (
	"reflect"
)

// FieldValuesGetter defines the ability to get a path/value set (bases).
type FieldValuesGetter interface {
	GetFieldValues(model ModelInterface) (FieldValues, []error)
}

type Node interface {
	BuildGetList(s *GetList)
	FeedFieldValues(
		srcID string,
		fieldValues FieldValues,
		value reflect.Value,
	)
	Fill(
		value reflect.Value,
		layers []FieldValues,
	) (PathLocations, error)
}

type FieldValue struct {
	value    reflect.Value
	location string
}

type FieldValues map[uint]*FieldValue

type Getter interface {
	Get(
		path string,
		_type reflect.Type,
	) (
		fieldValue *FieldValue,
		err error,
	)
}

type PoliciesGetter interface {
	GetPolicies() (constraintLayerPolicies, error)
}

type ModelInterface interface {
	TypeName() string
	ApplyOn(g Getter) (FieldValues, []error)
	FeedFieldValues(id string, v reflect.Value) FieldValues
	Fill(
		inputModelValue reflect.Value,
		layers []FieldValues,
	) (PathLocations, error)
}
