package walker

import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/walker/svalues"
)

// ErrInvalidType represent an error where ....
var ErrInvalidType = errors.New("invalid type")

// ErrParse represents an error indicating that a value cannot be parsed.
var ErrParse = errors.New("parse error")

// ErrAliasCollision represents an error indicating that an alias is colliding
// with an actual key in the structure.
var ErrAliasCollision = errors.New("alias collision")

// ErrUnboundKey represents an error indicating that a key is never bound to the
// structure.
var ErrUnboundKey = errors.New("unbound key")

// ErrOverriddenKey represents an error indicating that a potential key binding
// wha overridden in another layer.
var ErrOverriddenKey = errors.New("overridden key")

// ErrNilProvider is shitty...
var ErrNilProvider = errors.New("nil provider")

// StringBasedBuilder is a value bases builder depending on text values.
type StringBasedBuilder struct {
	internalOpts
	values svalues.StringValues
}

// ErrNoAliasesProvided represent an error where no aliases map was
// provided with option.
var ErrNoAliasesProvided = errors.New("no aliases provided")

// Option is processing option for string based binder.
type Option interface {
	apply(opts *internalOpts) error
}

type internalOpts struct {
	aliases map[string]string
}

// AliasesOption defines keys aliasing.
type AliasesOption map[string]string

func (o *internalOpts) applyOptions(options []Option) error {
	for i, option := range options {
		if err := option.apply(o); err != nil {
			return fmt.Errorf(
				"when processing option #%d: %w",
				i,
				err,
			)
		}
	}

	return nil
}

func (a AliasesOption) apply(opts *internalOpts) error {
	if len(a) > 0 {
		opts.aliases = a
		return nil
	}

	return ErrNoAliasesProvided
}

// WithAliases returns a keys aliasing option.
func WithAliases(mapping map[string]string) AliasesOption {
	return mapping
}

// NewStringBasedBuilder creates a bases builder for the provided path/text
// value set.
func NewStringBasedBuilder(
	provider svalues.StringValuesProvider,
	options ...Option,
) (
	*StringBasedBuilder,
	error,
) {
	if provider == nil {
		return nil, ErrNilProvider
	}

	internalOptions := internalOpts{}

	if err := internalOptions.applyOptions(options); err != nil {
		return nil, err
	}

	return &StringBasedBuilder{
		internalOpts: internalOptions,
		values:       provider.GetStringValues(),
	}, nil
}

func (s *StringBasedBuilder) Get(
	path string, _type reflect.Type,
) (fieldValue *FieldValue, err error) {
	const (
		errFmt  = "%s: %w"
		errFmt2 = "%s [%s]: %w"
	)

	convertedPath := convert(path)

	// check for alias collisions
	if _, found := s.internalOpts.aliases[convertedPath]; found {
		return nil, fmt.Errorf(
			errFmt,
			path,
			ErrAliasCollision,
		)
	}

	entry, found := s.values[convertedPath]
	if !found {
		return nil, nil
	}

	switch _type.Kind() { //nolint:exhaustive // it's expected
	case reflect.Pointer:
		tp := reflect.New(_type.Elem())

		if err := yaml.Unmarshal(
			[]byte(entry.Value), tp.Interface(),
		); err != nil {
			return nil, fmt.Errorf(
				"%s-<%s> %s: %w",
				path,
				dsco.LongTypeName(_type),
				entry.Location,
				ErrParse,
			)
		}

		delete(s.values, convertedPath)

		return &FieldValue{
			value:    tp,
			location: entry.Location,
		}, nil

	case reflect.Slice:
		tp := reflect.New(_type)

		if err := yaml.Unmarshal(
			[]byte(entry.Value), tp.Interface(),
		); err != nil {

			return nil, fmt.Errorf(
				errFmt2,
				path,
				entry.Location,
				ErrParse,
			)
		}

		delete(s.values, convertedPath)

		return &FieldValue{
			value:    tp.Elem(),
			location: entry.Location,
		}, nil

	default:
		return nil, fmt.Errorf(
			errFmt,
			path,
			ErrInvalidType,
		)
	}
}

// GetBaseFor creates the bases.
func (s *StringBasedBuilder) GetFieldValues(
	model *Model,
) (FieldValues, []error) {
	const errFmt = "%s: %w"
	var errs []error

	result, e := model.ApplyOn(s)

	if len(e) > 0 {
		errs = append(errs, e...)
	}

	for _, v := range s.values {
		errs = append(
			errs, fmt.Errorf(
				errFmt,
				v.Location,
				ErrUnboundKey,
			),
		)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return result, nil
}
