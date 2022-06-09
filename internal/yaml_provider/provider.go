package yaml_provider

import (
	"fmt"
	"io"
	"reflect"

	"gopkg.in/yaml.v3"
)

// Provider represents an interface provider that decodes yaml from incoming
// read performer.
type Provider struct {
	i interface{}
}

// GetInterface implements InterfaceProvider interface.
func (p *Provider) GetInterface() (interface{}, error) {
	return p.i, nil
}

// New creates an interface provider based on the given model and read performer.
//
// First parameter model MUST not be nil and MUST refer to a pointer on struct.
//
// Second parameter functor MUST not be nil.
//
func New(model interface{}, functor ReaderFunctor) (*Provider, error) {
	err := checkModel(model)
	if err != nil {
		return nil, err
	}

	k := reflect.New(reflect.TypeOf(model).Elem()).Interface()

	if functor == nil {
		return nil, ErrNilReaderFunctor
	}

	err = functor.Apply(
		func(reader io.Reader) error {
			dec := yaml.NewDecoder(
				reader,
			)

			if err2 := dec.Decode(k); err2 != nil {
				return fmt.Errorf("while parsing yaml buffer: %w", err2)
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &Provider{
		i: k,
	}, nil
}

func checkModel(model interface{}) error {
	if model == nil {
		return fmt.Errorf("model is nil: %w", ErrInvalidModel)
	}

	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("model not a pointer: %w", ErrInvalidModel)
	}

	te := t.Elem()
	if te.Kind() != reflect.Struct {
		return fmt.Errorf("model not a pointer on struct: %w", ErrInvalidModel)
	}

	return nil
}
