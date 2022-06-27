package walker

import (
	"reflect"
)

// FieldValuesGetter defines the ability to get a path/value set (bases).
type FieldValuesGetter interface {
	GetFieldValues(model ModelInterface) (FieldValues, []error)
}

type FillReporter interface {
	ReportUse(uid uint, path string, location string)
	ReportUnused(path string)
	ReportOverride(uid uint, location string)
	Result() (FillReport, error)
	ReportError(err error)
	Failed() bool
}

type Node interface {
	BuildGetList(s *GetList)
	FeedFieldValues(
		srcID string,
		fieldValues FieldValues,
		value reflect.Value,
	)
	Fill(
		fillReporter FillReporter,
		value reflect.Value,
		layers []FieldValues,
	)
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
	GetPolicies(fillReporter FillReporter) constraintLayerPolicies
}

type ModelInterface interface {
	TypeName() string
	ApplyOn(g Getter) (FieldValues, []error)
	FeedFieldValues(id string, v reflect.Value) FieldValues
	Fill(
		fillReporter FillReporter,
		inputModelValue reflect.Value,
		layers []FieldValues,
	)
}
