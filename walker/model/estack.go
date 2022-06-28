package model

import (
	"reflect"
)

type elemEmbedded struct {
	path  string
	index []int
	field reflect.StructField
	depth int
	order int
}

type stackEmbed []*elemEmbedded

type elems []*elemEmbedded

func (s *stackEmbed) push(e *elemEmbedded) {
	*s = append(*s, e)
}

func (s *stackEmbed) pop() *elemEmbedded {
	ls := len(*s) - 1
	ee := (*s)[ls]
	*s = (*s)[:ls]

	return ee
}

func (s *stackEmbed) more() bool {
	return len(*s) > 0
}
