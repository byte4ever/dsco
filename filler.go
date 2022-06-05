package dsco

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/utils"
)

type Filler struct {
	layers layersIFace
	report reportIface
}

var ErrInvalidLayers = errors.New("invalid layers")

func NewFiller(l ...Binder) (*Filler, error) {
	if len(l) < 1 {
		return nil, fmt.Errorf("at least on layer MUST be provided: %w", ErrInvalidLayers)
	}

	r := Report{}

	return &Filler{layers: layers(l), report: &r}, nil
}

//goland:noinspection SpellCheckingInspection
func (r *Filler) fillStruct(rootKey string, v reflect.Value) {
	t := v.Elem().Type()
	ve := v.Elem()

	for i := 0; i < ve.NumField(); i++ {
		f := ve.Field(i)
		ft := t.Field(i)

		key := utils.GetKeyName(rootKey, ft)

		switch ft.Type.String() {
		case "*time.Time":
			re := r.layers.bind(key, f)

			if re.isFound() {
				ve.Field(i).Set(re.Value)
			}

			r.report.addEntry(re)

			continue
		}

		e := ft.Type.Elem()
		if e.Kind() == reflect.Struct {
			fv := reflect.New(e)
			r.fillStruct(
				key,
				fv,
			)

			ve.Field(i).Set(fv)

			continue
		}

		re := r.layers.bind(key, f)

		if re.isFound() {
			ve.Field(i).Set(re.Value)
		}

		r.report.addEntry(re)
	}

	return
}

func (r *Filler) Fill(i interface{}) []error {
	if err := checkStruct(i); err != nil {
		return []error{err}
	}

	v := reflect.ValueOf(i)
	r.fillStruct("", v)

	return r.processReport()
}

var ErrUninitialized = errors.New("uninitialized")

func (r *Filler) processReport() (errs []error) {
	errs = append(errs, r.report.perEntryReport()...)
	errs = append(errs, r.layers.getPostProcessErrors()...)

	return
}
