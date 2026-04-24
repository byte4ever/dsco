package dsco

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco/internal/fvalue"
	"github.com/byte4ever/dsco/internal/ierror"
	"github.com/byte4ever/dsco/internal/merror"
	"github.com/byte4ever/dsco/internal/model"
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
	return "unbounded location " + a.Location
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
	keyFormatter   KeyFormatter
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

// newStringBasedBuilderWithFormatter is an internal constructor that
// behaves like NewStringBasedBuilder but records the KeyFormatter used to
// render the layer's keys in inventory reports.
func newStringBasedBuilderWithFormatter(
	provider StringValuesProvider,
	formatter KeyFormatter,
	options ...Option,
) (*StringBasedBuilder, error) {
	builder, err := NewStringBasedBuilder(provider, options...)
	if err != nil {
		return nil, err //nolint:wrapcheck // same-package constructor
	}

	builder.keyFormatter = formatter

	return builder, nil
}

func (s *StringBasedBuilder) ExpandStruct(
	path string,
	_type reflect.Type) (
	err error,
) {
	convertedPath := convert(path)

	// check for alias collisions
	if _, found := s.aliases[convertedPath]; found {
		return &AliasCollisionError{
			Path: path,
		}
	}

	entryToExpand, found := s.values[convertedPath]
	if !found {
		return nil
	}

	delete(s.values, convertedPath)

	tp := reflect.New(_type.Elem())

	// parse yaml struct
	if err = yaml.Unmarshal(
		[]byte(entryToExpand.Value), tp.Interface(),
	); err != nil {
		return &ParseError{
			path,
			_type,
			entryToExpand.Location,
		}
	}

	extractedModel, err := model.NewModel(_type)
	if err != nil {
		return fmt.Errorf("when expanding: %w", err)
	}

	for _, value := range extractedModel.GetFieldValuesFor(
		entryToExpand.Location,
		tp,
	) {
		s.expandedValues[path+"."+value.Path] = value
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
	if _, found := s.aliases[convertedPath]; found {
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

// ErrNilKeyFormatter indicates that ReportInventory was called on a
// StringBasedBuilder with no KeyFormatter set.
var ErrNilKeyFormatter = errors.New("nil key formatter")

// ErrUnknownKeyFormatterKind is returned by NewStringBasedBuilderForTest
// when the kind argument is not a recognised formatter kind.
var ErrUnknownKeyFormatterKind = errors.New("unknown key-formatter kind")

// ReportInventory implements InventoryReporter by walking the model's
// alias map and rendering each entry through the layer's KeyFormatter.
// No I/O is performed.
func (s *StringBasedBuilder) ReportInventory(
	mdl ModelInterface,
) (LayerInventory, error) {
	const errCtx = "reporting inventory"

	if s.keyFormatter == nil {
		// Defensive: any builder constructed via the layer wrappers in
		// builders.go has a non-nil formatter.
		return LayerInventory{}, fmt.Errorf(
			"%s: %w", errCtx, ErrNilKeyFormatter,
		)
	}

	aliases, err := collectAliases(mdl)
	if err != nil {
		return LayerInventory{}, fmt.Errorf("%s: %w", errCtx, err)
	}

	provides := make([]FieldProvision, 0, len(aliases))
	for fieldUID, aliasPath := range aliases {
		provides = append(provides, FieldProvision{
			FieldUID: fieldUID,
			Key:      s.keyFormatter.FormatKey(aliasPath),
		})
	}

	inv := LayerInventory{
		Name:     s.keyFormatter.LayerName(),
		Provides: provides,
	}

	if s.keyFormatter.LayerKind() == "" {
		inv.Note = "custom provider — keys not enumerable"
		// Drop key strings — they cannot be rendered for custom providers.
		for i := range inv.Provides {
			inv.Provides[i].Key = ""
		}
	}

	return inv, nil
}

// aliasRecorder implements internal.ValueGetter to capture each leaf
// field's (FieldUID, alias-path) without producing any value.
type aliasRecorder struct {
	aliases map[string]string
}

// Get records the alias path for path and returns (nil, nil) so the
// model treats the field as unfilled. No I/O is performed.
func (r *aliasRecorder) Get(
	path string,
	_ reflect.Type,
) (*fvalue.Value, error) {
	r.aliases[path] = convert(path)
	return nil, nil //nolint:nilnil // matches StringBasedBuilder.Get when nothing is found
}

// collectAliases returns the field-uid → alias-path map for the given
// model. The alias path is dash-separated, lowercase, matching the form
// StringBasedBuilder.values keys use.
//
// It builds a recording ValueGetter that captures (path, convert(path))
// for every leaf the model iterates, then runs the model's ApplyOn pass.
// No values are loaded; the recorder always returns (nil, nil).
func collectAliases(mdl ModelInterface) (map[string]string, error) {
	const errCtx = "collecting aliases"

	rec := &aliasRecorder{
		aliases: make(map[string]string),
	}

	if _, err := mdl.ApplyOn(rec); err != nil {
		return nil, fmt.Errorf("%s: %w", errCtx, err)
	}

	return rec.aliases, nil
}

// NewStringBasedBuilderForTest constructs a StringBasedBuilder with a
// synthetic KeyFormatter. Intended solely for tests that need to
// exercise ReportInventory without going through the layer wrappers in
// builders.go.
//
// kind must be one of "env", "cmdline", "file", or "" (nil formatter
// for custom-provider behaviour). For "env" / "file", metaOrPrefix is
// the prefix or file id.
func NewStringBasedBuilderForTest(
	provider StringValuesProvider,
	kind, metaOrPrefix string,
) (*StringBasedBuilder, error) {
	var kf KeyFormatter

	switch kind {
	case "env":
		kf = newEnvKeyFormatter(metaOrPrefix)
	case "cmdline":
		kf = newCmdlineKeyFormatter()
	case "":
		kf = newNilKeyFormatter(metaOrPrefix)
	default:
		return nil, fmt.Errorf(
			"%w: %q", ErrUnknownKeyFormatterKind, kind,
		)
	}

	return newStringBasedBuilderWithFormatter(provider, kf)
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
	_model ModelInterface,
) (
	fvalue.Values,
	error,
) {
	var errs GetError

	if err := _model.Expand(s); err != nil {
		errs.Add(err)
	}

	result, e := _model.ApplyOn(s)
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
