package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/walker/fvalues"
	"github.com/byte4ever/dsco/walker/ifaces"
)

func TestGetList_Push(t *testing.T) {
	t.Parallel()

	t.Run(
		"initial state",
		func(t *testing.T) {
			t.Parallel()

			var l GetList

			f1 := func(g ifaces.Getter) (
				uid uint,
				fieldValue *fvalues.FieldValue,
				err error,
			) {
				return
			}

			l.Push(
				f1,
			)

			require.Len(t, l, 1)
			require.Equal(
				t, reflect.ValueOf(f1).Pointer(),
				reflect.ValueOf(l[0]).Pointer(),
			)
		},
	)

	t.Run(
		"accum",
		func(t *testing.T) {
			t.Parallel()

			f1 := func(g ifaces.Getter) (
				uid uint,
				fieldValue *fvalues.FieldValue,
				err error,
			) {
				return
			}

			l := GetList{f1}

			f2 := func(g ifaces.Getter) (
				uid uint,
				fieldValue *fvalues.FieldValue,
				err error,
			) {
				return
			}

			l.Push(
				f2,
			)

			require.Len(t, l, 2)
			require.Equal(
				t, reflect.ValueOf(f2).Pointer(),
				reflect.ValueOf(l[1]).Pointer(),
			)
		},
	)

}

type FakeError1 struct {
	error
}

type FakeError2 struct {
	error
}

func TestGetList_ApplyOn(t *testing.T) {
	t.Parallel()

	t.Run(
		"initial state",
		func(t *testing.T) {
			t.Parallel()

			getter := NewMockGetter(t)
			f0 := NewMockGetOp(t)
			fv0 := &fvalues.FieldValue{
				Location: "fv0",
			}

			f1 := NewMockGetOp(t)
			f2 := NewMockGetOp(t)
			fv2 := &fvalues.FieldValue{
				Location: "fv2",
			}
			f3 := NewMockGetOp(t)

			f0.On("Execute", getter).Return(
				uint(0), fv0, nil,
			)

			f1.On("Execute", getter).Return(
				uint(1), nil, FakeError1{errMocked1},
			)

			f2.On("Execute", getter).Return(
				uint(2), fv2, nil,
			)

			f3.On("Execute", getter).Return(
				uint(3), nil, FakeError2{errMocked2},
			)

			l := GetList{
				func(g ifaces.Getter) (
					uid uint,
					fieldValue *fvalues.FieldValue,
					err error,
				) {
					return f0.Execute(g)
				},
				func(g ifaces.Getter) (
					uid uint,
					fieldValue *fvalues.FieldValue,
					err error,
				) {
					return f1.Execute(g)
				},
				func(g ifaces.Getter) (
					uid uint,
					fieldValue *fvalues.FieldValue,
					err error,
				) {
					return f2.Execute(g)
				},
				func(g ifaces.Getter) (
					uid uint,
					fieldValue *fvalues.FieldValue,
					err error,
				) {
					return f3.Execute(g)
				},
			}

			res, err := l.ApplyOn(getter)

			var e1 FakeError1

			require.ErrorAs(t, err, &e1)
			require.Equal(t, errMocked1, e1.error)

			var e2 FakeError2

			require.ErrorAs(t, err, &e2)
			require.Equal(t, errMocked2, e2.error)

			require.Equal(
				t, res, fvalues.FieldValues{
					0: fv0,
					2: fv2,
				},
			)

		},
	)
}
