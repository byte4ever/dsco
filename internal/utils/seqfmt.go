package utils

import (
	"fmt"
	"strings"
)

const (
	noSeqPanicMsg = "no sequence to format"
)

// FormatIndexSequence formats nicely slice of integers.
func FormatIndexSequence(indexes []int) string { //nolint:dupl // it's ok
	const (
		singleDigit = "#%d"
	)

	l := make([]string, len(indexes))
	for i, idx := range indexes {
		l[i] = fmt.Sprintf(singleDigit, idx)
	}

	return formatSequence(l)
}

// FormatStringSequence formats nicely slice of strings.
func FormatStringSequence(values []string) string { //nolint:dupl // it's ok
	const (
		singleString = "%q"
	)

	l := make([]string, len(values))
	for i, idx := range values {
		l[i] = fmt.Sprintf(singleString, idx)
	}

	return formatSequence(l)
}

// formatSequence formats nicely slice of strings.
func formatSequence(values []string) string { //nolint:dupl // it's ok
	const (
		comaSingleString = `, %s`
		andSingleString  = ` and %s`
	)

	indexesLen := len(values)

	if indexesLen == 0 {
		panic(noSeqPanicMsg)
	}

	var sb strings.Builder

	sb.WriteString(values[0])

	if indexesLen == 1 {
		return sb.String()
	}

	if indexesLen > 2 {
		for _, idx := range values[1 : indexesLen-1] {
			sb.WriteString(fmt.Sprintf(comaSingleString, idx))
		}
	}

	sb.WriteString(fmt.Sprintf(andSingleString, values[indexesLen-1]))

	return sb.String()
}
