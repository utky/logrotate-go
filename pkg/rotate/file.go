package rotate

import "os"

// File is a abstracted interface to manipulate file.
type File interface {
	ID() string
	AbsolutePath() string
	MoveTo(storage Storage) error
}

// OsFile has backend of OS file system.
type OsFile struct {
	file os.File
}
