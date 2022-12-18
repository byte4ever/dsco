package kfile

import (
	"errors"
	"fmt"
	"os"
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
	result := make(svalue.Values)

	var errs PathErrors

	_ = afero.Walk(
		fs, dirName,
		func(path string, info os.FileInfo, err error) error {
			// fmt.Println(path, info.Mode(), info.Sys(), err)

			if err != nil {
				if !opt.silentFileErrors {
					errs = append(
						errs, &pathError{
							path: path,
							err:  err,
						},
					)
				}
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if !reFileName.MatchString(info.Name()) {
				errs = append(
					errs,
					&pathError{
						path: path,
						err:  ErrInvalidFileName,
					},
				)
				return nil
			}

			fileContent, err := afero.ReadFile(fs, path)
			if err != nil {
				errs = append(
					errs,
					&pathError{
						path: path,
						err:  err,
					},
				)
				return nil
			}

			result[strings.ToLower(info.Name())] = &svalue.Value{
				Location: fmt.Sprintf("kfile[%s]", path),
				Value:    string(fileContent),
			}

			return nil
		},
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
			dirName,
		),
	}, nil
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
