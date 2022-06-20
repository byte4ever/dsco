package walker

import (
	"reflect"
)

// Get returns a pointer to a value and it's localtion.
// If none exists returns nil and "",
func (b Base) Get(id int) (*reflect.Value, string) {
	if e, found := b[id]; found {
		location := e.location
		value := e.value

		delete(b, id)

		return value, location
	}

	return nil, ""
}
