package utils

import (
	"reflect"
	"strings"
)

func GetKeyName(rootKey string, fieldType reflect.StructField) string {
	s := getName(fieldType)

	if s == "" {
		s = ToSnakeCase(fieldType.Name)
	}

	key := appendKey(rootKey, s)

	return key
}

func getName(fieldType reflect.StructField) string {
	return strings.Split(
		strings.ReplaceAll(
			fieldType.Tag.Get("yaml"),
			" ",
			"",
		),
		",",
	)[0]
}

func appendKey(a, b string) string {
	if a == "" {
		return b
	}

	return a + "-" + b
}
