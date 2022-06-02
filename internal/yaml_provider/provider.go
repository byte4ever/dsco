package yaml_provider

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"gopkg.in/yaml.v3"
)

var ErrNilInterfaces = errors.New("nil interfaces")

type ReadCloseProvider interface {
	ReadClose(perform func(r io.Reader) error) error
}

type Provider struct {
	i interface{}
}

func (p *Provider) GetInterface() (interface{}, error) {
	return p.i, nil
}

func Provide(model interface{}, rcProvider ReadCloseProvider) (*Provider, error) {
	if model == nil {
		return nil, ErrNilInterfaces
	}

	k := reflect.New(reflect.TypeOf(model).Elem()).Interface()

	err := rcProvider.ReadClose(
		func(reader io.Reader) error {
			dec := yaml.NewDecoder(
				reader,
			)

			if err := dec.Decode(k); err != nil {
				return fmt.Errorf("while parsing yaml buffer: %w", err)
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
