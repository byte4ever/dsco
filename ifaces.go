package dsco

import (
	"reflect"
)

// TODO :- lmartin 6/10/22 -: need some refactoring - to big and complex.
//  interface Bind should only take a type, it returns too many values.

type layersIFace interface {
	bind(
		key string,
		dstValue reflect.Value,
	) ReportEntry
	getPostProcessErrors() []error
}
type reportInterface interface {
	perEntryReport() (errs []error)
	addEntry(e ReportEntry)
}

// BindingAttempt is a bounded value.
type BindingAttempt struct {
	Error    error
	Value    reflect.Value
	Location string
}

// Binder2 defines new binder behaviour (simpler).
type Binder2 interface {
	Bind(
		key string,
		dstType reflect.Type,
	) BindingAttempt

	Use(
		key string,
	) error

	Errors() []error
}

// HasValue returns true when bounding attempts value is available.
func (ba *BindingAttempt) HasValue() bool {
	return ba.Value.IsValid()
}

// Binder2 defines the ability to bind/create a value based on a given key. When
// set is true the value is actually created otherwise it will only perform all
// checks and value is not created.
