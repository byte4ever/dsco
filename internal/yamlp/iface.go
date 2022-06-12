package yamlp

import (
	"io"
)

// ReaderFunctor enable applying action using a reader.
type ReaderFunctor interface {
	// Apply applies action using the reader.
	Apply(action func(r io.Reader) error) error
}
