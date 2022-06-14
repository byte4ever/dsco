package dsco

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/utils"
)

// Filler represents a structure Filler.
type Filler struct {
	layers layersIFace
	report reportInterface
}

// ErrInvalidLayers represents am error when layers are invalid at Filler
// creation time.
var ErrInvalidLayers = errors.New("invalid layers")

func checkLayers(layers []Binder2) error {
	layersLength := len(layers)
	if layersLength < 1 {
		return fmt.Errorf("no layers: %w", ErrInvalidLayers)
	}

	nilIndexes := make([]int, 0, layersLength)

	for i, layer := range layers {
		if layer == nil {
			nilIndexes = append(nilIndexes, i)
		}
	}

	switch nilIndexesLen := len(nilIndexes); nilIndexesLen {
	case 0:
		return nil

	case 1:
		return fmt.Errorf(
			"layer %s is nil: %w",
			utils.FormatIndexSequence(nilIndexes),
			ErrInvalidLayers,
		)

	default:
		return fmt.Errorf(
			"layers %s are nil: %w",
			utils.FormatIndexSequence(nilIndexes),
			ErrInvalidLayers,
		)
	}
}

// NewFiller creates a new filler using layers.
func NewFiller(l ...Binder2) (*Filler, error) {
	if err := checkLayers(l); err != nil {
		return nil, err
	}

	return &Filler{
		layers: layers(l),
		report: &Report{},
	}, nil
}

//goland:noinspection SpellCheckingInspection
func (filler *Filler) fillStruct(rootKey string, v reflect.Value) {
	t := v.Elem().Type()
	ve := v.Elem()

	for i := 0; i < ve.NumField(); i++ {
		field := ve.Field(i)
		fieldTyp := t.Field(i)

		key := utils.GetKeyName(rootKey, fieldTyp)

		//nolint:revive // will be removed soon
		switch fieldTyp.Type.String() {
		case "*time.Time":
			re := filler.layers.bind(key, field)

			if re.isFound() {
				ve.Field(i).Set(re.Value)
			}

			filler.report.addEntry(re)

			continue
		}

		e := fieldTyp.Type.Elem()
		if e.Kind() == reflect.Struct {
			fv := reflect.New(e)
			filler.fillStruct(
				key,
				fv,
			)

			ve.Field(i).Set(fv)

			continue
		}

		re := filler.layers.bind(key, field)

		if re.isFound() {
			ve.Field(i).Set(re.Value)
		}

		filler.report.addEntry(re)
	}
}

// Fill model based on layers. The parameter model must be a non nil interface
// and a non nil pointer to a struct.
func (filler *Filler) Fill(model interface{}) []error {
	if err := checkStruct(model); err != nil {
		return []error{err}
	}

	v := reflect.ValueOf(model)
	filler.fillStruct("", v)

	return filler.processReport()
}

func (filler *Filler) processReport() []error {
	errs := append([]error{}, filler.report.perEntryReport()...)
	errs = append(errs, filler.layers.getPostProcessErrors()...)

	return errs
}
