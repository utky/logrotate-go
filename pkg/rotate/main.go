package rotate

import (
	"fmt"
	"time"

	"github.com/utky/logproc-go/pkg/log"
)

var logger = log.New()

// Base is common structure which all stages of log should have.
type Base struct {
	config *Config
	file   File
}

// Source is a original file before rotated.
type Source struct {
	*Base
	owners []Owner
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
				logger.Warnf("Failed to query owner of file", log.Fields{"error": rlsErr})
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
	logger.Info("Notified signal to release handle", log.Fields{"file": source.file.AbsolutePath()})

	rlsErr := WaitOwnerRelease(source)
	if rlsErr != nil {
		return temp, rlsErr
	}
	logger.Info("Completed to wait release handle", log.Fields{"file": source.file.AbsolutePath()})
	temp = &Temp{
		Base: source.Base,
	}
	return temp, nil
}

// Temp is a temporary state of rotation.
type Temp struct {
	*Base
}

// Compress archives current temp file.
func (temp *Temp) Compress() (*Archive, error) {
	archive := &Archive{
		Base: temp.Base,
	}
	return archive, nil
}

// Archive is compressed file
type Archive struct {
	*Base
}

// Finalize runs post-process action.
func (archive *Archive) Finalize() error {
	return nil
}

// RunRotate runs log processing pipeline
func RunRotate(source *Source) error {

	logger.Info("Start evacuate",
		log.Fields{
			"source": source.file.Basename(),
		})
	timeStartEvacuate := time.Now()
	temp, evErr := source.Evacuate()
	timeEndEvacuate := time.Now()
	logger.Info("End evacuate",
		log.Fields{
			"file":    source.file.Basename(),
			"elapsed": timeEndEvacuate.Sub(timeStartEvacuate),
		})
	if evErr != nil {
		return evErr
	}

	logger.Info("Start compress",
		log.Fields{
			"file": temp.file.Basename(),
		})
	timeStartArchive := time.Now()
	archive, cmErr := temp.Compress()
	timeEndArchive := time.Now()
	logger.Info("End evacuate",
		log.Fields{
			"file":    temp.file.Basename(),
			"elapsed": timeEndArchive.Sub(timeStartArchive),
		})
	if cmErr != nil {
		return cmErr
	}

	logger.Info("Start finalize",
		log.Fields{
			"file": archive.file.Basename(),
		})
	timeStartFinalize := time.Now()
	fnErr := archive.Finalize()
	timeEndFinalize := time.Now()
	logger.Info("End finalize",
		log.Fields{
			"file":    archive.file.Basename(),
			"elapsed": timeEndFinalize.Sub(timeStartFinalize),
		})
	return fnErr
}

// NewSource creates source
func NewSource(config *Config, file File, owners []Owner) *Source {
	base := &Base{
		config: config,
		file:   file,
	}
	source := &Source{
		Base:   base,
		owners: owners,
	}
	return source
}
