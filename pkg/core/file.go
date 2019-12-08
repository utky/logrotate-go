package core

import (
	"os"
	"path/filepath"
)

// File is a abstracted interface to manipulate file.
// This file does not mean opened.
type File interface {
	Base() string
	Abs() string
	Move(dest string) error
}

// OsFile has backend of OS file system.
type OsFile struct {
	path string
}

// NewOsFile creates file backed by OS file system.
func NewOsFile(path string) (*OsFile, error) {
	var osFile *OsFile
	apath, err := filepath.Abs(path)
	if err != nil {
		return osFile, err
	}
	osFile = &OsFile{
		path: apath,
	}
	return osFile, nil
}

// Base implements File interface
func (f *OsFile) Base() string {
	return filepath.Base(f.path)
}

// Abs implements File interface
func (f *OsFile) Abs() string {
	return f.path
}

// Move implements File interface
func (f *OsFile) Move(dest string) error {
	return os.Rename(f.Abs(), dest)
}

// NewFileFunc is callback to build implementation of File interface
type NewFileFunc = func(path string) (File, NewFileError)

// NewFileError should be raise if NewFileFunc fails or tells skipping
type NewFileError struct {
	err  error
	skip bool
}

// FailNewFile notifies failure to create file.
func FailNewFile(err error) NewFileError {
	return NewFileError{
		err:  err,
		skip: false,
	}
}

// SkipFile notifies caller should skip this result.
func SkipFile() NewFileError {
	return NewFileError{
		err:  nil,
		skip: true,
	}
}

// Empty creates nil values of NewFileError
func Empty() NewFileError {
	return NewFileError{
		err:  nil,
		skip: false,
	}
}

// CollectFiles reads file system and build matched files.
func CollectFiles(callback NewFileFunc, pattern string) ([]File, error) {
	result := make([]File, 0)
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return result, err
	}
	for _, p := range paths {
		f, nfe := callback(p)
		if nfe.err != nil {
			return result, nfe.err
		}
		if nfe.skip {
			continue
		}
		result = append(result, f)
	}
	return result, nil
}
