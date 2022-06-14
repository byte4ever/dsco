package utils

import (
	"fmt"
	"strings"
)

const noSeqPanicMsg = "no sequence to format"

// FormatIndexSequence formats nicely slice of integers.
func FormatIndexSequence(indexes []int) string { //nolint:dupl // it's ok
	const (
		single     = "#%d"
		comaSingle = ", " + single
		andSingle  = " and " + single
	)

	indexesLen := len(indexes)

	if indexesLen == 0 {
		panic(noSeqPanicMsg)
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

// FormatStringSequence formats nicely slice of strings.
func FormatStringSequence(values []string) string { //nolint:dupl // it's ok
	const (
		single     = `"%s"`
		comaSingle = ", " + single
		andSingle  = " and " + single
	)

	indexesLen := len(values)

	if indexesLen == 0 {
		panic(noSeqPanicMsg)
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(single, values[0]))

	if indexesLen == 1 {
		return sb.String()
	}

	if indexesLen > 2 {
		for _, idx := range values[1 : indexesLen-1] {
			sb.WriteString(fmt.Sprintf(comaSingle, idx))
		}
	}

	sb.WriteString(fmt.Sprintf(andSingle, values[indexesLen-1]))

	return sb.String()
}
