package walker

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/byte4ever/dsco"
)

const errFmt = "%s: %w"

// ErrUnsupportedType represent an error where ....
var ErrUnsupportedType = errors.New("unsupported type")

// ErrEmbeddedPointer represent an error where ....
var ErrEmbeddedPointer = errors.New("embedded pointer")

// ErrExpectPointerOnStruct represent an error where ....
var ErrExpectPointerOnStruct = errors.New("expect pointer on struct")

// ErrFieldNameCollision represent an error where ....
var ErrFieldNameCollision = errors.New("field name collision")

// ErrNotNilValue represent an error where ....
var ErrNotNilValue = errors.New("not nil value")

// ErrNilInterface represent an error where ....
var ErrNilInterface = errors.New("nil interface")

type elem struct {
	value reflect.Value
	path  string
	field reflect.StructField
	depth int
	order int
}

type stack []*elem

type elems []*elem

type actionFunc func(
	order int,
	path string,
	value *reflect.Value,
) error

type walker struct {
	fieldAction actionFunc
	isGetter    bool
}

func newSetter(action actionFunc) *walker {
	return &walker{
		fieldAction: action,
	}
}

/* WIP
func newGetter(action actionFunc) *walker {
	return &walker{
		fieldAction: action,
		isGetter:    true,
	}
}
*/

// Len implement sort.Interface.
func (e elems) Len() int {
	return len(e)
}

// Less implements sort.Interface.
func (e elems) Less(i, j int) bool {
	return e[i].order < e[j].order
}

// Swap implements sort.Interface.
func (e elems) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (s *stack) push(e *elem) {
	*s = append(*s, e)
}

func (s *stack) pop() (*elem, bool) {
	if ls := len(*s) - 1; ls >= 0 {
		ee := (*s)[ls]
		*s = (*s)[:ls]

		return ee, true
	}

	return nil, false
}

func concatKey(currentPath string, fieldName string) string {
	if currentPath == "" {
		return fieldName
	}

	var sb strings.Builder

	sb.WriteString(currentPath)
	sb.WriteRune('.')
	sb.WriteString(fieldName)

	return sb.String()
}

func pushToStack(
	depth int,
	path string,
	curStack *stack,
	value reflect.Value,
) error {
	valueType := value.Type()

	if valueType.Kind() != reflect.Struct {
		return ErrEmbeddedPointer
	}

	for i := valueType.NumField() - 1; i >= 0; i-- {
		field := valueType.Field(i)
		vv := value.Field(i)

		curStack.push(
			&elem{
				depth: depth,
				field: field,
				path:  path,
				value: vv,
			},
		)
	}

	return nil
}

func setStruct(toSet any, action actionFunc) error {
	var maxId int

	if toSet == nil {
		return ErrNilInterface
	}

	return newSetter(action).walkRec(
		&maxId,
		"",
		reflect.ValueOf(toSet),
	)
}

/* WIP
func getStruct(toGet any, action actionFunc) error {
	var maxId int

	if toGet == nil {
		return ErrNilInterface
	}

	return newGetter(action).walkRec(
		&maxId,
		"",
		reflect.ValueOf(toGet),
	)
}
*/

func (w *walker) walkRec(
	id *int, curPath string,
	value reflect.Value,
) error {
	// checking value
	t := value.Type()

	if t.Kind() != reflect.Pointer {
		return ErrExpectPointerOnStruct
	}

	te := t.Elem()
	if te.Kind() != reflect.Struct {
		return ErrExpectPointerOnStruct
	}

	// flatten all fields to get field visiblity when struct are embedded and
	// to check field name collisions
	var st stack

	_ = pushToStack(0, "", &st, value.Elem())

	var (
		order     int
		some      bool
		toProcess *elem
	)

	processed := make(map[string]*elem)

	for {
		toProcess, some = st.pop()

		if !some {
			break
		}

		// don't process unexported fields
		if !toProcess.field.IsExported() {
			continue
		}

		localFieldName := concatKey(
			toProcess.path,
			toProcess.field.Name,
		)

		// deal with embedded structs
		if toProcess.field.Anonymous {
			// pay attention to error to detect embedded pointer structs
			err := pushToStack(
				toProcess.depth+1,
				localFieldName,
				&st,
				toProcess.value,
			)

			if err != nil {
				return fmt.Errorf("%s.%s: %w", curPath, localFieldName, err)
			}

			continue
		}

		toProcess.order = order
		order++

		// filter field visiblity
		prevDecl, found := processed[toProcess.field.Name]
		if (found && prevDecl.depth >= toProcess.depth) || !found {
			// detecting field collision
			if found {
				return fmt.Errorf(
					"%q %q: %w",
					concatKey(
						curPath,
						localFieldName,
					),
					concatKey(
						curPath,
						concatKey(
							prevDecl.path,
							prevDecl.field.Name,
						),
					),
					ErrFieldNameCollision,
				)
			}

			processed[toProcess.field.Name] = toProcess
		}
	}

	// reorder processed fields
	fieldValues := make(elems, 0, len(processed))
	for _, e := range processed {
		fieldValues = append(fieldValues, e)
	}

	sort.Sort(fieldValues)

	// now we can continue the walk on visible fields.
	for _, fieldValue := range fieldValues {
		ck := concatKey(
			curPath,
			fieldValue.field.Name,
		)

		// manage pointer on struct case
		if fieldValue.value.Kind() == reflect.Pointer &&
			fieldValue.value.Type().Elem().Kind() == reflect.Struct {
			// getter behaviour
			if w.isGetter {
				if !fieldValue.value.IsNil() {
					if err := w.walkRec(
						id,
						ck,
						fieldValue.value,
					); err != nil {
						return err
					}
				}

				continue
			}

			// setter behaviour
			if !fieldValue.value.IsNil() {
				return fmt.Errorf(
					errFmt,
					ck,
					ErrNotNilValue,
				)
			}

			n := reflect.New(fieldValue.value.Type().Elem())
			fieldValue.value.Set(n)

			if err := w.walkRec(
				id,
				ck,
				n,
			); err != nil {
				return err
			}

			continue
		}

		// manage slice case and registered types
		if dsco.TypeIsRegistered(fieldValue.value.Type()) || fieldValue.value.
			Kind() == reflect.Slice {
			// getter behaviour
			if w.isGetter {
				if !fieldValue.value.IsNil() {
					if err := w.fieldAction(
						*id, ck, &fieldValue.value,
					); err != nil {
						return err
					}
				}

				*id++

				continue
			}

			// setter behaviour
			if !fieldValue.value.IsNil() {
				return fmt.Errorf(
					errFmt,
					ck,
					ErrNotNilValue,
				)
			}

			if err := w.fieldAction(
				*id, ck, &fieldValue.value,
			); err != nil {
				return err
			}

			*id++

			continue
		}

		// type is not supported
		return fmt.Errorf(
			"%s-<%s>: %w",
			ck,
			fieldValue.value.Type().String(),
			ErrUnsupportedType,
		)
	}

	return nil
}
