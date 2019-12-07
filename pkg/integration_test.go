// test
package rotate

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/utky/logproc-go/pkg/core"
)

func prepareSource(t *testing.T, dir string, filename string, content string) {
	p := path.Join(dir, filename)
	err := ioutil.WriteFile(p, []byte(content), 0644)
	if err != nil {
		t.Error("Failed to write content")
	}
}

func Test01Normal(t *testing.T) {
	d := core.CreateTemp(t)
	defer os.RemoveAll(d)
	prepareSource(t, d, "access.log", "hello")

}

func Test02HasTemp(t *testing.T) {
}

func Test03HasArchive(t *testing.T) {
}

func Test04HasRemove(t *testing.T) {

}
