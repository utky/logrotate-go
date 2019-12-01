package rotate

import (
	"os"
	"path/filepath"
)

// File is a abstracted interface to manipulate file.
type File interface {
	Basename() string
	AbsolutePath() string
	MoveTo(dest string) error
}

// OsFile has backend of OS file system.
type OsFile struct {
	path string
}

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

func (f *OsFile) Basename() string {
	return filepath.Base(f.path)
}
func (f *OsFile) AbsolutePath() string {
	return f.path
}
func (f *OsFile) MoveTo(dest string) error {
	return os.Rename(f.AbsolutePath(), dest)
}
