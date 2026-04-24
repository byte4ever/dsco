package inventory

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/internal/fvalue"
)

type (
	// Report is the static inventory of a config struct against a layer set.
	Report struct {
		Type   string  `json:"type"   yaml:"type"`
		Fields []Field `json:"fields" yaml:"fields"`
	}

	// Field describes the canonical key (and any baked-in default) for one
	// leaf field of the config struct.
	Field struct {
		Satisfied *Satisfaction `json:"satisfied,omitempty" yaml:"satisfied,omitempty"`
		Key       *KeySpec      `json:"key,omitempty"       yaml:"key,omitempty"`
		Path      string        `json:"path"                yaml:"path"`
		GoType    string        `json:"go_type"             yaml:"go_type"`
	}

	// Satisfaction records that a struct layer already provides a value
	// for this field.
	Satisfaction struct {
		Value   any    `json:"value"    yaml:"value"`
		LayerID string `json:"layer_id" yaml:"layer_id"`
	}

	// KeySpec is the canonical (highest-precedence) key form a string-based
	// layer would accept to supply this field.
	KeySpec struct {
		Layer string `json:"layer" yaml:"layer"`
		Key   string `json:"key"   yaml:"key"`
	}

	// leaf describes one scalar leaf field of the config struct.
	leaf struct {
		path   string
		goType string
	}

	// leafRecorder implements internal.ValueGetter to capture (path, type)
	// for every leaf field the model iterates. No values are produced.
	leafRecorder struct {
		leaves []leaf
	}
)

// Get records a leaf entry and returns (nil, nil) so the model treats
// the field as unfilled.
func (r *leafRecorder) Get(
	path string,
	fieldType reflect.Type,
) (*fvalue.Value, error) {
	r.leaves = append(r.leaves, leaf{
		path:   path,
		goType: fieldType.String(),
	})

	return nil, nil //nolint:nilnil // matches StringBasedBuilder.Get when nothing is found
}

// Compute walks the model and layer builders exactly like dsco.Fill, but
// instead of loading values it returns the canonical key each
// string-based layer would accept for every required leaf field of cfg.
// No environment variables, command-line arguments, or files are read.
//
// cfg must be **T (pointer to a pointer to a struct), mirroring the Fill
// calling convention: var c *MyConfig; Compute(&c, layers...).
//
// Pattern: Factory — assembles a Report from model + layers without I/O.
func Compute(cfg any, layers ...dsco.Layer) (*Report, error) {
	const errCtx = "computing inventory"

	// Dereference **T → *T so PrepareInventoryWalk receives a plain
	// pointer-to-struct, mirroring dscoContext.generateModel in Fill.
	rv := reflect.ValueOf(cfg)
	if rv.Kind() != reflect.Pointer {
		return nil, fmt.Errorf(
			"%s: %w",
			errCtx,
			errors.Join(dsco.ErrFiller, dsco.ErrCfgMustBePointer),
		)
	}

	inner := rv.Elem().Interface()

	walk, err := dsco.PrepareInventoryWalk(inner, layers...)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %w", errCtx, errors.Join(dsco.ErrFiller, err),
		)
	}

	report, err := computeFromWalk(walk)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCtx, err)
	}

	return report, nil
}

// computeFromWalk runs the reporter loop and reduce phase for a prepared
// InventoryWalk. Extracted so in-package tests can inject a walk with a
// synthetic failing reporter without going through the public layer API.
func computeFromWalk(walk *dsco.InventoryWalk) (*Report, error) {
	const errCtx = "computing inventory"

	perLayer := make([]dsco.LayerInventory, 0, len(walk.Reporters))

	for idx, reporter := range walk.Reporters {
		inv, invErr := reporter.ReportInventory(walk.Model)
		if invErr != nil {
			return nil, fmt.Errorf(
				"%s: layer #%d: %w",
				errCtx,
				idx,
				errors.Join(dsco.ErrFiller, invErr),
			)
		}

		perLayer = append(perLayer, inv)
	}

	return reduce(walk.Model, perLayer), nil
}

// reduce collapses per-layer reports into one Field per leaf, applying
// precedence rules: the last string-based layer that can supply a field
// wins for Key; any struct layer that bakes in a value populates
// Satisfied.
func reduce(
	mdl dsco.ModelInterface,
	perLayer []dsco.LayerInventory,
) *Report {
	leaves := collectLeaves(mdl)

	fields := make([]Field, 0, len(leaves))

	for _, lf := range leaves {
		field := Field{Path: lf.path, GoType: lf.goType}

		// Walk layers in declaration order; later ones override earlier
		// for Key (per dsco precedence).
		for _, inv := range perLayer {
			for _, prov := range inv.Provides {
				if prov.FieldUID != lf.path {
					continue
				}

				if prov.Value != nil {
					field.Satisfied = &Satisfaction{
						LayerID: trimStructPrefix(inv.Name),
						Value:   prov.Value,
					}
				}

				if prov.Key != "" {
					field.Key = &KeySpec{
						Layer: layerKindFromName(inv.Name),
						Key:   prov.Key,
					}
				}
			}
		}

		fields = append(fields, field)
	}

	sort.Slice(fields, func(ii, jj int) bool {
		return fields[ii].Path < fields[jj].Path
	})

	return &Report{
		Type:   mdl.TypeName(),
		Fields: fields,
	}
}

// collectLeaves walks mdl and returns one leaf entry per scalar field.
func collectLeaves(mdl dsco.ModelInterface) []leaf {
	rec := &leafRecorder{}
	_, _ = mdl.ApplyOn(rec) //nolint:errcheck // recorder never errors

	return rec.leaves
}

// trimStructPrefix strips the "struct:" prefix from a layer Name to
// expose just the user-supplied id.
func trimStructPrefix(name string) string {
	const prefix = "struct:"
	if len(name) > len(prefix) && name[:len(prefix)] == prefix {
		return name[len(prefix):]
	}

	return name
}

// layerKindFromName extracts the kind (e.g. "env") from a layer Name
// like "env:MYAPP". Returns the whole name if no colon is present
// (e.g. "cmdline").
func layerKindFromName(name string) string {
	if idx := strings.IndexByte(name, ':'); idx >= 0 {
		return name[:idx]
	}

	return name
}
