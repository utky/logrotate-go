package upload

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/utky/logproc-go/pkg/core"
	"github.com/utky/logproc-go/pkg/log"
)

// Archive is compressed file
type Archive struct {
	*core.Entry
	suffix string
}

func Run(archive *Archive) error {
	archive.Logger.Info("Start finalize",
		log.Fields{
			"file": archive.File.Base(),
		})
	timeStartFinalize := time.Now()
	fnErr := archive.Finalize()
	timeEndFinalize := time.Now()
	archive.Logger.Info("End finalize",
		log.Fields{
			"file":    archive.File.Base(),
			"elapsed": timeEndFinalize.Sub(timeStartFinalize),
		})
	return fnErr
}

// Finalize runs post-process action.
func (archive *Archive) Finalize() error {
	for _, cmdStr := range archive.Config.FinalizeCommands {
		timeStartFinalize := time.Now()
		cmd := exec.Command("sh", "-c", cmdStr)
		outputBytes, cmdErr := cmd.CombinedOutput()
		timeEndFinalize := time.Now()
		elapsedFinalize := timeEndFinalize.Sub(timeStartFinalize)
		output := string(outputBytes)
		if cmdErr != nil {
			archive.Logger.Errorf(
				"Finalize failed",
				log.Fields{
					"output":  output,
					"cmd":     cmdStr,
					"elapsed": elapsedFinalize,
				})
			return errors.Wrap(
				cmdErr,
				fmt.Sprintf("Finalize failed with stdout and stderr: %s", output))
		}
		archive.Logger.Infof(
			"Finalize succeeded",
			log.Fields{
				"output":  output,
				"cmd":     cmdStr,
				"elapsed": elapsedFinalize,
			})
	}
	return nil
}
