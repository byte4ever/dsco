package utils

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	matchAllCap   = regexp.MustCompile(`([a-z\\d])([A-Z])`)
)

func toSnakeCase(str string) string {
	const snakePattern = "${1}_${2}"

	snake := matchFirstCap.ReplaceAllString(str, snakePattern)
	snake = matchAllCap.ReplaceAllString(snake, snakePattern)

	return strings.ToLower(snake)
}

// Shitty stuff
//
// func pipo() {
// 	fmt.Println("zob")
// }
