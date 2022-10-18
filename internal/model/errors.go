package model

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/registry"
)

type FieldNameCollisionError struct {
	Path1 string
	Path2 string
}

func (u FieldNameCollisionError) Error() string {
	return fmt.Sprintf(
		"field collision between %s and %s",
		u.Path1,
		u.Path2,
	)
}

type UnsupportedTypeError struct {
	Type reflect.Type
	Path string
}

func (u UnsupportedTypeError) Error() string {
	return fmt.Sprintf(
		"struct field %s with unsupported type %s",
		u.Path,
		registry.LongTypeName(u.Type),
	)
}

// ErrInvalidEmbedded represent an error where ....
var ErrInvalidEmbedded = errors.New("invalid embedded")

// ErrUninitializedKey represent an error where ....
var ErrUninitializedKey = errors.New("uninitialized key")
