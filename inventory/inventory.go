package inventory

type (
	// Report is the static inventory of a config struct against a layer set.
	Report struct {
		Type   string  `json:"type"   yaml:"type"`
		Fields []Field `json:"fields" yaml:"fields"`
	}

	// Field describes the canonical key (and any baked-in default) for one
	// leaf field of the config struct.
	Field struct {
		Path      string        `json:"path"                yaml:"path"`
		GoType    string        `json:"go_type"             yaml:"go_type"`
		Satisfied *Satisfaction `json:"satisfied,omitempty" yaml:"satisfied,omitempty"`
		Key       *KeySpec      `json:"key,omitempty"       yaml:"key,omitempty"`
	}

	// Satisfaction records that a struct layer already provides a value
	// for this field.
	Satisfaction struct {
		LayerID string `json:"layer_id" yaml:"layer_id"`
		Value   any    `json:"value"    yaml:"value"`
	}

	// KeySpec is the canonical (highest-precedence) key form a string-based
	// layer would accept to supply this field.
	KeySpec struct {
		Layer string `json:"layer" yaml:"layer"`
		Key   string `json:"key"   yaml:"key"`
	}
)
