package core

import (
	"io/ioutil"
	"path"
	"testing"
)

// CreateTemp creates temp directory.
func CreateTemp(t *testing.T) string {
	d, tempErr := ioutil.TempDir("", "rotate-test")
	if tempErr != nil {
		t.Error("Failed to create temp dir")
	}
	return d
}

// PrepareFile creates OS backed file.
func PrepareFile(t *testing.T, dir string, filename string, content string) *OsFile {
	p := path.Join(dir, filename)
	err := ioutil.WriteFile(p, []byte(content), 0644)
	if err != nil {
		t.Errorf("Failed to write content: %s", err)
	}
	f, newErr := NewOsFile(p)
	if newErr != nil {
		t.Errorf("Failed to instantiate OsFile: %s", newErr)
	}
	return f
}
