package yaml_provider

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

var errMocked = errors.New("mocked error")

type failReaderCloser struct{}

func (f *failReaderCloser) ReadClose(func(r io.Reader) error) error {
	return errMocked
}

type bufferReaderCloser struct {
	buf []byte
}

func (b *bufferReaderCloser) ReadClose(f func(r io.Reader) error) error {
	return f(bytes.NewReader(b.buf))
}

type T1Root struct {
	A int
	B float64
}

func TestProvide(t *testing.T) {
	t.Run(
		"return error if interface is nil", func(t *testing.T) {
			provider, err := Provide(nil, nil)
			require.ErrorIs(t, err, ErrNilInterfaces)
			require.Nil(t, provider)
		},
	)

	t.Run(
		"reader provider internal failure", func(t *testing.T) {
			k := &struct{}{}
			mrc := &failReaderCloser{}
			provider, err := Provide(k, mrc)
			require.ErrorIs(t, err, errMocked)
			require.Nil(t, provider)
		},
	)

	t.Run(
		"invalid yaml content", func(t *testing.T) {
			mrc := &bufferReaderCloser{
				buf: []byte("invalid yaml content"),
			}
			k := &T1Root{}
			provider, err := Provide(k, mrc)
			require.Nil(t, provider)
			var e *yaml.TypeError
			require.ErrorAs(t, err, &e)
			require.ErrorContains(t, e, "line 1")
			require.ErrorContains(t, e, "yaml")
		},
	)

	t.Run(
		"success", func(t *testing.T) {
			mrc := &bufferReaderCloser{
				buf: []byte(`
a: 123
b: 999.99
`),
			}
			k := &T1Root{}
			provider, err := Provide(k, mrc)
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
