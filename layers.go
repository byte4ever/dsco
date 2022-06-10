package dsco

import (
	"reflect"
)

var _ layersIFace = layers{}

type layers []Binder

func (l layers) bind(
	key string,
	dstValue reflect.Value,
) ReportEntry {
	var (
		idxFound         = -1
		externalKeyFound = ""
		outVal           reflect.Value
	)

	errs := make([]error, 0, len(l))

	for idx, binder := range l {
		// todo :- lmartin 6/5/22 -: to many results here, should be simplified
		_, keyOut, success, v, err := binder.Bind(key, idxFound == -1, dstValue)

		if err == nil && idxFound == -1 && success {
			idxFound = idx
			externalKeyFound = keyOut
			outVal = v
		}

		errs = append(errs, err)
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
		if e := layer.GetPostProcessErrors(); len(e) > 0 {
			errs = append(errs, e...)
		}
	}

	return errs
}
