package dsco

import (
	"reflect"
)

var _ layersIFace = layers{}

type layers []Binder2

func (l layers) bind(
	key string,
	dstValue reflect.Value,
) ReportEntry {
	var (
		idxFound         = -1
		externalKeyFound = ""
		outVal           reflect.Value
	)

	dstValType := dstValue.Type()

	errs := make([]error, 0, len(l))

	for idx, binder := range l {
		bindingAttempt := binder.Bind(key, dstValType)
		if bindingAttempt.Error == nil &&
			idxFound == -1 &&
			bindingAttempt.HasValue() {
			idxFound = idx
			externalKeyFound = bindingAttempt.Location
			outVal = bindingAttempt.Value

			if err := binder.Use(key); err != nil {
				panic(err)
			}
		}

		errs = append(errs, bindingAttempt.Error)
	}

	return ReportEntry{
		Value:       outVal,
		Key:         key,
		ExternalKey: externalKeyFound,
		Idx:         idxFound,
		Errors:      errs,
	}
}

func (l layers) getPostProcessErrors() []error {
	var errs []error

	for _, layer := range l {
		if e := layer.Errors(); len(e) > 0 {
			errs = append(errs, e...)
		}
	}

	return errs
}
