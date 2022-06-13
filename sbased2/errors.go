package sbased2

import (
	"errors"
	"fmt"
	"sort"
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

// ErrNilProvider is shitty...
var ErrNilProvider = errors.New("nil provider")

// Errors returns all errors encountered during processing of the
// layer.
func (s *Binder) Errors() []error {
	const errFormat = "%s: %w"

	if len(s.values) < 1 {
		return nil
	}

	var errs []error

	sortedKeys := make([]string, 0, len(s.values))

	for key := range s.values {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		entry := s.values[key]
		switch entry.state {
		case unbounded:
			errs = append(
				errs, fmt.Errorf(
					errFormat,
					entry.location,
					ErrUnboundKey,
				),
			)

		case unused:
			errs = append(
				errs,
				fmt.Errorf(
					errFormat,
					entry.location,
					ErrOverriddenKey,
				),
			)

		case used:
		}
	}

	return errs
}
