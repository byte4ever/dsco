package sbased2

import (
	"errors"
)

// ErrParse represents an error indicating that a value cannot be parsed.
var ErrParse = errors.New("parse error")

// ErrAliasCollision represents an error indicating that an alias is colliding
// with an actual key in the structure.
var ErrAliasCollision = errors.New("alias collision")

// ErrUnboundKey represents an error indicating that a key is never bound to the
// structure.
var ErrUnboundKey = errors.New("unbound key")

// ErrOverriddenKey represents am error indicating that a potential key binding
// wha overridden in another layer.
var ErrOverriddenKey = errors.New("overridden key")
