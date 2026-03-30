package utils

import (
	"reflect"
	"strings"
)

// GetKeyName returns key name based on root prefix and struct field.
// It first attempts to extract the key from YAML tags, then falls back to
// converting the Go field name to snake_case. If a prefix is provided,
// it's prepended with a hyphen separator.
func GetKeyName(prefix string, fieldType reflect.StructField) (keyName string) {
	// Extract field name from YAML tag or use field name conversion
	fn := fieldName(fieldType)

	// If no YAML tag found, convert Go field name to snake_case
	if fn == "" {
		fn = ToSnakeCase(fieldType.Name)
	}

	// Return plain field name if no prefix specified
	if prefix == "" {
		return fn
	}

	// Build prefixed key name: "prefix-fieldname"
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteRune('-') // Use hyphen as separator
	sb.WriteString(fn)

	return sb.String()
}

// fieldName extracts the key name from a struct field's YAML tag.
// It handles YAML tags like "key,omitempty" by taking only the key part.
func fieldName(fieldType reflect.StructField) string {
	return strings.Split(
		// Remove spaces from YAML tag to handle malformed tags
		strings.ReplaceAll(
			fieldType.Tag.Get("yaml"), // Get the YAML struct tag
			" ",                       // Remove any spaces
			"",
		),
		",", // Split on comma to separate key from options like "omitempty"
	)[0] // Take only the first part (the key name)
}
