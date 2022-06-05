package dsco

import (
	"reflect"
)

type Binder interface {
	Bind(
		key string,
		set bool,
		dstValue *reflect.Value,
	) (
		origin Origin,
		keyOut string,
		succeed bool,
		err error,
	)
	GetPostProcessErrors() []error
}

type layersIFace interface {
	bind(
		key string,
		dstValue *reflect.Value,
	) ReportEntry
	getPostProcessErrors() []error
}
type reportIface interface {
	perEntryReport() (errs []error)
	addEntry(e ReportEntry)
}
