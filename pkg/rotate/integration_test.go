package rotate

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func createTemp(t *testing.T) string {
	d, tempErr := ioutil.TempDir("test/", "integration-test")
	if tempErr != nil {
		t.Error("Failed to create temp dir")
	}
	return d
}

func prepareSource(t *testing.T, dir string, filename string, content string) {
	p := path.Join(dir, filename)
	err := ioutil.WriteFile(p, []byte(content), 0644)
	if err != nil {
		t.Error("Failed to write content")
	}
}

func Test01(t *testing.T) {
	d := createTemp(t)
	defer os.RemoveAll(d)
	prepareSource(t, d, "access.log", "hello")

}
