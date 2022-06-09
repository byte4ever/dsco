package utils

import (
	"reflect"
	"strings"
)

// GetKeyName returns key name based on root prefix and struct field.
func GetKeyName(prefix string, fieldType reflect.StructField) (keyName string) {
	fn := fieldName(fieldType)

	if fn == "" {
		fn = toSnakeCase(fieldType.Name)
	}

	if prefix == "" {
		return fn
	}

	var sb strings.Builder

	sb.WriteString(prefix)
	sb.WriteRune('-')
	sb.WriteString(fn)

	return sb.String()
}

func fieldName(fieldType reflect.StructField) string {
	return strings.Split(
		strings.ReplaceAll(
			fieldType.Tag.Get("yaml"),
			" ",
			"",
		),
		",",
	)[0]
}
