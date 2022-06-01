package searcher

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// ErrConfNotFound is returned when no file is found.
var ErrConfNotFound = errors.New("no configuration file found")
var ErrNoSearchPath = errors.New("no search path")

type MultiPathFileSearcher struct {
	file  *os.File
	paths []string
}

func NewMultiPath(paths []string, fileName string) (*MultiPathFileSearcher, error) {
	if len(paths) == 0 {
		return nil, ErrNoSearchPath
	}

	for _, path := range paths {
		fp := filepath.Join(path, fileName)
		input, err := os.Open(fp)

		if err == nil {
			return &MultiPathFileSearcher{
				file:  input,
				paths: paths,
			}, nil
		}
	}

	return nil, ErrConfNotFound
}

func (m *MultiPathFileSearcher) ReadClose(perform func(r io.Reader) error) (err error) {
	defer func() {
		errClose := m.file.Close()
		if errClose != nil && err == nil {
			err = errClose
		}
	}()

	return perform(m.file)
}

func (m *MultiPathFileSearcher) ProvideFile() (*os.File, error) {
	return m.file, nil
}
