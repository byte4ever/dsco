package dsco

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_layers_getPostProcessErrors(t *testing.T) {
	b1 := NewMockBinder(t)
	b2 := NewMockBinder(t)
	b3 := NewMockBinder(t)

	b1.On("GetPostProcessErrors").Return([]error{err1, err2}).Once()
	b2.On("GetPostProcessErrors").Return(nil).Once()
	b3.On("GetPostProcessErrors").Return([]error{err3, err4, err5}).Once()

	l := layers{b1, b2, b3}

	require.Equal(t, []error{err1, err2, err3, err4, err5}, l.getPostProcessErrors())
}
