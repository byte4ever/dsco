package dsco

import (
	"fmt"
	"os"
	"reflect"

	"github.com/byte4ever/dsco/cmdline"
	"github.com/byte4ever/dsco/env"
	"github.com/byte4ever/dsco/ifaces"
	"github.com/byte4ever/dsco/internal/ierror"
)

type layerBuilder struct {
	idDedup  map[string]int
	builders []constraintLayerPolicy
}

type Layers []Layer

func (layers Layers) GetPolicies() (constraintLayerPolicies, error) {
	var errs LayerErrors

	bo := newLayerBuilder(len(layers))

	for index, layer := range layers {
		err := layer.register(bo)

		if err != nil {
			errs.Add(
				ierror.IError{
					Index: index,
					Info:  "layer",
					Err:   err,
				},
			)
		}
	}

	if errs.None() {
		return bo.builders, nil
	}

	return nil, errs
}

// Layer defines a configuration layer.
type Layer interface {
	register(to *layerBuilder) error
}

// StrictCmdlineLayer is a strict command line layer.
type StrictCmdlineLayer struct {
	options []Option
}

// CmdlineLayer is a command line layer.
type CmdlineLayer struct {
	options []Option
}

// StrictEnvLayer is a strict environnement variables layer.
type StrictEnvLayer struct {
	prefix  string
	options []Option
}

// EnvLayer is an environnement variable layer.
type EnvLayer struct {
	prefix  string
	options []Option
}

// StrictStructLayer is a strict structure layer.
type StrictStructLayer struct {
	input any
	id    string
}

// StructLayer is a structure layer.
type StructLayer struct {
	input any
	id    string
}

// CmdLine builds a command line manager.
func CmdLine(options ...Option) (
	*StringBasedBuilder,
	error,
) {
	cmdLine, err := cmdline.NewEntriesProvider(
		os.Args[1:],
	)
	if err != nil {
		return nil, fmt.Errorf("cmdline builder: %w", err)
	}

	return NewStringBasedBuilder(cmdLine, options...)
}

// Env builds a env manager.
func Env(prefix string, options ...Option) (
	*StringBasedBuilder,
	error,
) {
	envProvider1, err := env.NewEntriesProvider(prefix)
	if err != nil {
		return nil, fmt.Errorf("env builder: %w", err)
	}

	return NewStringBasedBuilder(envProvider1, options...)
}

func newLayerBuilder(l int) *layerBuilder {
	return &layerBuilder{
		builders: make(constraintLayerPolicies, 0, l),
		idDedup:  make(map[string]int),
	}
}

func (o *layerBuilder) addBuilder(b constraintLayerPolicy) {
	o.builders = append(o.builders, b)
}

func (o *layerBuilder) curPos() int {
	return len(o.builders)
}

func (o *layerBuilder) dedupId(id string) *int {
	if idx, found := o.idDedup[id]; found {
		return &idx
	}

	o.idDedup[id] = o.curPos()

	return nil
}

func wrapCmdlineBuild(
	to *layerBuilder,
	wrap func(ifaces.FieldValuesGetter) constraintLayerPolicy,
	options []Option,
) error {
	if idx := to.dedupId("cmdLine"); idx != nil {
		return CmdlineAlreadyUsedError{
			Index: *idx,
		}
	}

	builder, err := CmdLine(options...)
	if err != nil {
		return err
	}

	to.builders = append(to.builders, wrap(builder))

	return nil
}

func (o *StrictCmdlineLayer) register(to *layerBuilder) error {
	return wrapCmdlineBuild(to, newStrictLayer, o.options)
}

// WithStrictCmdlineLayer creates a strict command line layer.
// It can be used only once.
func WithStrictCmdlineLayer(options ...Option) *StrictCmdlineLayer {
	return &StrictCmdlineLayer{
		options: options,
	}
}

func (o *CmdlineLayer) register(to *layerBuilder) error {
	return wrapCmdlineBuild(to, newNormalLayer, o.options)
}

// WithCmdlineLayer creates a command line layer.
func WithCmdlineLayer(options ...Option) *CmdlineLayer {
	return &CmdlineLayer{
		options: options,
	}
}

// ///////////////////////////////////////////////////

func wrapEnvBuild(
	to *layerBuilder,
	wrap func(ifaces.FieldValuesGetter) constraintLayerPolicy,
	prefix string,
	options []Option,
) error {
	if idx := to.dedupId(
		fmt.Sprintf(
			"env(%s)",
			prefix,
		),
	); idx != nil {
		return DuplicateEnvPrefixError{
			Index:  *idx,
			Prefix: prefix,
		}
	}

	builder, err := Env(prefix, options...)
	if err != nil {
		return err
	}

	to.addBuilder(wrap(builder))

	return nil
}

func (o *StrictEnvLayer) register(to *layerBuilder) error {
	return wrapEnvBuild(to, newStrictLayer, o.prefix, o.options)
}

// WithStrictEnvLayer creates a new strict environment layer.
func WithStrictEnvLayer(prefix string, options ...Option) *StrictEnvLayer {
	return &StrictEnvLayer{
		options: options,
		prefix:  prefix,
	}
}

func (o *EnvLayer) register(to *layerBuilder) error {
	return wrapEnvBuild(to, newNormalLayer, o.prefix, o.options)
}

// WithEnvLayer creates an environment variable layer.
func WithEnvLayer(prefix string, options ...Option) *EnvLayer {
	return &EnvLayer{
		options: options,
		prefix:  prefix,
	}
}

// ///////////////////////////////////////////////////////////////////

func wrapStructBuild(
	to *layerBuilder,
	wrap func(ifaces.FieldValuesGetter) constraintLayerPolicy,
	input any,
	id string,
) error {
	builder, err := NewStructBuilder(input, id)
	if err != nil {
		return err
	}

	ptr := reflect.ValueOf(input).Pointer()

	if idx := to.dedupId(
		fmt.Sprintf(
			"structPtr(%d)",
			ptr,
		),
	); idx != nil {
		return DuplicateInputStructError{
			Index: *idx,
		}
	}

	if idx := to.dedupId(
		fmt.Sprintf(
			"structId(%s)",
			id,
		),
	); idx != nil {
		return DuplicateStructIDError{
			Index: *idx,
			ID:    id,
		}
	}

	to.addBuilder(wrap(builder))

	return nil
}

func (o *StrictStructLayer) register(to *layerBuilder) error {
	return wrapStructBuild(to, newStrictLayer, o.input, o.id)
}

// WithStrictStructLayer creates a new strict structure layer.
func WithStrictStructLayer(input any, id string) *StrictStructLayer {
	return &StrictStructLayer{
		input: input,
		id:    id,
	}
}

func (o *StructLayer) register(to *layerBuilder) error {
	return wrapStructBuild(to, newNormalLayer, o.input, o.id)
}

// WithStructLayer creates a structure layer.
func WithStructLayer(input any, id string) *StructLayer {
	return &StructLayer{
		input: input,
		id:    id,
	}
}
