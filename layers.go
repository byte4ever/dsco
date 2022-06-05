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
		e                []error
		idxFound         = -1
		ExternalKeyFound = ""
	)

	for idx, binder := range l {
		_, keyOut, success, err := binder.Bind(key, idxFound == -1, &dstValue)

		if err == nil && idxFound == -1 && success {
			idxFound = idx
			ExternalKeyFound = keyOut
		}

		e = append(e, err)
	}

	return ReportEntry{
		Value:       dstValue,
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
