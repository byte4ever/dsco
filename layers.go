package dsco

import (
	"reflect"
)

func (r *Filler) bind(
	key string,
	dstValue *reflect.Value,
) bool {
	var (
		e                []error
		idxFound         = -1
		ExternalKeyFound = ""
	)

	for idx, binder := range r.layers {
		_, keyOut, success, err := binder.Bind(key, idxFound == -1, dstValue)

		if err == nil && idxFound == -1 && success {
			idxFound = idx
			ExternalKeyFound = keyOut
		}

		e = append(e, err)
	}

	r.m = append(
		r.m, ReportEntry{
			Key:         key,
			ExternalKey: ExternalKeyFound,
			Idx:         idxFound,
			Errors:      e,
		},
	)

	return idxFound != -1
}
