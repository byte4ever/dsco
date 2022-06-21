package walker

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/byte4ever/dsco/walker/cmdline"
	"github.com/byte4ever/dsco/walker/env"
)

const duplicateFmt = "layer #%d and #%d: %w"

// ErrCmdlineAlreadyUsed represent an error where ....
var ErrCmdlineAlreadyUsed = errors.New("cmdline already used")

// ErrDuplicateEnvPrefix represent an error where ....
var ErrDuplicateEnvPrefix = errors.New("duplicate env prefix")

// ErrDuplicateInputStruct represent an error where ....
var ErrDuplicateInputStruct = errors.New("duplicate input struct")

// ErrDuplicateStructID represent an error where ....
var ErrDuplicateStructID = errors.New("duplicate struct id")

type layerBuilder struct {
	idDedup  map[string]int
	builders []constraintLayerPolicy
}

type Layers []Layer

func (layers Layers) GetPolicies(fillReporter FillReporter) constraintLayerPolicies {
	bo := newLayerBuilder()

	for _, layer := range layers {
		err := layer.register(bo)
		if err != nil {
			fillReporter.ReportError(err)
			return nil //nolint:wrapcheck // error is clear enough
		}
	}

	return bo.builders
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

func newLayerBuilder() *layerBuilder {
	return &layerBuilder{
		builders: nil,
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
	wrap func(FieldValuesGetter) *constraintLayer,
	options []Option,
) error {
	if idx := to.dedupId("cmdLine"); idx != nil {
		return fmt.Errorf(
			duplicateFmt,
			*idx,
			to.curPos(),
			ErrCmdlineAlreadyUsed,
		)
	}

	builder, err := CmdLine(options...)
	if err != nil {
		return err
	}

	to.builders = append(to.builders, wrap(builder))

	return nil
}

func (o *StrictCmdlineLayer) register(to *layerBuilder) error {
	return wrapCmdlineBuild(to, strictLayer, o.options)
}

// WithStrictCmdlineLayer creates a strict command line layer.
// It can be used only once.
func WithStrictCmdlineLayer(options ...Option) *StrictCmdlineLayer {
	return &StrictCmdlineLayer{
		options: options,
	}
}

func (o *CmdlineLayer) register(to *layerBuilder) error {
	return wrapCmdlineBuild(to, normalLayer, o.options)
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
	wrap func(FieldValuesGetter) *constraintLayer,
	prefix string,
	options []Option,
) error {
	if idx := to.dedupId(
		fmt.Sprintf(
			"env(%s)",
			prefix,
		),
	); idx != nil {
		return fmt.Errorf(
			duplicateFmt, *idx, to.curPos(),
			ErrDuplicateEnvPrefix,
		)
	}

	builder, err := Env(prefix, options...)
	if err != nil {
		return err
	}

	to.addBuilder(wrap(builder))

	return nil
}

func (o *StrictEnvLayer) register(to *layerBuilder) error {
	return wrapEnvBuild(to, strictLayer, o.prefix, o.options)
}

// WithStrictEnvLayer creates a new strict environment layer.
func WithStrictEnvLayer(prefix string, options ...Option) *StrictEnvLayer {
	return &StrictEnvLayer{
		options: options,
		prefix:  prefix,
	}
}

func (o *EnvLayer) register(to *layerBuilder) error {
	return wrapEnvBuild(to, normalLayer, o.prefix, o.options)
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
	wrap func(FieldValuesGetter) *constraintLayer,
	input any,
	id string,
) error {
	ptr := reflect.ValueOf(input).Pointer()

	for _, issue := range []*struct {
		err error
		uid string
	}{
		{
			uid: fmt.Sprintf(
				"structPtr(%d)",
				ptr,
			),
			err: ErrDuplicateInputStruct,
		},
		{
			uid: fmt.Sprintf(
				"structId(%s)",
				id,
			),
			err: ErrDuplicateStructID,
		},
	} {
		if idx := to.dedupId(issue.uid); idx != nil {
			return fmt.Errorf(
				duplicateFmt, *idx, to.curPos(),
				issue.err,
			)
		}
	}

	builder, err := NewStructBuilder(input, id)
	if err != nil {
		return err
	}

	to.addBuilder(wrap(builder))

	return nil
}

func (o *StrictStructLayer) register(to *layerBuilder) error {
	return wrapStructBuild(to, strictLayer, o.input, o.id)
}

// WithStrictStructLayer creates a new strict structure layer.
func WithStrictStructLayer(input any, id string) *StrictStructLayer {
	return &StrictStructLayer{
		input: input,
		id:    id,
	}
}

func (o *StructLayer) register(to *layerBuilder) error {
	return wrapStructBuild(to, normalLayer, o.input, o.id)
}

// WithStructLayer creates a structure layer.
func WithStructLayer(input any, id string) *StructLayer {
	return &StructLayer{
		input: input,
		id:    id,
	}
}
