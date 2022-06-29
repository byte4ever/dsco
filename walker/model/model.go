package model

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/merror"
	"github.com/byte4ever/dsco/walker/fvalues"
	"github.com/byte4ever/dsco/walker/ifaces"
	"github.com/byte4ever/dsco/walker/plocation"
)

type Model struct {
	accelerator Node
	typeName    string
	getList     GetList
}

func (m *Model) TypeName() string {
	return m.typeName
}

type ModelError struct {
	merror.MError
}

var ErrModel = errors.New("")

func (e ModelError) Is(err error) bool {
	return errors.Is(err, ErrModel)
}

func NewModel(inputModelType reflect.Type) (*Model, error) {
	var maxUID uint

	accelerator, errs := scan(&maxUID, "", inputModelType)

	if !errs.None() {
		return nil, errs
	}

	var getList GetList
	accelerator.BuildGetList(&getList)

	return &Model{
		typeName:    dsco.LongTypeName(inputModelType),
		accelerator: accelerator,
		getList:     getList,
	}, nil
}

func (m *Model) ApplyOn(g ifaces.Getter) (fvalues.FieldValues, error) {
	return m.getList.ApplyOn(g)
}

func (m *Model) FeedFieldValues(id string, v reflect.Value) fvalues.FieldValues {
	k := make(fvalues.FieldValues, len(m.getList))

	m.accelerator.FeedFieldValues(
		id,
		k,
		v,
	)

	return k
}

func (m *Model) Fill(
	inputModelValue reflect.Value, layers []fvalues.FieldValues,
) (plocation.PathLocations, error) {
	return m.accelerator.Fill(
		inputModelValue,
		layers,
	)
}

func scan(uid *uint, path string, t reflect.Type) (Node, ModelError) {
	var errs ModelError

	switch {
	case t.Kind() == reflect.Slice || dsco.TypeIsRegistered(t):
		n := &ValueNode{
			UID:         *uid,
			Type:        t,
			VisiblePath: path,
		}
		*uid++

		return n, errs

	case t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct:
		node := &StructNode{
			Type: t,
		}

		visibleFields, lErrs := getVisibleFieldList(path, t)
		if len(lErrs) > 0 {
			errs.MError = append(errs.MError, lErrs...)
		}

		for _, field := range visibleFields {
			subNode, subErrs := scan(
				uid, pathTo(
					path,
					field.field.Name,
				), field.field.Type,
			)

			if !subErrs.None() {
				errs.MError = append(errs.MError, lErrs...)
				continue
			}

			node.PushSubNodes(field.index, subNode)
		}

		return node, errs
	default:
		return nil, ModelError{
			MError: []error{
				fmt.Errorf("%s: %w", dsco.LongTypeName(t),
					ErrUnsupportedType,
				),
			},
		}
	}
}

func (s *stackEmbed) pushToStack(
	index []int, depth int, path string, _type reflect.Type,
) error {
	if _type.Kind() != reflect.Struct {
		return fmt.Errorf("%s: %w", path, ErrInvalidEmbedded)
	}

	for i := _type.NumField() - 1; i >= 0; i-- {
		field := _type.Field(i)

		ni := make([]int, len(index)+1)
		copy(ni, index)
		ni[len(index)] = i

		s.push(
			&elemEmbedded{
				index: ni,
				depth: depth,
				field: field,
				path:  path,
			},
		)
	}

	return nil
}

func getVisibleFieldList(path string, t reflect.Type) (elems, []error) {
	var errs []error

	st := make(stackEmbed, 0, 16)

	_ = st.pushToStack(nil, 0, "", t.Elem())

	var (
		order int
	)

	processed := make(map[string]*elemEmbedded)

	for st.more() {
		toProcess := st.pop()

		if !toProcess.field.IsExported() {
			continue
		}

		localFieldName := pathTo(
			toProcess.path,
			toProcess.field.Name,
		)

		// deal with embedded structs
		if toProcess.field.Anonymous {
			// pay attention to error to detect embedded pointer structs
			err := st.pushToStack(
				toProcess.index,
				toProcess.depth+1,
				localFieldName,
				toProcess.field.Type,
			)

			if err != nil {
				errs = append(errs, err)
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
				errs = append(
					errs,
					fmt.Errorf(
						"%q %q: %w",
						pathTo(
							path,
							localFieldName,
						),
						pathTo(
							path,
							pathTo(
								prevDecl.path,
								prevDecl.field.Name,
							),
						),
						ErrFieldNameCollision,
					),
				)

				continue
			}

			processed[toProcess.field.Name] = toProcess
		}
	}

	// reorder processed fields
	fieldValues := make(elems, 0, len(processed))
	for _, e := range processed {
		fieldValues = append(fieldValues, e)
	}

	sort.Slice(
		fieldValues, func(i, j int) bool {
			return fieldValues[i].order < fieldValues[j].order
		},
	)

	return fieldValues, errs
}
