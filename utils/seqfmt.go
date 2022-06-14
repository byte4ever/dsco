package utils

import (
	"fmt"
	"strings"
)

// FormatIndexSequence formats nicely sequences of indexes.
func FormatIndexSequence(indexes []int) string {
	const (
		single     = "#%d"
		comaSingle = ", " + single
		andSingle  = " and " + single
	)

	indexesLen := len(indexes)

	if indexesLen == 0 {
		panic("no sequence to format")
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(single, indexes[0]))

	if indexesLen == 1 {
		return sb.String()
	}

	if indexesLen > 2 {
		for _, idx := range indexes[1 : indexesLen-1] {
			sb.WriteString(fmt.Sprintf(comaSingle, idx))
		}
	}

	sb.WriteString(fmt.Sprintf(andSingle, indexes[indexesLen-1]))

	return sb.String()
}
