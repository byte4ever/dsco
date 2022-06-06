package cmdline

import (
	"errors"
)

// ErrParamFormat represents an error when creating the provider with an
// invalid command line option.
var ErrParamFormat = errors.New("options param not in --xxx=val format")

// ErrDuplicateParam represents an error when creating the provider with
// duplicated options in the command line.
var ErrDuplicateParam = errors.New("duplicate param")
