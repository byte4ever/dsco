package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/internal"
	"github.com/byte4ever/dsco/internal/fvalue"
)

func TestGetList_Push(t *testing.T) {
	t.Parallel()

	t.Run(
		"initial state",
		func(t *testing.T) {
			t.Parallel()

			var l GetList

			f1 := func(g internal.Getter) (
				uid uint,
				fieldValue *fvalue.Value,
				err error,
			) {
				//nolint:revive // don't care
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

			f1 := func(g internal.Getter) (
				uid uint,
				fieldValue *fvalue.Value,
				err error,
			) {
				//nolint:revive // don't care
				return
			}

			l := GetList{f1}

			f2 := func(g internal.Getter) (
				uid uint,
				fieldValue *fvalue.Value,
				err error,
			) {
				//nolint:revive // don't care
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
		"success with errors collected",
		func(t *testing.T) {
			t.Parallel()

			getter := NewMockGetter(t)
			f0 := NewMockGetOp(t)
			fv0 := &fvalue.Value{
				Location: "fv0",
			}

			f1 := NewMockGetOp(t)
			f2 := NewMockGetOp(t)
			fv2 := &fvalue.Value{
				Location: "fv2",
			}
			f3 := NewMockGetOp(t)

			f0.On("Execute", getter).Return(
				uint(0), fv0, nil,
			)

			f1.On("Execute", getter).Return(
				uint(1), nil, Mocked1Error{},
			)

			f2.On("Execute", getter).Return(
				uint(2), fv2, nil,
			)

			f3.On("Execute", getter).Return(
				uint(3), nil, Mocked2Error{},
			)

			l := GetList{
				func(g internal.Getter) (
					uid uint,
					fieldValue *fvalue.Value,
					err error,
				) {
					return f0.Execute(g)
				},
				func(g internal.Getter) (
					uid uint,
					fieldValue *fvalue.Value,
					err error,
				) {
					return f1.Execute(g)
				},
				func(g internal.Getter) (
					uid uint,
					fieldValue *fvalue.Value,
					err error,
				) {
					return f2.Execute(g)
				},
				func(g internal.Getter) (
					uid uint,
					fieldValue *fvalue.Value,
					err error,
				) {
					return f3.Execute(g)
				},
			}

			res, err := l.ApplyOn(getter)

			checkAsMockedError1(t, err)
			checkAsMockedError2(t, err)

			require.Equal(
				t, res, fvalue.Values{
					0: fv0,
					2: fv2,
				},
			)

		},
	)

	t.Run(
		"success with no errors",
		func(t *testing.T) {
			t.Parallel()

			getter := NewMockGetter(t)

			f0 := NewMockGetOp(t)
			fv0 := &fvalue.Value{
				Location: "fv0",
			}

			f1 := NewMockGetOp(t)
			fv1 := &fvalue.Value{
				Location: "fv1",
			}

			f2 := NewMockGetOp(t)
			fv2 := &fvalue.Value{
				Location: "fv2",
			}

			f3 := NewMockGetOp(t)
			fv3 := &fvalue.Value{
				Location: "fv3",
			}

			f0.On("Execute", getter).Return(
				uint(0), fv0, nil,
			)

			f1.On("Execute", getter).Return(
				uint(1), fv1, nil,
			)

			f2.On("Execute", getter).Return(
				uint(2), fv2, nil,
			)

			f3.On("Execute", getter).Return(
				uint(3), fv3, nil,
			)

			l := GetList{
				func(g internal.Getter) (
					uid uint,
					fieldValue *fvalue.Value,
					err error,
				) {
					return f0.Execute(g)
				},
				func(g internal.Getter) (
					uid uint,
					fieldValue *fvalue.Value,
					err error,
				) {
					return f1.Execute(g)
				},
				func(g internal.Getter) (
					uid uint,
					fieldValue *fvalue.Value,
					err error,
				) {
					return f2.Execute(g)
				},
				func(g internal.Getter) (
					uid uint,
					fieldValue *fvalue.Value,
					err error,
				) {
					return f3.Execute(g)
				},
			}

			res, err := l.ApplyOn(getter)

			require.NoError(t, err)
			require.Equal(
				t, res, fvalue.Values{
					0: fv0,
					1: fv1,
					2: fv2,
					3: fv3,
				},
			)

		},
	)
}

func TestApplyError_Is(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, errMocked1, ErrApply)
	require.ErrorIs(t, ApplyError{}, ErrApply)
}
