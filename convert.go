package dsco

import (
	"strings"

	"github.com/byte4ever/dsco/internal/utils"
)

// convert transforms a dot-separated field path into a dash-separated
// snake_case string suitable for environment variable or command line
// argument naming. Each path segment is converted to snake_case and
// segments are joined with dashes.
//
// Example: "field.subField" becomes "field-sub_field"
func convert(s string) string {
	var sb strings.Builder

	for i, s2 := range strings.Split(s, ".") {
		if i != 0 {
			sb.WriteRune('-')
		}

		sb.WriteString(utils.ToSnakeCase(s2))
	}

	return sb.String()
}
