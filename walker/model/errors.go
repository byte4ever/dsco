package model

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco"
)

// ErrFieldNameCollision represent an error where ....
var ErrFieldNameCollision = errors.New("field name collision")

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

// ErrUnsupportedType represent an error where ....
var ErrUnsupportedType = errors.New("unsupported type")

type UnsupportedTypeError struct {
	Path string
	Type reflect.Type
}

func (u UnsupportedTypeError) Error() string {
	return fmt.Sprintf(
		"struct field %s with unsupported type %s",
		u.Path,
		dsco.LongTypeName(u.Type),
	)
}

// ErrInvalidEmbedded represent an error where ....
var ErrInvalidEmbedded = errors.New("invalid embedded")

// ErrUninitializedKey represent an error where ....
var ErrUninitializedKey = errors.New("uninitialized key")
