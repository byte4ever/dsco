package yamlp

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

var errMocked = errors.New("mocked error")

type T1Root struct {
	A int
	B float64
}

type failReaderCloser struct{}

type bufferReaderCloser struct {
	buf []byte
}

func (*failReaderCloser) Apply(func(r io.Reader) error) error {
	return errMocked
}

func (b *bufferReaderCloser) Apply(f func(r io.Reader) error) error {
	return f(bytes.NewReader(b.buf))
}

func TestProvide(t *testing.T) {
	t.Parallel()

	t.Run(
		"error when model interface is nil",
		func(t *testing.T) {
			t.Parallel()

			provider, err := New(nil, nil)
			require.ErrorIs(t, err, ErrInvalidModel)
			require.ErrorContains(t, err, "nil")
			require.Nil(t, provider)
		},
	)

	t.Run(
		"error when model interface not a pointer",
		func(t *testing.T) {
			t.Parallel()

			provider, err := New(123, nil)
			require.ErrorIs(t, err, ErrInvalidModel)
			require.ErrorContains(t, err, "pointer")
			require.Nil(t, provider)
		},
	)

	t.Run(
		"error when model interface not a pointer on struct",

		func(t *testing.T) {
			t.Parallel()

			v := 5
			provider, err := New(&v, nil)
			require.ErrorIs(t, err, ErrInvalidModel)
			require.ErrorContains(t, err, "struct")
			require.Nil(t, provider)
		},
	)

	t.Run(
		"error when performer is nil",
		func(t *testing.T) {
			t.Parallel()

			k := &struct{}{}
			provider, err := New(k, nil)
			require.ErrorIs(t, err, ErrNilReaderFunctor)
			require.ErrorContains(t, err, "nil")
			require.Nil(t, provider)
		},
	)

	t.Run(
		"reader provider internal failure",
		func(t *testing.T) {
			t.Parallel()

			k := &struct{}{}
			mrc := &failReaderCloser{}
			provider, err := New(k, mrc)
			require.ErrorIs(t, err, errMocked)
			require.Nil(t, provider)
		},
	)

	t.Run(
		"invalid yaml content",
		func(t *testing.T) {
			t.Parallel()

			mrc := &bufferReaderCloser{
				buf: []byte("invalid yaml content"),
			}
			k := &T1Root{}
			provider, err := New(k, mrc)
			require.Nil(t, provider)
			var e *yaml.TypeError
			require.ErrorAs(t, err, &e)
			require.ErrorContains(t, e, "line 1")
			require.ErrorContains(t, e, "yaml")
		},
	)

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			mrc := &bufferReaderCloser{
				buf: []byte(`
a: 123
b: 999.99
`),
			}
			k := &T1Root{}
			provider, err := New(k, mrc)
			require.NoError(t, err)
			require.NotNil(t, provider)

			ri, err := provider.GetInterface()
			require.NoError(t, err)

			require.NotNil(t, ri)
			require.IsType(t, &T1Root{}, ri)
			fi, ok := ri.(*T1Root)
			require.True(t, ok)
			require.Equal(t, 123, fi.A)
			require.Equal(t, 999.99, fi.B)
		},
	)
}
