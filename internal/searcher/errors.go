package searcher

import (
	"errors"
)

// ErrConfNotFound represents an error indicating that no configuration file was
// found.
var ErrConfNotFound = errors.New("no configuration file found")

// ErrNoSearchPath represents an error indicating that no search path was
// provided.
var ErrNoSearchPath = errors.New("no search path")
