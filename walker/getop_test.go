package walker

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetList_Push(t *testing.T) {
	t.Parallel()

	t.Run(
		"initial state",
		func(t *testing.T) {
			t.Parallel()

			var l GetList

			f1 := func(g Getter) (
				uid uint,
				fieldValue *FieldValue,
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

			f1 := func(g Getter) (
				uid uint,
				fieldValue *FieldValue,
				err error,
			) {
				return
			}

			l := GetList{f1}

			f2 := func(g Getter) (
				uid uint,
				fieldValue *FieldValue,
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

func TestGetList_ApplyOn(t *testing.T) {
	t.Parallel()

	t.Run(
		"initial state",
		func(t *testing.T) {
			t.Parallel()

			getter := NewMockGetter(t)
			f0 := NewMockGetOp(t)
			fv0 := &FieldValue{
				location: "fv0",
			}

			f1 := NewMockGetOp(t)
			f2 := NewMockGetOp(t)
			fv2 := &FieldValue{
				location: "fv2",
			}
			f3 := NewMockGetOp(t)

			f0.On("Execute", getter).Return(
				uint(0), fv0, nil,
			)

			f1.On("Execute", getter).Return(
				uint(1), nil, errMocked1,
			)

			f2.On("Execute", getter).Return(
				uint(2), fv2, nil,
			)

			f3.On("Execute", getter).Return(
				uint(3), nil, errMocked2,
			)

			l := GetList{
				func(g Getter) (
					uid uint,
					fieldValue *FieldValue,
					err error,
				) {
					return f0.Execute(g)
				},
				func(g Getter) (
					uid uint,
					fieldValue *FieldValue,
					err error,
				) {
					return f1.Execute(g)
				},
				func(g Getter) (
					uid uint,
					fieldValue *FieldValue,
					err error,
				) {
					return f2.Execute(g)
				},
				func(g Getter) (
					uid uint,
					fieldValue *FieldValue,
					err error,
				) {
					return f3.Execute(g)
				},
			}

			res, errs := l.ApplyOn(getter)
			require.Len(t, errs, 2)
			require.Equal(t, errs[0], errMocked1)
			require.Equal(t, errs[1], errMocked2)
			require.Equal(
				t, res, FieldValues{
					0: fv0,
					2: fv2,
				},
			)

		},
	)
}
