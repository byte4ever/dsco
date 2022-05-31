package goconf

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/byte4ever/goconf/utils"
)

// Layers is dummy...
type Layers []Binder

type ReportEntry struct {
	Key         string
	ExternalKey string
	Idx         int
	Errors      []error
}

type Report []ReportEntry

// Filler is dummy...
type Filler struct {
	layers Layers
	m      Report
}

// CheckedConfig is dummy
type CheckedConfig struct{}

// NewFiller is dummy
func NewFiller(
	_ *CheckedConfig,
	layers Layers,
) (
	*Filler,
	error,
) {
	return &Filler{layers: layers}, nil
}

//goland:noinspection SpellCheckingInspection
func (r *Filler) initializeStruct(rootKey string, t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)

		name := strings.Split(strings.Replace(ft.Tag.Get("yaml"), " ", "", -1), ",")[0]

		var s string
		if name != "" {
			s = name
		} else {
			s = utils.ToSnakeCase(ft.Name)
		}

		key := appendKey(rootKey, s)

		if ft.Type.Kind() != reflect.Ptr && ft.Type.Kind() != reflect.Slice {
			panic(fmt.Sprintf("opts support for type %v (%v)", key, ft))
		}

		switch ft.Type.String() {
		case
			"*time.Time":
			vv := f.Interface()
			vvValue := reflect.ValueOf(vv)

			if r.bind(key, ft.Type, &vvValue) {
				f.Set(vvValue)
			}

			continue
		}

		e := ft.Type.Elem()

		if e.Kind() == reflect.Struct {
			if f.IsNil() {
				fv := reflect.New(e)
				r.initializeStruct(
					key,
					e,
					fv.Elem(),
				)

				f.Set(fv)

				continue
			}

			r.initializeStruct(key, e, f.Elem())

			continue
		}

		vv := f.Interface()
		vvValue := reflect.ValueOf(vv)

		if r.bind(key, ft.Type, &vvValue) {
			f.Set(vvValue)
		}
	}

	return
}

func appendKey(a, b string) string {
	if a == "" {
		return b
	}

	return a + "-" + b
}

// Fill is dummy....
func (r *Filler) Fill(i interface{}) []error {
	if err := checkStruct(i); err != nil {
		return []error{err}
	}

	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	r.initializeStruct("", t.Elem(), v.Elem())

	return r.processReport()
}

var ErrUninitialized = errors.New("uninitialized")

func (r *Filler) processReport() (errs []error) {
	for _, entry := range r.m {
		for _, err := range entry.Errors {
			if err != nil {
				errs = append(errs, err)
			}
		}

		if entry.Idx == -1 {
			errs = append(errs, fmt.Errorf("key <%v>: %w", entry.Key, ErrUninitialized))
		}
	}

	for _, layer := range r.layers {
		if e := layer.GetPostProcessErrors(); len(e) > 0 {
			errs = append(errs, e...)
		}
	}

	return
}
