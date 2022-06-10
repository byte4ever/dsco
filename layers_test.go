package dsco

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_layers_getPostProcessErrors(t *testing.T) {
	t.Parallel()
	b1 := NewMockBinder(t)
	b2 := NewMockBinder(t)
	b3 := NewMockBinder(t)

	b1.On("GetPostProcessErrors").Return([]error{err1, err2}).Once()
	b2.On("GetPostProcessErrors").Return(nil).Once()
	b3.On("GetPostProcessErrors").Return([]error{err3, err4, err5}).Once()

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
			fakeOrigin1 := Origin("mocked1")
			fakeOrigin2 := Origin("mocked2")
			fakeOrigin3 := Origin("mocked3")

			val1 := 1
			vValue1 := reflect.ValueOf(val1)

			val2 := 2
			vValue2 := reflect.ValueOf(val2)

			val3 := 3
			vValue3 := reflect.ValueOf(val3)

			targetVal := -1
			targetValValue := reflect.ValueOf(targetVal)

			b1 := NewMockBinder(t)
			b1.On(
				"Bind",
				"key_name",
				true,
				targetValValue,
			).Return(
				fakeOrigin1,
				"keyOut1",
				true,
				vValue1,
				nil,
			).Once()

			b2 := NewMockBinder(t)
			b2.On(
				"Bind",
				"key_name",
				false,
				targetValValue,
			).Return(
				fakeOrigin2,
				"keyOut2",
				true,
				vValue2,
				nil,
			).Once()

			b3 := NewMockBinder(t)
			b3.On(
				"Bind",
				"key_name",
				false,
				targetValValue,
			).Return(
				fakeOrigin3,
				"keyOut3",
				false,
				vValue3,
				err1,
			).Once()

			l := layers{b1, b2, b3}
			re := l.bind("key_name", targetValValue)
			require.True(t, re.isFound())
			require.Equal(t, vValue1, re.Value)
			require.Equal(t, "keyOut1", re.ExternalKey)
			require.Equal(t, "key_name", re.Key)
			require.Equal(t, []error{nil, nil, err1}, re.Errors)
		},
	)
}
