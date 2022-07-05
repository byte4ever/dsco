package model

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/byte4ever/dsco/ifaces"
	"github.com/byte4ever/dsco/internal/fvalue"
	"github.com/byte4ever/dsco/internal/merror"
	"github.com/byte4ever/dsco/internal/plocation"
	"github.com/byte4ever/dsco/internal/utils"
)

type Model struct {
	accelerator Node
	getList     GetListInterface
	typeName    string
	fieldCount  uint
}

func (m *Model) TypeName() string {
	return m.typeName
}

type ModelError struct {
	merror.MError
}

var ErrModel = errors.New("")

func (ModelError) Is(err error) bool {
	return errors.Is(err, ErrModel)
}

func NewModel(inputModelType reflect.Type) (*Model, error) {
	var maxUID uint

	accelerator, errs := scan(&maxUID, "", inputModelType)

	if !errs.None() {
		return nil, ModelError{
			MError: errs,
		}
	}

	getList := make(GetList, 0, maxUID)
	accelerator.BuildGetList(&getList)

	return &Model{
		fieldCount:  maxUID,
		typeName:    utils.LongTypeName(inputModelType),
		accelerator: accelerator,
		getList:     &getList,
	}, nil
}

func (m *Model) ApplyOn(g ifaces.Getter) (fvalue.Values, error) {
	return m.getList.ApplyOn(g) //nolint:wrapcheck // don't wrap it
}

func (m *Model) GetFieldValuesFor(
	id string,
	value reflect.Value,
) fvalue.Values {
	fieldValues := make(fvalue.Values, m.fieldCount)

	m.accelerator.FeedFieldValues(
		id,
		fieldValues,
		value,
	)

	return fieldValues
}

func (m *Model) Fill(
	inputModelValue reflect.Value,
	layers []fvalue.Values,
) (plocation.Locations, error) {
	//nolint:wrapcheck // wrap no required
	return m.accelerator.Fill(
		inputModelValue,
		layers,
	)
}

func (s *stackEmbed) pushToStack(
	index []int, depth int, path string, _type reflect.Type,
) error {
	// Todo :- lmartin 7/4/22 -:  create InvalidEmbededError type
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

const preallocateStack = 16

func getVisibleFieldList(path string, t reflect.Type) (elems, []error) {
	var errs []error

	st := make(stackEmbed, 0, preallocateStack)

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
					FieldNameCollisionError{
						Path1: pathTo(
							path,
							localFieldName,
						),
						Path2: pathTo(
							path,
							pathTo(
								prevDecl.path,
								prevDecl.field.Name,
							),
						),
					},
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
