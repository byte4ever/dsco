package dsco

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco/internal/fvalue"
	"github.com/byte4ever/dsco/internal/ierror"
	"github.com/byte4ever/dsco/internal/merror"
	model2 "github.com/byte4ever/dsco/internal/model"
	"github.com/byte4ever/dsco/registry"
	"github.com/byte4ever/dsco/svalue"
)

// ErrInvalidType represent an error where ....
var ErrInvalidType = errors.New("invalid type")

// ErrParse represents an error indicating that a value cannot be parsed.
var ErrParse = errors.New("parse error")

type ParseError struct {
	Path     string
	Type     reflect.Type
	Location string
}

func (a ParseError) Error() string {
	return fmt.Sprintf(
		"parse error on %s-<%s> %s",
		a.Path,
		registry.LongTypeName(a.Type),
		a.Location,
	)
}

func (ParseError) Is(err error) bool {
	return errors.Is(err, ErrParse)
}

// ErrAliasCollision represents an error indicating that an alias is colliding
// with an actual key in the structure.
var ErrAliasCollision = errors.New("alias collision")

type AliasCollisionError struct {
	Path string
}

func (a AliasCollisionError) Error() string {
	return fmt.Sprintf("alias %s collides with structure", a.Path)
}

func (AliasCollisionError) Is(err error) bool {
	return errors.Is(err, ErrAliasCollision)
}

// ErrUnboundedLocation represents an error indicating that a key is never
// bound to the
// structure.
var ErrUnboundedLocation = errors.New("unbound key")

type UnboundedLocationErrors []UnboundedLocationError

func (u UnboundedLocationErrors) Len() int {
	return len(u)
}

func (u UnboundedLocationErrors) Less(i, j int) bool {
	return u[i].Location < u[j].Location
}

func (u UnboundedLocationErrors) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

type UnboundedLocationError struct {
	Location string
}

func (a UnboundedLocationError) Error() string {
	return fmt.Sprintf("unbounded location %s", a.Location)
}

func (UnboundedLocationError) Is(err error) bool {
	return errors.Is(err, ErrUnboundedLocation)
}

// ErrOverriddenKey represents an error indicating that a potential key binding
// wha overridden in another layer.
var ErrOverriddenKey = errors.New("overridden key")

type OverriddenKeyError struct {
	Path             string
	Location         string
	OverrideLocation string
}

func (a OverriddenKeyError) Error() string {
	return fmt.Sprintf(
		"for path %s %s is override by %s",
		a.Path,
		a.Location,
		a.OverrideLocation,
	)
}

func (OverriddenKeyError) Is(err error) bool {
	return errors.Is(err, ErrOverriddenKey)
}

// ErrNilProvider is shitty...
var ErrNilProvider = errors.New("nil provider")

// StringBasedBuilder is a value bases builder depending on text values.
type StringBasedBuilder struct {
	internalOpts
	values         svalue.Values
	expandedValues map[string]*fvalue.Value
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
			return ierror.IError{
				Index: i,
				Info:  "when processing option",
				Err:   err,
			}
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
	provider StringValuesProvider,
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

	values := provider.GetStringValues()

	if len(internalOptions.aliases) == 0 {
		return &StringBasedBuilder{
			internalOpts:   internalOptions,
			values:         values,
			expandedValues: make(map[string]*fvalue.Value),
		}, nil
	}

	converted := make(svalue.Values, len(values))

	for n, value := range values {
		if target, found := internalOptions.aliases[n]; found {
			converted[target] = value
			continue
		}

		converted[n] = value
	}

	return &StringBasedBuilder{
		internalOpts:   internalOptions,
		values:         converted,
		expandedValues: make(map[string]*fvalue.Value),
	}, nil
}

func (s *StringBasedBuilder) Expand(
	path string,
	_type reflect.Type) (
	err error,
) {
	convertedPath := convert(path)

	// check for alias collisions
	if _, found := s.internalOpts.aliases[convertedPath]; found {
		return AliasCollisionError{
			Path: path,
		}
	}

	entry, found := s.values[convertedPath]
	if !found {
		return nil
	}

	delete(s.values, convertedPath)

	tp := reflect.New(_type.Elem())

	err = yaml.Unmarshal(
		[]byte(entry.Value), tp.Interface(),
	)
	if err != nil {
		return ParseError{
			path,
			_type,
			entry.Location,
		}
	}

	model, err := model2.NewModel(_type)
	if err != nil {
		return fmt.Errorf("when expanding: %w", err)
	}

	valuesFor := model.GetFieldValuesFor(entry.Location, tp)
	for _, value := range valuesFor {
		s.expandedValues[strings.Join([]string{path, value.Path}, ".")] = value
	}

	return nil
}

func (s *StringBasedBuilder) Get(
	path string,
	_type reflect.Type,
) (
	fieldValue *fvalue.Value,
	err error,
) {
	convertedPath := convert(path)

	// check for alias collisions
	if _, found := s.internalOpts.aliases[convertedPath]; found {
		return nil, AliasCollisionError{
			Path: path,
		}
	}

	expandedEntry, found := s.expandedValues[path]
	if found {
		delete(s.expandedValues, path)
		return expandedEntry, nil
	}

	entry, found := s.values[convertedPath]
	if !found {
		return nil, nil //nolint:nilnil // required when nothing is found
	}

	switch _type.Kind() { //nolint:exhaustive // it's expected
	case reflect.Pointer:
		tp := reflect.New(_type.Elem())

		delete(s.values, convertedPath)

		if err := yaml.Unmarshal(
			[]byte(entry.Value), tp.Interface(),
		); err != nil {
			return nil, ParseError{
				path,
				_type,
				entry.Location,
			}
		}

		return &fvalue.Value{
			Value:    tp,
			Location: entry.Location,
		}, nil

	case reflect.Slice:
		tp := reflect.New(_type)

		delete(s.values, convertedPath)

		if err := yaml.Unmarshal(
			[]byte(entry.Value), tp.Interface(),
		); err != nil {
			return nil, ParseError{
				path,
				_type,
				entry.Location,
			}
		}

		return &fvalue.Value{
			Value:    tp.Elem(),
			Location: entry.Location,
		}, nil

	default:
		return nil, fmt.Errorf(
			"%s: %w",
			path,
			ErrInvalidType,
		)
	}
}

type GetError struct {
	merror.MError
}

var ErrGet = errors.New("")

func (GetError) Is(err error) bool {
	return errors.Is(err, ErrGet)
}

// GetBaseFor creates the bases.
func (s *StringBasedBuilder) GetFieldValuesFrom(
	model ModelInterface,
) (
	fvalue.Values,
	error,
) {
	var errs GetError

	if err := model.Expand(s); err != nil {
		errs.Add(err)
	}

	result, e := model.ApplyOn(s)
	if e != nil {
		errs.Add(e)
	}

	var e2s UnboundedLocationErrors

	for _, v := range s.values {
		e2s = append(
			e2s,
			UnboundedLocationError{
				Location: v.Location,
			},
		)
	}

	for _, v := range s.expandedValues {
		e2s = append(
			e2s,
			UnboundedLocationError{
				Location: v.Location,
			},
		)
	}

	if len(e2s) > 0 {
		sort.Sort(e2s)

		for _, e2 := range e2s {
			errs.Add(e2)
		}
	}

	if errs.None() {
		return result, nil
	}

	return nil, errs
}
