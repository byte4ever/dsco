package walker

import (
	"strings"
)

func pathTo(currentPath string, fieldName string) string {
	if currentPath == "" {
		return fieldName
	}

	var sb strings.Builder

	sb.WriteString(currentPath)
	sb.WriteRune('.')
	sb.WriteString(fieldName)

	return sb.String()
}
