package kfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"

	"github.com/byte4ever/dsco/svalue"
)

const (
	fileNameExp = `^[A-Z][A-Z\d]*([-_][A-Z][A-Z\d]*)*$`
)

var (
	reFileName = regexp.MustCompile(fileNameExp)

	// ErrInvalidFileName represents an error indicating that the file name is
	// invalid.
	ErrInvalidFileName = errors.New("invalid kfile name")
)

type options struct {
	// silentDirErrors  bool
	silentFileErrors bool
}

func newProvider(
	fs afero.Fs,
	dirName string,
	opt *options,
) (
	*EntriesProvider,
	error,
) {
	cleanDirName := filepath.Clean(dirName)

	dirToSkip := len(
		strings.Split(
			cleanDirName,
			string(filepath.Separator),
		),
	)

	result := make(svalue.Values)

	var errs PathErrors

	_ = afero.Walk(
		fs, cleanDirName,
		walkFunc(fs, opt, &errs, result, cleanDirName, dirToSkip),
	)

	if len(errs) > 0 {
		return nil, errs
	}

	if len(result) == 0 {
		return &EntriesProvider{}, nil
	}

	return &EntriesProvider{
		values: result,
		name: fmt.Sprintf(
			"kfile(%s)",
			cleanDirName,
		),
	}, nil
}

func walkFunc(
	fs afero.Fs,
	opt *options,
	errs *PathErrors,
	result svalue.Values,
	cleanDirName string,
	dirToSkip int,
) func(
	path string,
	info os.FileInfo,
	err error,
) error {
	appendError := func(path string, err error) {
		*errs = append(
			*errs,
			&pathError{
				path: path,
				err:  err,
			},
		)
	}

	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if !opt.silentFileErrors {
				appendError(path, err)
			}

			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !reFileName.MatchString(info.Name()) {
			appendError(path, ErrInvalidFileName)

			return nil
		}

		fileContent, err := afero.ReadFile(fs, path)
		if err != nil {
			appendError(path, err)

			return nil
		}

		result[strings.ToLower(info.Name())] = &svalue.Value{
			Location: fmt.Sprintf(
				"kfile[%s]:%s",
				cleanDirName,
				filepath.Join(
					strings.Split(path, string(filepath.Separator))[dirToSkip:]...,
				),
			),
			Value: string(fileContent),
		}

		return nil
	}
}

type EntriesProvider struct {
	values svalue.Values
	name   string
}

func (e *EntriesProvider) GetName() string {
	return e.name
}

func (e *EntriesProvider) GetStringValues() svalue.Values {
	return e.values
}

func NewEntriesProvider(
	path string,
) (
	*EntriesProvider,
	error,
) {
	return newProvider(
		afero.NewReadOnlyFs(afero.NewOsFs()),
		path,
		&options{},
	)
}
