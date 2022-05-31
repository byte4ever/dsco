package cmdline

import (
	"errors"
)

// ErrFormatParam is dummy.
var ErrFormatParam = errors.New("options param not in --xxx=val format")
