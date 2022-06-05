package cmdline

import (
	"errors"
)

var ErrFormatParam = errors.New("options param not in --xxx=val format")
