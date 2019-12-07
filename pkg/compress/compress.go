package compress

import (
	"fmt"
	"os/exec"
	"path"
	"time"

	"github.com/utky/logproc-go/pkg/core"
	"github.com/utky/logproc-go/pkg/log"
)

// Temp is a temporary state of rotation.
type Temp struct {
	*core.Entry
}

// NewTemp creates source
func NewTemp(config *core.Config, file core.File) *Temp {
	logger := log.NewWithFields(log.Fields{"file": file.Base()})
	base := &core.Entry{
		Config: config,
		File:   file,
		Logger: logger,
	}
	temp := &Temp{
		Entry: base,
	}
	return temp
}

// Compress archives current temp file.
func (temp *Temp) Compress() error {
	dest := path.Join(temp.Config.TempStorage, temp.File.Base()+".gz")
	gzip := &Gzip{}
	gzipErr := gzip.Run(temp.File, dest)
	return gzipErr
}

// Run archives temp file
func Run(temp *Temp) error {
	temp.Logger.Info("Start compress",
		log.Fields{
			"file": temp.File.Base(),
		})
	timeStartArchive := time.Now()
	cmErr := temp.Compress()
	timeEndArchive := time.Now()
	temp.Logger.Info("End evacuate",
		log.Fields{
			"file":    temp.File.Base(),
			"elapsed": timeEndArchive.Sub(timeStartArchive),
		})
	return cmErr
}

// Compressor provides a way to compress file.
type Compressor interface {
	Run(file core.File, dest string) error
}

// Gzip uses gzip as file compressor
type Gzip struct {
}

// Run uses gzip via sh
func (gzip *Gzip) Run(file core.File, dest string) error {
	logger := log.NewWithFields(log.Fields{"file": file.Abs()})

	src := file.Abs()
	args := []string{"-c", fmt.Sprintf("gzip -q -c '%s' > '%s'", src, dest)}
	cmd := exec.Command("sh", args...)
	out, err := cmd.CombinedOutput()
	logger.Debugf("output of gzip", log.Fields{"out": out})
	return err
}
