package utils

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	matchAllCap   = regexp.MustCompile(`([a-z\\d])([A-Z])`)
)

// ToSnakeCase converts s to snake case.
func ToSnakeCase(s string) string {
	const snakePattern = "${1}_${2}"

	snake := matchFirstCap.ReplaceAllString(s, snakePattern)
	snake = matchAllCap.ReplaceAllString(snake, snakePattern)

	return strings.ToLower(snake)
}
