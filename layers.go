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
		ExternalKeyFound = ""
		outVal           reflect.Value
	)

	e := make([]error, 0, len(l))

	for idx, binder := range l {
		// todo :- lmartin 6/5/22 -: to many results here, should be simplified
		_, keyOut, success, v, err := binder.Bind(key, idxFound == -1, dstValue)

		if err == nil && idxFound == -1 && success {
			idxFound = idx
			ExternalKeyFound = keyOut
			outVal = v
		}

		e = append(e, err)
	}

	return ReportEntry{
		Value:       outVal,
		Key:         key,
		ExternalKey: ExternalKeyFound,
		Idx:         idxFound,
		Errors:      e,
	}
}

func (l layers) getPostProcessErrors() (errs []error) {
	for _, layer := range l {
		if e := layer.GetPostProcessErrors(); len(e) > 0 {
			errs = append(errs, e...)
		}
	}

	return
}
