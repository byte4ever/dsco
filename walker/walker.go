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

func (e elems) Len() int {
	return len(e)
}

func (e elems) Less(i, j int) bool {
	return e[i].order < e[j].order
}

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

func (w *walker) walk(id *int, curPath string, value reflect.Value) error {
	t := value.Type()

	if t.Kind() != reflect.Pointer {
		return errors.New("expect pointer")
	}

	te := t.Elem()
	if te.Kind() != reflect.Struct {
		return errors.New("expect pointer on struct")
	}

	var st stack
	{
	}

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
				return errors.New(
					fmt.Sprintf(
						"colision %q %q",
						concatKey(toProcess.path, toProcess.field.Name),
						concatKey(prevDecl.path, prevDecl.field.Name),
					),
				)
			}

			processed[toProcess.field.Name] = toProcess
		}
	}

	ll := make(elems, 0, len(processed))
	for _, e := range processed {
		ll = append(ll, e)
	}

	sort.Sort(ll)

	for _, e := range ll {
		if e.value.Kind() == reflect.Slice ||
			dsco.TypeIsRegistered(e.value.Type()) {
			ck := concatKey(curPath, e.field.Name)

			if w.scanMode {
				if !e.value.IsNil() {
					if err := w.walkFunc(*id, ck, &e.value); err != nil {
						return err
					}
				}
				*id++
				continue
			}

			if !e.value.IsNil() {
				return errors.New(
					fmt.Sprintf(
						"detecting non nil values %q",
						ck,
					),
				)
			}

			if err := w.walkFunc(
				*id, ck, &e.value,
			); err != nil {
				return err
			}

			*id++

			continue
		}

		if e.value.Kind() == reflect.Pointer &&
			e.value.Type().Elem().Kind() == reflect.Struct {
			ck := concatKey(
				curPath,
				e.field.Name,
			)

			if w.scanMode {
				if !e.value.IsNil() {
					if err := w.walk(
						id,
						ck,
						e.value,
					); err != nil {
						return err
					}
				}

				continue
			}

			if !e.value.IsNil() {
				return errors.New(
					fmt.Sprintf(
						"detecting non nil values %q",
						ck,
					),
				)
			}

			n := reflect.New(e.value.Type().Elem())
			e.value.Set(n)

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
