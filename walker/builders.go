package walker

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/byte4ever/dsco/walker/cmdline"
	"github.com/byte4ever/dsco/walker/env"
)

func CmdLine(options ...Option) (
	*StringBasedBuilder,
	error,
) {
	cmdLine, err := cmdline.NewEntriesProvider(
		os.Args[1:],
	)
	if err != nil {
		return nil, err
	}

	return NewStringBasedBuilder(cmdLine, options...)
}

func Env(prefix string, options ...Option) (
	*StringBasedBuilder,
	error,
) {
	envProvider1, err := env.NewEntriesProvider(prefix)
	if err != nil {
		return nil, err
	}

	return NewStringBasedBuilder(envProvider1, options...)
}

type layerBuilder struct {
	builders []ConstraintLayerPolicy
	idDedup  map[string]int
}

func newLayerBuilder() *layerBuilder {
	return &layerBuilder{
		builders: nil,
		idDedup:  make(map[string]int),
	}
}

func (o *layerBuilder) addBuilder(b ConstraintLayerPolicy) {
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

type Layer interface {
	register(to *layerBuilder) error
}

type StrictCmdlineLayer struct {
	options []Option
}

// ErrCmdlineAlreadyUsed represent an error where ....
var ErrCmdlineAlreadyUsed = errors.New("cmdline already used")

func wrapCmdlineBuild(
	to *layerBuilder,
	wrap func(BaseGetter) *ConstraintLayer,
	options []Option,
) error {
	if idx := to.dedupId("cmdLine"); idx != nil {
		return fmt.Errorf(
			"defined #%d and #%d: %w",
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
	return wrapCmdlineBuild(to, StrictLayer, o.options)
}

func WithStrictCmdlineLayer(options ...Option) *StrictCmdlineLayer {
	return &StrictCmdlineLayer{
		options: options,
	}
}

type CmdlineLayer struct {
	options []Option
}

func (o *CmdlineLayer) register(to *layerBuilder) error {
	return wrapCmdlineBuild(to, NormalLayer, o.options)
}

func WithCmdlineLayer(options ...Option) *CmdlineLayer {
	return &CmdlineLayer{
		options: options,
	}
}

// ///////////////////////////////////////////////////

type StrictEnvLayer struct {
	options []Option
	prefix  string
}

// ErrDuplicateEnvPrefix represent an error where ....
var ErrDuplicateEnvPrefix = errors.New("duplicate env prefix")

func wrapEnvBuild(
	to *layerBuilder,
	wrap func(BaseGetter) *ConstraintLayer,
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
			"layer #%d and #%d: %w", *idx, to.curPos(),
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
	return wrapEnvBuild(to, StrictLayer, o.prefix, o.options)
}

func WithStrictEnvLayer(prefix string, options ...Option) *StrictEnvLayer {
	return &StrictEnvLayer{
		options: options,
		prefix:  prefix,
	}
}

type EnvLayer struct {
	options []Option
	prefix  string
}

func (o *EnvLayer) register(to *layerBuilder) error {
	return wrapEnvBuild(to, NormalLayer, o.prefix, o.options)
}

func WithEnvLayer(prefix string, options ...Option) *EnvLayer {
	return &EnvLayer{
		options: options,
		prefix:  prefix,
	}
}

// ///////////////////////////////////////////////////////////////////

type StrictStructLayer struct {
	input any
	id    string
}

// ErrDuplicateStructPrefix represent an error where ....
var ErrDuplicateInputStruct = errors.New("duplicate input struct")

// ErrDuplicateStructPrefix represent an error where ....
var ErrDuplicateStructID = errors.New("duplicate struct id")

func wrapStructBuild(
	to *layerBuilder,
	wrap func(BaseGetter) *ConstraintLayer,
	input any,
	id string,
) error {
	ptr := reflect.ValueOf(input).Pointer()

	for _, e := range []struct {
		uid string
		err error
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
		if idx := to.dedupId(e.uid); idx != nil {
			return fmt.Errorf(
				"layer #%d and #%d: %w", *idx, to.curPos(),
				e.err,
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
	return wrapStructBuild(to, StrictLayer, o.input, o.id)
}

func WithStrictStructLayer(input any, id string) *StrictStructLayer {
	return &StrictStructLayer{
		input: input,
		id:    id,
	}
}

type StructLayer struct {
	input any
	id    string
}

func (o *StructLayer) register(to *layerBuilder) error {
	return wrapStructBuild(to, NormalLayer, o.input, o.id)
}

func WithStructLayer(input any, id string) *StructLayer {
	return &StructLayer{
		input: input,
		id:    id,
	}
}
