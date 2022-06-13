package sbased2

import (
	"fmt"
)

// Use implements the dsco.Binder2 interface.
func (s *Binder) Use(
	key string,
) error {
	const errFmt = "%s : %w"

	entry := s.values[key]

	if entry == nil {
		return fmt.Errorf(errFmt, key, ErrKeyNotFound)
	}

	if entry.state != unused {
		return fmt.Errorf(errFmt, entry.location, ErrNotUnused)
	}

	entry.state = used

	return nil
}
