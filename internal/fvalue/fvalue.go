package fvalue

import (
	"reflect"
)

type Value struct {
	Value    reflect.Value
	Location string
	Path     string
}

type Values map[uint]*Value
