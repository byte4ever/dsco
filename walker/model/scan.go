package model

import (
	"reflect"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/merror"
)

//nolint:ireturn // expected to build abstract tree nodes
func scan(
	uid *uint,
	path string,
	t reflect.Type,
) (Node, merror.MError) {

	switch {
	case t.Kind() == reflect.Slice || dsco.TypeIsRegistered(t):
		n := &ValueNode{
			UID:         *uid,
			Type:        t,
			VisiblePath: path,
		}
		*uid++

		return n, nil

	case t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct:
		var errs merror.MError

		node := &StructNode{
			Type: t,
		}

		visibleFields, lErrs := getVisibleFieldList(path, t)
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
				node.PushSubNodes(field.index, subNode)
			}
		}

		return node, errs
	default:
		return nil, merror.MError{
			UnsupportedTypeError{
				Path: path,
				Type: t,
			},
		}
	}
}
