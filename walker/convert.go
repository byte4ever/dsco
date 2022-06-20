package walker

import (
	"strings"

	"github.com/byte4ever/dsco/utils"
)

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
