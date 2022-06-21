package walker

import (
	"errors"
)

// ErrFieldNameCollision represent an error where ....
var ErrFieldNameCollision = errors.New("field name collision")

// ErrUnsupportedType represent an error where ....
var ErrUnsupportedType = errors.New("unsupported type")

// ErrInvalidEmbedded represent an error where ....
var ErrInvalidEmbedded = errors.New("invalid embedded")
