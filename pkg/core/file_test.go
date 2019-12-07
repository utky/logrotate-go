package core

import (
	"os"
	"path"
	"testing"
)

func TestMove(t *testing.T) {
	srcDir := CreateTemp(t)
	defer os.RemoveAll(srcDir)
	destDir := CreateTemp(t)
	defer os.RemoveAll(destDir)
	f := PrepareFile(t, srcDir, "TestMove", "TestMove")
	if e := f.Move(path.Join(destDir, "TestMove")); e != nil {
		t.Errorf("Failed to move file: %s", e)
	}
}
