package cmdline

import (
	"errors"
	"fmt"
	"strings"

	"github.com/byte4ever/dsco/internal/utils"
)

// ErrInvalidFormat represents an error when creating the provider with an
// invalid command line option.
var ErrInvalidFormat = errors.New("invalid format")

// ErrDuplicateParam represents an error when creating the provider with
// duplicated options in the command line.
var ErrDuplicateParam = errors.New("duplicate param")

// ParamError represents an error when creating the provider with invalid
// parameters.
type ParamError struct {
	Positions []int
	Errs      []error
}

func (e *ParamError) Error() string {
	lp := len(e.Positions)
	le := len(e.Errs)

	// pre-condition
	if lp != le || lp == 0 {
		panic(
			fmt.Sprintf(
				"malformed ParamError "+
					"#positions=%d and "+
					"#positions=%d",
				lp, le,
			),
		)
	}

	var sb strings.Builder
	if lp == 1 {
		sb.WriteString("cmdline issue at position ")
	} else {
		sb.WriteString("cmdline issues at positions ")
	}

	sb.WriteString(utils.FormatIndexSequence(e.Positions))
	sb.WriteString(": ")

	for i, err := range e.Errs {
		if i != 0 {
			sb.WriteString(" / ")
		}

		sb.WriteString(err.Error())
	}

	return sb.String()
}
