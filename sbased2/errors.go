package sbased2

import (
	"errors"
	"fmt"
)

// ErrParse represents an error indicating that a value cannot be parsed.
var ErrParse = errors.New("parse error")

// ErrAliasCollision represents an error indicating that an alias is colliding
// with an actual key in the structure.
var ErrAliasCollision = errors.New("alias collision")

// ErrUnboundKey represents an error indicating that a key is never bound to the
// structure.
var ErrUnboundKey = errors.New("unbound key")

// ErrOverriddenKey represents an error indicating that a potential key binding
// wha overridden in another layer.
var ErrOverriddenKey = errors.New("overridden key")

var ErrNilProvider = errors.New("nil provider")

// Errors returns all errors encountered during processing of the
// layer.
func (s *Binder) Errors() []error {
	var errs []error

	const errFormat = "%s: %w"

	for _, entry := range s.values {
		if !entry.bounded {
			errs = append(
				errs, fmt.Errorf(
					errFormat,
					entry.location,
					ErrUnboundKey,
				),
			)

			continue
		}

		if !entry.used {
			errs = append(
				errs,
				fmt.Errorf(
					errFormat,
					entry.location,
					ErrOverriddenKey,
				),
			)

			continue
		}
	}

	return errs
}
