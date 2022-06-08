package searcher

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

var errMocked = errors.New("mocked error")

func buildTestFile(t *testing.T, fileName string, fileContent []byte) (string, string, string, string) {
	t.Helper()

	rootName := t.TempDir()
	p1 := path.Join(rootName, "d1")
	p2 := path.Join(rootName, "d2")
	p3 := path.Join(rootName, "d3")
	pf := path.Join(p1, fileName)

	require.NoError(t, os.Mkdir(p1, 0777))
	require.NoError(t, os.Mkdir(p2, 0777))
	require.NoError(t, os.Mkdir(p3, 0777))

	f, err := os.Create(pf)
	require.NoError(t, err)
	_, err = f.Write(fileContent)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	return p1, p2, p3, pf
}

func Test_NewSearchPath(t *testing.T) {
	t.Run(
		"success", func(t *testing.T) {
			p1, p2, p3, pf := buildTestFile(t, "f", []byte("content"))

			paths := []string{p1, p2, p3}

			for i := 0; i < 3; i++ {
				of, err := searchfile(paths, "f")
				require.NoError(t, err)

				require.Equal(t, pf, of.Name())
				require.NoError(t, of.Close())
				paths = append(paths[2:], paths[:2]...)

			}
		},
	)

	t.Run(
		"err when no search path is provided (nil)", func(t *testing.T) {
			s, err := NewMultiPath(nil, "f")
			require.ErrorIs(t, err, ErrNoSearchPath)
			require.Nil(t, s)
		},
	)

	t.Run(
		"err when no search path is provided (empty)", func(t *testing.T) {
			s, err := NewMultiPath(make([]string, 0), "f")
			require.ErrorIs(t, err, ErrNoSearchPath)
			require.Nil(t, s)
		},
	)

	t.Run(
		"file not found", func(t *testing.T) {
			p1, p2, p3, _ := buildTestFile(t, "f", []byte("content"))

			paths := []string{p1, p2, p3}

			s, err := NewMultiPath(paths, "not-present-file")
			require.ErrorIs(t, err, ErrConfNotFound)
			require.Nil(t, s)
		},
	)
}

func TestMultiPathFileSearcher_ReadClose(t *testing.T) {
	t.Run(
		"success", func(t *testing.T) {
			expectedContent := []byte("content")
			p1, p2, p3, _ := buildTestFile(t, "f", expectedContent)

			paths := []string{p1, p2, p3}
			s, err := NewMultiPath(paths, "f")
			require.NoError(t, err)
			require.NoError(
				t, s.ReadClose(
					func(reader io.Reader) error {
						content, err := ioutil.ReadAll(reader)
						require.NoError(t, err)
						require.Equal(t, expectedContent, content)
						return nil
					},
				),
			)
		},
	)

	t.Run(
		"perform error", func(t *testing.T) {
			expectedContent := []byte("content")

			p1, p2, p3, _ := buildTestFile(t, "f", expectedContent)

			paths := []string{p1, p2, p3}
			s, err := NewMultiPath(paths, "f")
			require.NoError(t, err)
			require.ErrorIs(
				t,
				s.ReadClose(
					func(reader io.Reader) error {
						return errMocked
					},
				),
				errMocked,
			)
		},
	)

	t.Run(
		"handle close error", func(t *testing.T) {
			expectedContent := []byte("content")
			p1, p2, p3, _ := buildTestFile(t, "f", expectedContent)

			paths := []string{p1, p2, p3}
			s, err := NewMultiPath(paths, "f")
			require.NoError(t, err)

			require.NoError(t, s.file.Close())

			require.ErrorIs(
				t,
				s.ReadClose(
					func(reader io.Reader) error {
						return nil
					},
				),
				os.ErrClosed,
			)
		},
	)
}
