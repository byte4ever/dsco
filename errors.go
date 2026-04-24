package dsco

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/byte4ever/dsco/internal/merror"
	"github.com/byte4ever/dsco/registry"
)

// ErrNilInput is dummy...
var ErrNilInput = errors.New("nil input")

// ErrInvalidInput is the sentinel error for invalid input types.
var ErrInvalidInput = errors.New("")

// InvalidInputError represents an error where the input type is invalid.
type InvalidInputError struct {
	Type reflect.Type
}

// ErrCmdlineAlreadyUsed is the sentinel error for cmdline already used.
var ErrCmdlineAlreadyUsed = errors.New("")

// CmdlineAlreadyUsedError represents an error where cmdline is already used.
type CmdlineAlreadyUsedError struct {
	Index int
}

// ErrDuplicateEnvPrefix is the sentinel error for duplicate env prefix.
var ErrDuplicateEnvPrefix = errors.New("")

// DuplicateEnvPrefixError represents an error where env prefix is duplicated.
type DuplicateEnvPrefixError struct {
	Prefix string
	Index  int
}

// ErrDuplicateInputStruct is the sentinel error for duplicate input struct.
var ErrDuplicateInputStruct = errors.New("")

// DuplicateInputStructError represents an error where input struct is
// duplicated.
type DuplicateInputStructError struct {
	Index int
}

// ErrDuplicateStructID is the sentinel error for duplicate struct ID.
var ErrDuplicateStructID = errors.New("")

// DuplicateStructIDError represents an error where struct ID is duplicated.
type DuplicateStructIDError struct {
	ID    string
	Index int
}

// InvalidInputError methods.
func (c InvalidInputError) Error() string {
	return fmt.Sprintf(
		"type %s is not a valid pointer on struct",
		registry.LongTypeName(c.Type),
	)
}

func (InvalidInputError) Is(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// CmdlineAlreadyUsedError methods.
func (c CmdlineAlreadyUsedError) Error() string {
	return fmt.Sprintf(
		"cmdline already used in position #%d",
		c.Index,
	)
}

func (CmdlineAlreadyUsedError) Is(err error) bool {
	return errors.Is(err, ErrCmdlineAlreadyUsed)
}

// DuplicateEnvPrefixError methods.
func (c DuplicateEnvPrefixError) Error() string {
	return fmt.Sprintf(
		"layer #%d has same prefix=%s",
		c.Index,
		c.Prefix,
	)
}

func (DuplicateEnvPrefixError) Is(err error) bool {
	return errors.Is(err, ErrDuplicateEnvPrefix)
}

// DuplicateInputStructError methods.
func (c DuplicateInputStructError) Error() string {
	return fmt.Sprintf(
		"struct layer #%d is using same pointer",
		c.Index,
	)
}

func (DuplicateInputStructError) Is(err error) bool {
	return errors.Is(err, ErrDuplicateInputStruct)
}

// DuplicateStructIDError methods.
func (c DuplicateStructIDError) Error() string {
	return fmt.Sprintf(
		"struct layer #%d is using same id=%q",
		c.Index,
		c.ID,
	)
}

func (DuplicateStructIDError) Is(err error) bool {
	return errors.Is(err, ErrDuplicateStructID)
}

// ErrDuplicateStringProvider is the sentinel error for duplicate string
// provider.
var ErrDuplicateStringProvider = errors.New("duplicate string provider")

// ErrLayer is the sentinel error for layer errors.
var ErrLayer = errors.New("")

// LayerErrors represents multiple layer errors.
type LayerErrors struct {
	merror.MError
}

func (LayerErrors) Is(err error) bool {
	return errors.Is(err, ErrLayer)
}
