package walker

import (
	"reflect"
)

// FieldValuesGetter defines the ability to get a path/value set (bases).
type FieldValuesGetter interface {
	GetFieldValues(
		model *Model,
	) (FieldValues, []error)
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

type GetList []GetOp

func (s GetList) ApplyOn(g Getter) (FieldValues, []error) {
	var errs []error

	res := make(FieldValues, len(s))

	for _, op := range s {
		uid, fieldValue, err := op(g)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		if fieldValue != nil {
			res[uid] = fieldValue
		}
	}

	return res, errs
}

func (s *GetList) Push(o GetOp) {
	*s = append(*s, o)
}

type GetOp func(g Getter) (uid uint, fieldValue *FieldValue, err error)

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
