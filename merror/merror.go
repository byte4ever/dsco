package merror

import (
	"errors"
	"strings"
)

var Err = errors.New("")

type MError []error

func (m MError) Error() string {
	if len(m) == 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(m[0].Error())

	for _, err := range m[1:] {
		sb.WriteRune('\n')
		sb.WriteString(err.Error())
	}

	return sb.String()
}

func (m MError) Is(err error) bool {
	return errors.Is(err, Err)
}

func (m MError) As(errAs any) bool {
	for _, err := range m {
		if errors.As(err, errAs) {
			return true
		}
	}

	return false
}

func (m *MError) Add(err error) {
	*m = append(*m, err)
}

func (m *MError) None() bool {
	return len(*m) == 0
}
