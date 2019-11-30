package rotate

import "os"

// Storage is a collection of files.
type Storage interface {
	Id() string
	Has(entry File) bool
	Add(entry File) error
	Delete(entry File) error
}

// Directory as storage
type Directory struct {
	directory os.File
}
