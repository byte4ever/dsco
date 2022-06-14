package dsco

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errMocked = errors.New("mocked error")

func Test_layers_getPostProcessErrors(t *testing.T) {
	t.Parallel()
	b1 := NewMockBinder2(t)
	b2 := NewMockBinder2(t)
	b3 := NewMockBinder2(t)

	b1.On("Errors").Return([]error{err1, err2}).Once()
	b2.On("Errors").Return(nil).Once()
	b3.On("Errors").Return([]error{err3, err4, err5}).Once()

	l := layers{b1, b2, b3}

	require.Equal(
		t,
		[]error{err1, err2, err3, err4, err5},
		l.getPostProcessErrors(),
	)
}

func Test_layers_bind(t *testing.T) {
	t.Parallel()

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			val1 := 1
			vValue1 := reflect.ValueOf(val1)

			val2 := 2
			vValue2 := reflect.ValueOf(val2)

			targetVal := -1
			targetValType := reflect.TypeOf(targetVal)
			targetValValue := reflect.ValueOf(targetVal)

			b1 := NewMockBinder2(t)
			b1.
				On(
					"Bind",
					"key_name",
					mock.MatchedBy(
						func(vt reflect.Type) bool {
							return targetValType.String() == vt.String()
						},
					),
				).
				Return(
					BindingAttempt{
						Value:    vValue1,
						Location: "loc-b1",
					},
				).
				Once()
			b1.
				On(
					"Use",
					"key_name",
				).
				Return(nil).
				Once()

			b2 := NewMockBinder2(t)
			b2.
				On(
					"Bind",
					"key_name",
					mock.MatchedBy(
						func(vt reflect.Type) bool {
							return targetValType.String() == vt.String()
						},
					),
				).
				Return(
					BindingAttempt{
						Value:    vValue2,
						Location: "loc-b2",
					},
				).
				Once()

			b3 := NewMockBinder2(t)
			b3.On(
				"Bind",
				"key_name",
				mock.MatchedBy(
					func(vt reflect.Type) bool {
						return targetValType.String() == vt.String()
					},
				),
			).Return(
				BindingAttempt{
					Error: err1,
				},
			).Once()

			l := layers{b1, b2, b3}
			re := l.bind("key_name", targetValValue)
			require.True(t, re.isFound())
			require.Equal(t, vValue1, re.Value)
			require.Equal(t, "loc-b1", re.ExternalKey)
			require.Equal(t, "key_name", re.Key)
			require.Equal(t, []error{nil, nil, err1}, re.Errors)
		},
	)

	t.Run(
		"success with no key found first",
		func(t *testing.T) {
			t.Parallel()

			val2 := 2
			vValue2 := reflect.ValueOf(val2)

			targetVal := -1
			targetValType := reflect.TypeOf(targetVal)
			targetValValue := reflect.ValueOf(targetVal)

			b1 := NewMockBinder2(t)
			b1.
				On(
					"Bind",
					"key_name",
					mock.MatchedBy(
						func(vt reflect.Type) bool {
							return targetValType.String() == vt.String()
						},
					),
				).
				Return(
					BindingAttempt{},
				).
				Once()

			b2 := NewMockBinder2(t)
			b2.
				On(
					"Bind",
					"key_name",
					mock.MatchedBy(
						func(vt reflect.Type) bool {
							return targetValType.String() == vt.String()
						},
					),
				).
				Return(
					BindingAttempt{
						Value:    vValue2,
						Location: "loc-b2",
					},
				).
				Once()

			b2.
				On(
					"Use",
					"key_name",
				).
				Return(nil).
				Once()

			b3 := NewMockBinder2(t)
			b3.On(
				"Bind",
				"key_name",
				mock.MatchedBy(
					func(vt reflect.Type) bool {
						return targetValType.String() == vt.String()
					},
				),
			).Return(
				BindingAttempt{
					Error: err1,
				},
			).Once()

			l := layers{b1, b2, b3}
			re := l.bind("key_name", targetValValue)
			require.True(t, re.isFound())
			require.Equal(t, vValue2, re.Value)
			require.Equal(t, "loc-b2", re.ExternalKey)
			require.Equal(t, "key_name", re.Key)
			require.Equal(t, []error{nil, nil, err1}, re.Errors)
		},
	)

	t.Run(
		"use panic",
		func(t *testing.T) {
			t.Parallel()

			val2 := 2
			vValue2 := reflect.ValueOf(val2)

			targetVal := -1
			targetValType := reflect.TypeOf(targetVal)
			targetValValue := reflect.ValueOf(targetVal)

			b1 := NewMockBinder2(t)
			b1.
				On(
					"Bind",
					"key_name",
					mock.MatchedBy(
						func(vt reflect.Type) bool {
							return targetValType.String() == vt.String()
						},
					),
				).
				Return(
					BindingAttempt{},
				).
				Once()

			b2 := NewMockBinder2(t)
			b2.
				On(
					"Bind",
					"key_name",
					mock.MatchedBy(
						func(vt reflect.Type) bool {
							return targetValType.String() == vt.String()
						},
					),
				).
				Return(
					BindingAttempt{
						Value:    vValue2,
						Location: "loc-b2",
					},
				).
				Once()

			b2.
				On(
					"Use",
					"key_name",
				).
				Return(errMocked).
				Once()

			b3 := NewMockBinder2(t)

			l := layers{b1, b2, b3}

			require.Panics(t, func() {
				_ = l.bind("key_name", targetValValue)
			})
		},
	)
}
