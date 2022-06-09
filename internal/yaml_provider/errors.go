package yaml_provider

import (
	"errors"
)

// ErrInvalidModel represents error when input model is not valid.
var ErrInvalidModel = errors.New("invalid model")

// ErrNilReaderFunctor represent an error where read functor is not valid.
var ErrNilReaderFunctor = errors.New("nil read functor")
