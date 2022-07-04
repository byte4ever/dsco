package utils

import (
	"fmt"
	"strings"
)

const noSeqPanicMsg = "no sequence to format"

// FormatIndexSequence formats nicely slice of integers.
func FormatIndexSequence(indexes []int) string { //nolint:dupl // it's ok
	const (
		singleDigit     = "#%d"
		comaSingleDigit = ", " + singleDigit
		andSingleDigit  = " and " + singleDigit
	)

	indexesLen := len(indexes)

	if indexesLen == 0 {
		panic(noSeqPanicMsg)
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(singleDigit, indexes[0]))

	if indexesLen == 1 {
		return sb.String()
	}

	if indexesLen > 2 {
		for _, idx := range indexes[1 : indexesLen-1] {
			sb.WriteString(fmt.Sprintf(comaSingleDigit, idx))
		}
	}

	sb.WriteString(fmt.Sprintf(andSingleDigit, indexes[indexesLen-1]))

	return sb.String()
}

// FormatStringSequence formats nicely slice of strings.
func FormatStringSequence(values []string) string { //nolint:dupl // it's ok
	const (
		singleString     = `"%s"`
		comaSingleString = ", " + singleString
		andSingleString  = " and " + singleString
	)

	indexesLen := len(values)

	if indexesLen == 0 {
		panic(noSeqPanicMsg)
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(singleString, values[0]))

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
