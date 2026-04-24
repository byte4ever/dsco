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
)

// Compile-time interface assertions.
var (
	_ gojson.Marshaler          = Satisfaction{}
	_ goyaml.InterfaceMarshaler = Satisfaction{}
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
