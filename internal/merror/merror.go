package merror

import (
	"errors"
	"strings"
)

var Err = errors.New("")

// TODO :- lmartin 7/5/22 -: rename

type MError []error //nolint:errname // need to be fixed

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

func (MError) Is(err error) bool {
	return errors.Is(err, Err)
}

func (m MError) As(errAs any) bool {
	for _, err := range m {
		//goland:noinspection GoErrorsAs
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

func (m *MError) Count() int {
	return len(*m)
}
