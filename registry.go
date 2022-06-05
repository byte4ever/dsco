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

var ErrInvalidLayers = errors.New("invalid layers")

// NewFiller is dummy
func NewFiller(layers ...Binder) (*Filler, error) {
	if len(layers) < 1 {
		return nil, fmt.Errorf("at least on layer MUST be provided: %w", ErrInvalidLayers)
	}

	return &Filler{layers: layers}, nil
}

//goland:noinspection SpellCheckingInspection
func (r *Filler) fillStruct(rootKey string, v reflect.Value) {
	t := v.Elem().Type()
	ve := v.Elem()

	for i := 0; i < ve.NumField(); i++ {
		f := ve.Field(i)
		ft := t.Field(i)
		s := getName(ft)
		key := appendKey(rootKey, s)

		switch ft.Type.String() {
		case "*time.Time":
			if r.bind(key, &f) {
				ve.Field(i).Set(f)
			}

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

		if r.bind(key, &f) {
			ve.Field(i).Set(f)
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
