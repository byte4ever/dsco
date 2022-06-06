package env

import (
	"errors"
)

// ErrInvalidPrefix represents an error when creating the provider with an
// invalid prefix.
var ErrInvalidPrefix = errors.New("invalid prefix")

// ErrInvalidKeyFormat represent an error when a key starts with a valid
// prefix but with invalid syntax.
var ErrInvalidKeyFormat = errors.New("invalid key Format")
