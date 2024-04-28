package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/internal"
	"github.com/byte4ever/dsco/internal/merror"
)

func TestExpandList_Count(t *testing.T) {
	tests := []struct {
		name string
		s    ExpandList
		want int
	}{
		{
			name: "nil",
			s:    nil,
		},
		{
			name: "empty",
			s:    ExpandList{},
		},
		{
			name: "2 items",
			s: ExpandList{
				func(g internal.StructExpander) (err error) {
					return nil
				},
				func(g internal.StructExpander) (err error) {
					return nil
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.s.Count(), "Count()")
		})
	}
}

func TestExpandList_Push(t *testing.T) {
	t.Parallel()

	type args struct {
		o ExpandOp
	}

	tests := []struct {
		name string
		s    ExpandList
		args args
	}{
		{
			name: "",
			s:    nil,
			args: args{
				o: func(g internal.StructExpander) (err error) {
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sizeBefore := len(tt.s)
			tt.s.Push(tt.args.o)
			require.Equal(t, sizeBefore+1, len(tt.s))
		})
	}
}

func TestExpandList_ApplyOn(t *testing.T) {
	t.Run(
		"catch error", func(t *testing.T) {
			expand := newMockStructExpander(t)

			op1 := NewMockExpandOp(t)
			op1.
				EXPECT().
				Execute(expand).
				Return(errMocked1)
			op2 := NewMockExpandOp(t)
			op2.
				EXPECT().
				Execute(expand).
				Return(nil)
			op3 := NewMockExpandOp(t)
			op3.
				EXPECT().
				Execute(expand).
				Return(errMocked2)

			expandList := ExpandList{
				op1.Execute,
				op2.Execute,
				op3.Execute,
			}

			err := expandList.ApplyOn(expand)

			var asErr ApplyError
			require.ErrorAs(t, err, &asErr)
			require.Equal(t, asErr.MError, merror.MError([]error{
				errMocked1,
				errMocked2,
			}))
		},
	)

	t.Run(
		"success", func(t *testing.T) {
			expand := newMockStructExpander(t)

			op1 := NewMockExpandOp(t)
			op1.
				EXPECT().
				Execute(expand).
				Return(nil)
			op2 := NewMockExpandOp(t)
			op2.
				EXPECT().
				Execute(expand).
				Return(nil)
			op3 := NewMockExpandOp(t)
			op3.
				EXPECT().
				Execute(expand).
				Return(nil)

			expandList := ExpandList{
				op1.Execute,
				op2.Execute,
				op3.Execute,
			}

			err := expandList.ApplyOn(expand)
			require.NoError(t, err)
		},
	)
}
