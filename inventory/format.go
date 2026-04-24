package inventory

import (
	"fmt"
	"io"

	"github.com/goccy/go-json"
	"github.com/goccy/go-yaml"
)

// WriteJSON writes the inventory as JSON via github.com/goccy/go-json.
// Indentation is two spaces; output ends with a trailing newline.
func (r *Report) WriteJSON(writer io.Writer) error {
	const errCtx = "writing JSON inventory"

	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	if _, err = writer.Write(out); err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	if _, err = writer.Write([]byte{'\n'}); err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	return nil
}

// WriteYAML writes the inventory as YAML via github.com/goccy/go-yaml.
// Output ends with a trailing newline.
func (r *Report) WriteYAML(writer io.Writer) error {
	const errCtx = "writing YAML inventory"

	out, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	if _, err = writer.Write(out); err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	return nil
}
