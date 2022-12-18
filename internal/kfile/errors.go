package kfile

import (
	"strings"
)

type pathError struct {
	err  error
	path string
}

type PathErrors []*pathError

func (p PathErrors) Error() string {
	var sb strings.Builder

	for _, ep := range p {
		sb.WriteString(ep.path)
		sb.WriteString(": ")
		sb.WriteString(ep.err.Error())
		sb.WriteRune('\n')
	}

	return sb.String()
}
