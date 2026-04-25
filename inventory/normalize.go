package inventory

import (
	"fmt"

	gojson "github.com/goccy/go-json"
	goyaml "github.com/goccy/go-yaml"
)

type (
	// satisfactionJSON is a helper struct for JSON marshaling of Satisfaction.
	// It mirrors Satisfaction but holds the normalized value.
	satisfactionJSON struct {
		Value   any    `json:"value"`
		LayerID string `json:"layer_id"`
	}

	// satisfactionYAML is a helper struct for YAML marshaling of Satisfaction.
	// It mirrors Satisfaction but holds the normalized value.
	satisfactionYAML struct {
		Value   any    `yaml:"value"`
		LayerID string `yaml:"layer_id"`
	}

	// fieldJSON is a helper struct for JSON marshaling of Field.
	// Fields are emitted in human-readable order: path, go_type, satisfied, key.
	// Field order is intentional for output readability; fieldalignment is
	// secondary to serialization contract.
	//nolint:govet // fieldalignment: output field order takes priority over struct padding
	fieldJSON struct {
		Path      string        `json:"path"`
		GoType    string        `json:"go_type"`
		Satisfied *Satisfaction `json:"satisfied,omitempty"`
		Key       *KeySpec      `json:"key,omitempty"`
	}

	// fieldYAML is a helper struct for YAML marshaling of Field.
	// Fields are emitted in human-readable order: path, go_type, satisfied, key.
	// Field order is intentional for output readability; fieldalignment is
	// secondary to serialization contract.
	//nolint:govet // fieldalignment: output field order takes priority over struct padding
	fieldYAML struct {
		Path      string        `yaml:"path"`
		GoType    string        `yaml:"go_type"`
		Satisfied *Satisfaction `yaml:"satisfied,omitempty"`
		Key       *KeySpec      `yaml:"key,omitempty"`
	}
)

// Compile-time interface assertions.
var (
	_ gojson.Marshaler          = Satisfaction{}
	_ goyaml.InterfaceMarshaler = Satisfaction{}
	_ gojson.Marshaler          = Field{}
	_ goyaml.InterfaceMarshaler = Field{}
)

// normalizeValue converts fmt.Stringer values (time.Duration, time.Time,
// *url.URL, …) to their String() form so JSON / YAML / text output stays
// human-readable. Primitives and plain structs pass through unchanged.
func normalizeValue(val any) any {
	if val == nil {
		return nil
	}

	if stringer, ok := val.(fmt.Stringer); ok {
		return stringer.String()
	}

	return val
}

// MarshalJSON implements json.Marshaler so Satisfaction.Value is
// normalized before serialization.
func (s Satisfaction) MarshalJSON() ([]byte, error) {
	raw, err := gojson.Marshal(satisfactionJSON{
		LayerID: s.LayerID,
		Value:   normalizeValue(s.Value),
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling satisfaction: %w", err)
	}

	return raw, nil
}

// MarshalYAML implements yaml.InterfaceMarshaler so Satisfaction.Value
// is normalized before serialization.
func (s Satisfaction) MarshalYAML() (any, error) {
	return satisfactionYAML{
		LayerID: s.LayerID,
		Value:   normalizeValue(s.Value),
	}, nil
}

// MarshalJSON implements json.Marshaler so Field keys are emitted in
// human-readable order: path, go_type, satisfied, key.
func (f Field) MarshalJSON() ([]byte, error) {
	raw, err := gojson.Marshal(fieldJSON{
		Path:      f.Path,
		GoType:    f.GoType,
		Satisfied: f.Satisfied,
		Key:       f.Key,
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling field: %w", err)
	}

	return raw, nil
}

// MarshalYAML implements yaml.InterfaceMarshaler so Field keys are emitted in
// human-readable order: path, go_type, satisfied, key.
func (f Field) MarshalYAML() (any, error) {
	return fieldYAML{
		Path:      f.Path,
		GoType:    f.GoType,
		Satisfied: f.Satisfied,
		Key:       f.Key,
	}, nil
}
