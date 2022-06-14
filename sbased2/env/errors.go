package env

import (
	"errors"
)

// ErrInvalidPrefix represents an error when creating the provider with an
// invalid prefix.
var ErrInvalidPrefix = errors.New("invalid prefix")

// ErrAmbiguousKey represent an error when a key starts with a valid
// prefix but with invalid syntax.
var ErrAmbiguousKey = errors.New("is ambiguous")

// ErrAmbiguousKeys represent an error when multiple env keys  starts with a
// valid prefix but with invalid syntax.
var ErrAmbiguousKeys = errors.New("are ambiguous")
