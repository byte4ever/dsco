package fvalues

import (
	"reflect"
)

type FieldValue struct {
	Value    reflect.Value
	Location string
}

type FieldValues map[uint]*FieldValue
