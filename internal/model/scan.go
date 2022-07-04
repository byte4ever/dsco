package model

import (
	"reflect"

	"github.com/byte4ever/dsco/internal/merror"
	"github.com/byte4ever/dsco/internal/utils"
)

//nolint:ireturn // expected to build abstract tree nodes
func scan(
	uid *uint,
	path string,
	_type reflect.Type,
) (Node, merror.MError) {
	switch {
	case _type.Kind() == reflect.Slice || utils.TypeIsRegistered(_type):
		valueNode := &ValueNode{
			UID:         *uid,
			Type:        _type,
			VisiblePath: path,
		}
		*uid++

		return valueNode, nil

	case _type.Kind() == reflect.Pointer && _type.Elem().Kind() == reflect.Struct:
		var errs merror.MError

		structNode := &StructNode{
			Type: _type,
		}

		visibleFields, lErrs := getVisibleFieldList(path, _type)
		if len(lErrs) > 0 {
			errs = append(errs, lErrs...)
		}

		for _, field := range visibleFields {
			subNode, subErrs := scan(
				uid, pathTo(
					path,
					field.field.Name,
				), field.field.Type,
			)

			if !subErrs.None() {
				errs = append(errs, subErrs...)
			}

			if subNode != nil {
				structNode.PushSubNodes(field.index, subNode)
			}
		}

		return structNode, errs
	default:
		return nil, merror.MError{
			UnsupportedTypeError{
				Path: path,
				Type: _type,
			},
		}
	}
}
