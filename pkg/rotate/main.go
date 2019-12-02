package rotate

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/utky/logproc-go/pkg/log"
)

// Base is common structure which all stages of log should have.
type Base struct {
	config *Config
	file   File
	logger *log.Logger
}

// Source is a original file before rotated.
type Source struct {
	*Base
	owners []Owner
}

// Temp is a temporary state of rotation.
type Temp struct {
	*Base
}

// Archive is compressed file
type Archive struct {
	*Base
}

// WaitOwnerRelease wait with check interval and timeout that all owner release handle of the file.
func WaitOwnerRelease(source *Source) error {
	var err error
	err = nil
	timeoutCh := time.After(source.config.OwnerReleaseTimeout)
	intervalCh := time.Tick(source.config.OwnerReleaseInterval)
	for wait := true; wait; {
		select {
		case <-intervalCh:
			released, rlsErr := AllOwnersReleased(source.owners, source.file)
			if released && rlsErr == nil {
				wait = false
				err = rlsErr
			}
			if rlsErr != nil {
				source.logger.Warnf("Failed to query owner of file", log.Fields{"error": rlsErr})
			}
		case <-timeoutCh:
			wait = false
			err = fmt.Errorf("Timedout to wait file handle released: %s", source.file.AbsolutePath())
		default:
		}
	}
	return err
}

// Evacuate moves original log to temporary storage and wait owner to release handle to the file.
func (source *Source) Evacuate() (*Temp, error) {
	var temp *Temp
	notifyErr := NotifyAll(source.owners, source.file)
	if notifyErr != nil {
		return temp, notifyErr
	}
	source.logger.Info("Notified signal to release handle")

	rlsErr := WaitOwnerRelease(source)
	if rlsErr != nil {
		return temp, rlsErr
	}
	source.logger.Info("Completed to wait release handle")
	temp = &Temp{
		Base: source.Base,
	}
	return temp, nil
}

// Compress archives current temp file.
func (temp *Temp) Compress() (*Archive, error) {
	archive := &Archive{
		Base: temp.Base,
	}
	return archive, nil
}

// Finalize runs post-process action.
func (archive *Archive) Finalize() error {
	for _, cmdStr := range archive.config.FinalizeCommands {
		timeStartFinalize := time.Now()
		cmd := exec.Command("sh", "-c", cmdStr)
		outputBytes, cmdErr := cmd.CombinedOutput()
		timeEndFinalize := time.Now()
		elapsedFinalize := timeEndFinalize.Sub(timeStartFinalize)
		output := string(outputBytes)
		if cmdErr != nil {
			archive.logger.Errorf(
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
		archive.logger.Infof(
			"Finalize succeeded",
			log.Fields{
				"output":  output,
				"cmd":     cmdStr,
				"elapsed": elapsedFinalize,
			})
	}
	return nil
}

// RunRotate runs log processing pipeline
func RunRotate(source *Source) error {

	source.logger.Info("Start evacuate",
		log.Fields{
			"source": source.file.Basename(),
		})
	timeStartEvacuate := time.Now()
	temp, evErr := source.Evacuate()
	timeEndEvacuate := time.Now()
	source.logger.Info("End evacuate",
		log.Fields{
			"file":    source.file.Basename(),
			"elapsed": timeEndEvacuate.Sub(timeStartEvacuate),
		})
	if evErr != nil {
		return evErr
	}

	temp.logger.Info("Start compress",
		log.Fields{
			"file": temp.file.Basename(),
		})
	timeStartArchive := time.Now()
	archive, cmErr := temp.Compress()
	timeEndArchive := time.Now()
	temp.logger.Info("End evacuate",
		log.Fields{
			"file":    temp.file.Basename(),
			"elapsed": timeEndArchive.Sub(timeStartArchive),
		})
	if cmErr != nil {
		return cmErr
	}

	archive.logger.Info("Start finalize",
		log.Fields{
			"file": archive.file.Basename(),
		})
	timeStartFinalize := time.Now()
	fnErr := archive.Finalize()
	timeEndFinalize := time.Now()
	archive.logger.Info("End finalize",
		log.Fields{
			"file":    archive.file.Basename(),
			"elapsed": timeEndFinalize.Sub(timeStartFinalize),
		})
	return fnErr
}

// NewSource creates source
func NewSource(config *Config, file File, owners []Owner) *Source {
	logger := log.NewWithFields(log.Fields{"file": file.Basename})
	base := &Base{
		config: config,
		file:   file,
		logger: logger,
	}
	source := &Source{
		Base:   base,
		owners: owners,
	}
	return source
}
