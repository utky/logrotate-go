package compress

import (
	"os"
	"path"
	"testing"

	"github.com/utky/logproc-go/pkg/core"
)

func TestCompress(t *testing.T) {
	d := core.CreateTemp(t)
	defer os.RemoveAll(d)
	f := core.PrepareFile(t, d, "TestCompress", "TestCompress")
	dest := path.Join(d, "TestCompress.gzip")
	gzip := &Gzip{}
	gzipErr := gzip.Run(f, dest)
	if gzipErr != nil {
		t.Errorf("Failed to compress gzip: %s", gzipErr)
	}
	if _, absentErr := os.Stat(dest); os.IsNotExist(absentErr) {
		t.Errorf("Cannot find compressed file.: %s", absentErr)
	}
}
