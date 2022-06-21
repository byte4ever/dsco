package walker

import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/walker/svalues"
)

// ErrKeyNotFound represent an error where ....
// var ErrKeyNotFound = errors.New("key not found")

// ErrNotUnused represent an error where ....
// var ErrNotUnused = errors.New("not unused")

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

type assignedValue struct {
	path     string
	location string
	value    *reflect.Value
}

// Base is a set of values
type Base map[int]assignedValue

// StringBasedBuilder is a value base builder depending on text values.
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

// NewStringBasedBuilder creates a base builder for the provided path/text
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

// GetBaseFor creates the base.
func (s *StringBasedBuilder) GetBaseFor(
	inputModel any,
) (Base, []error) {
	const (
		errFmt  = "%s: %w"
		errFmt2 = "%s [%s]: %w"
	)

	var (
		errs  []error
		maxId int
	)

	model := reflect.New(reflect.ValueOf(inputModel).Type().Elem())
	result := make(Base)

	wlkr := walker{
		fieldAction: func(
			id int,
			path string,
			value *reflect.Value,
		) error {
			convertedPath := convert(path)

			// check for alias collisions
			if _, found := s.internalOpts.aliases[convertedPath]; found {
				errs = append(
					errs,
					fmt.Errorf(errFmt, path, ErrAliasCollision),
				)
				return nil
			}

			entry, found := s.values[convertedPath]
			if !found {
				return nil
			}

			var tp reflect.Value

			dstType := value.Type()

			switch dstType.Kind() { //nolint:exhaustive // it's expected
			case reflect.Pointer:
				tp = reflect.New(dstType.Elem())

				if err := yaml.Unmarshal(
					[]byte(entry.Value), tp.Interface(),
				); err != nil {
					errs = append(
						errs,
						fmt.Errorf(
							"%s-<%s> %s: %w",
							path,
							dsco.LongTypeName(dstType),
							entry.Location,
							ErrParse,
						),
					)

					return nil
				}

				delete(s.values, convertedPath)

				result[id] = assignedValue{
					path:     path,
					location: entry.Location,
					value:    &tp,
				}

			case reflect.Slice:
				tp = reflect.New(dstType)

				if err := yaml.Unmarshal(
					[]byte(entry.Value), tp.Interface(),
				); err != nil {
					errs = append(
						errs,
						fmt.Errorf(
							errFmt2,
							path,
							entry.Location,
							ErrParse,
						),
					)

					return nil
				}

				delete(s.values, convertedPath)

				te := tp.Elem()
				result[id] = assignedValue{
					path:     path,
					location: entry.Location,
					value:    &te,
				}

			default:

				errs = append(
					errs,
					fmt.Errorf(
						errFmt,
						path,
						ErrInvalidType,
					),
				)

				return nil
			}

			return nil
		},
	}

	err := wlkr.walkRec(
		&maxId,
		"",
		model,
	)

	if err != nil {
		errs = append(errs, err)
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
