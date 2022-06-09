package searcher

import (
	"io"
	"os"
	"path/filepath"
)

// MultiPathFileSearcher represents a configuration file searcher that supports
// multiple paths.
type MultiPathFileSearcher struct {
	file  *os.File
	paths []string
}

func searchfile(paths []string, fileName string) (*os.File, error) {
	if len(paths) == 0 {
		return nil, ErrNoSearchPath
	}

	for _, path := range paths {
		fp := filepath.Join(path, fileName)

		if f, err := os.Open(fp); err == nil {
			return f, nil
		}
	}

	return nil, ErrConfNotFound
}

// NewMultiPath creates a configuration searcher.
func NewMultiPath(paths []string, fileName string) (*MultiPathFileSearcher, error) {
	f, err := searchfile(paths, fileName)
	if err != nil {
		return nil, err
	}

	return &MultiPathFileSearcher{
		file:  f,
		paths: paths,
	}, nil
}

// Apply applies action using the reader.
func (m *MultiPathFileSearcher) Apply(action func(r io.Reader) error) (err error) {
	defer func() {
		errClose := m.file.Close()
		if errClose != nil && err == nil {
			err = errClose
		}
	}()

	return action(m.file)
}
