package walker

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/byte4ever/dsco"
)

type elem struct {
	value reflect.Value
	path  string
	field reflect.StructField
	depth int
	order int
}

type stack []*elem

type elems []*elem

type WalkFunc func(
	order int,
	path string,
	value *reflect.Value,
) error

type walker struct {
	walkFunc WalkFunc
	scanMode bool
}

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
) {
	t := value.Type()

	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expect struct got %s", t.Kind().String()))
	}

	for i := t.NumField() - 1; i >= 0; i-- {
		field := t.Field(i)
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
}

// ErrExpectPointerOnStruct represent an error where ....
var ErrExpectPointerOnStruct = errors.New("expect pointer on struct")

// ErrFieldNameCollision represent an error where ....
var ErrFieldNameCollision = errors.New("field name collision")

// ErrNotNilValue represent an error where ....
var ErrNotNilValue = errors.New("not nil value")

func (w *walker) walk(id *int, curPath string, value reflect.Value) error {
	t := value.Type()

	if t.Kind() != reflect.Pointer {
		return ErrExpectPointerOnStruct
	}

	te := t.Elem()
	if te.Kind() != reflect.Struct {
		return ErrExpectPointerOnStruct
	}

	var st stack

	pushToStack(0, "", &st, value.Elem())

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

		if !toProcess.field.IsExported() {
			continue
		}

		if toProcess.field.Anonymous {
			pushToStack(
				toProcess.depth+1,
				concatKey(
					toProcess.path, toProcess.field.Name,
				),
				&st,
				toProcess.value,
			)

			continue
		}

		toProcess.order = order
		order++

		prevDecl, found := processed[toProcess.field.Name]
		if (found && prevDecl.depth >= toProcess.depth) || !found {
			// detecting field collision
			if found {
				return fmt.Errorf(
					"%q %q: %w",
					concatKey(toProcess.path, toProcess.field.Name),
					concatKey(prevDecl.path, prevDecl.field.Name),
					ErrFieldNameCollision,
				)
			}

			processed[toProcess.field.Name] = toProcess
		}
	}

	fieldValues := make(elems, 0, len(processed))
	for _, e := range processed {
		fieldValues = append(fieldValues, e)
	}

	sort.Sort(fieldValues)

	for _, fieldValue := range fieldValues {
		if fieldValue.value.Kind() == reflect.Slice ||
			dsco.TypeIsRegistered(fieldValue.value.Type()) {
			ck := concatKey(curPath, fieldValue.field.Name)

			if w.scanMode {
				if !fieldValue.value.IsNil() {
					if err := w.walkFunc(
						*id, ck, &fieldValue.value,
					); err != nil {
						return err
					}
				}

				*id++

				continue
			}

			if !fieldValue.value.IsNil() {
				return fmt.Errorf(
					"%q: %w",
					ck,
					ErrNotNilValue,
				)
			}

			if err := w.walkFunc(
				*id, ck, &fieldValue.value,
			); err != nil {
				return err
			}

			*id++

			continue
		}

		if fieldValue.value.Kind() == reflect.Pointer &&
			fieldValue.value.Type().Elem().Kind() == reflect.Struct {
			ck := concatKey(
				curPath,
				fieldValue.field.Name,
			)

			if w.scanMode {
				if !fieldValue.value.IsNil() {
					if err := w.walk(
						id,
						ck,
						fieldValue.value,
					); err != nil {
						return err
					}
				}

				continue
			}

			if !fieldValue.value.IsNil() {
				return fmt.Errorf(
					"%q: %w",
					ck,
					ErrNotNilValue,
				)
			}

			n := reflect.New(fieldValue.value.Type().Elem())
			fieldValue.value.Set(n)

			if err := w.walk(
				id,
				ck,
				n,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
