package dsco

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/byte4ever/dsco/utils"
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
func (r *Filler) fillStruct(rootKey string, v reflect.Value) {
	t := v.Elem().Type()
	v = v.Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)

		s := getName(ft)

		key := appendKey(rootKey, s)

		switch ft.Type.String() {
		case
			"*time.Time":
			vv := f.Interface()
			vvValue := reflect.ValueOf(vv)

			if r.bind(key, &vvValue) {
				f.Set(vvValue)
			}

			continue
		}

		e := ft.Type.Elem()

		if e.Kind() == reflect.Struct {
			if f.IsNil() {
				fv := reflect.New(e)
				r.fillStruct(
					key,
					fv,
				)

				f.Set(fv)

				continue
			}

			r.fillStruct(key, f)

			continue
		}

		vv := f.Interface()
		vvValue := reflect.ValueOf(vv)

		if r.bind(key, &vvValue) {
			f.Set(vvValue)
		}
	}

	return
}

func getName(ft reflect.StructField) string {
	name := strings.Split(strings.Replace(ft.Tag.Get("yaml"), " ", "", -1), ",")[0]

	var s string
	if name != "" {
		s = name
	} else {
		s = utils.ToSnakeCase(ft.Name)
	}

	return s
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

	v := reflect.ValueOf(i)
	r.fillStruct("", v)

	return r.processReport()
}

var ErrUninitialized = errors.New("uninitialized")

func (r *Filler) processReport() (errs []error) {
	errs = r.perEntryReport(errs)
	errs = r.perLayerReport(errs)

	return
}

func (r *Filler) perLayerReport(errs []error) []error {
	for _, layer := range r.layers {
		if e := layer.GetPostProcessErrors(); len(e) > 0 {
			errs = append(errs, e...)
		}
	}

	return errs
}

func (r *Filler) perEntryReport(errs []error) []error {
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

	return errs
}
