package dsco

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/merror"
	"github.com/byte4ever/dsco/utils"
)

// ErrNilInput is dummy...
var ErrNilInput = errors.New("nil input")

// ErrCmdlineAlreadyUsed represent an error where ....
var ErrInvalidInput = errors.New("")

type InvalidInputError struct {
	Type reflect.Type
}

// ErrCmdlineAlreadyUsed represent an error where ....
var ErrCmdlineAlreadyUsed = errors.New("")

type CmdlineAlreadyUsedError struct {
	Index int
}

// ErrDuplicateEnvPrefix represent an error where ....
var ErrDuplicateEnvPrefix = errors.New("")

type DuplicateEnvPrefixError struct {
	Index  int
	Prefix string
}

// ErrDuplicateInputStruct represent an error where ....
var ErrDuplicateInputStruct = errors.New("")

type DuplicateInputStructError struct {
	Index int
}

// ErrDuplicateStructID represent an error where ....
var ErrDuplicateStructID = errors.New("")

type DuplicateStructIDError struct {
	Index int
	ID    string
}

func (c InvalidInputError) Error() string {
	return fmt.Sprintf(
		"type %s is not a valid pointer on struct",
		utils.LongTypeName(c.Type),
	)
}

func (InvalidInputError) Is(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

func (c CmdlineAlreadyUsedError) Error() string {
	return fmt.Sprintf(
		"cmdline already used in position #%d",
		c.Index,
	)
}

func (c CmdlineAlreadyUsedError) Is(err error) bool {
	return errors.Is(err, ErrCmdlineAlreadyUsed)
}

func (c DuplicateEnvPrefixError) Error() string {
	return fmt.Sprintf(
		"layer #%d has same prefix=%s",
		c.Index,
		c.Prefix,
	)
}

func (c DuplicateEnvPrefixError) Is(err error) bool {
	return errors.Is(err, ErrDuplicateEnvPrefix)
}

func (c DuplicateInputStructError) Error() string {
	return fmt.Sprintf(
		"struct layer #%d is using same pointer",
		c.Index,
	)
}

func (c DuplicateInputStructError) Is(err error) bool {
	return errors.Is(err, ErrDuplicateInputStruct)
}

func (c DuplicateStructIDError) Error() string {
	return fmt.Sprintf(
		"struct layer #%d is using same id=%q",
		c.Index,
		c.ID,
	)
}

func (c DuplicateStructIDError) Is(err error) bool {
	return errors.Is(err, ErrDuplicateStructID)
}

type LayerErrors struct {
	merror.MError
}

var ErrLayer = errors.New("")

func (m LayerErrors) Is(err error) bool {
	return errors.Is(err, ErrLayer)
}
