package evacuate

import (
	"fmt"
	"path"
	"time"

	"github.com/utky/logproc-go/pkg/core"
	"github.com/utky/logproc-go/pkg/log"
)

// Source is a original file before rotated.
type Source struct {
	*core.Entry
	owners []Owner
}

// WaitOwnerRelease wait with check interval and timeout that all owner release handle of the file.
func WaitOwnerRelease(source *Source) error {
	var err error
	err = nil
	timeoutCh := time.After(source.Config.OwnerReleaseTimeout)
	intervalCh := time.Tick(source.Config.OwnerReleaseInterval)
	for wait := true; wait; {
		select {
		case <-intervalCh:
			released, rlsErr := AllOwnersReleased(source.owners, source.File)
			if released && rlsErr == nil {
				wait = false
				err = rlsErr
			}
			if rlsErr != nil {
				source.Logger.Warnf("Failed to query owner of file", log.Fields{"error": rlsErr})
			}
		case <-timeoutCh:
			wait = false
			err = fmt.Errorf("Timedout to wait file handle released: %s", source.File.Abs())
		default:
		}
	}
	return err
}

// Evacuate moves original log to temporary storage and wait owner to release handle to the file.
func (source *Source) Evacuate() error {
	notifyErr := NotifyAll(source.owners, source.File)
	if notifyErr != nil {
		return notifyErr
	}
	source.Logger.Info("Notified signal to release handle")

	rlsErr := WaitOwnerRelease(source)
	if rlsErr != nil {
		return rlsErr
	}
	source.Logger.Info("Completed to wait release handle")

	tempPath := path.Join(source.Config.TempStorage, source.File.Base())
	if mvErr := source.File.Move(tempPath); mvErr != nil {
		return mvErr
	}
	return nil
}

// Run runs log processing pipeline
func Run(source *Source) error {

	source.Logger.Info("Start evacuate",
		log.Fields{
			"source": source.File.Base(),
		})
	timeStartEvacuate := time.Now()
	evErr := source.Evacuate()
	timeEndEvacuate := time.Now()
	source.Logger.Info("End evacuate",
		log.Fields{
			"file":    source.File.Base(),
			"elapsed": timeEndEvacuate.Sub(timeStartEvacuate),
		})
	return evErr
}

// NewSource creates source
func NewSource(config *core.Config, file core.File, owners []Owner) *Source {
	logger := log.NewWithFields(log.Fields{"file": file.Base()})
	base := &core.Entry{
		Config: config,
		File:   file,
		Logger: logger,
	}
	source := &Source{
		Entry:  base,
		owners: owners,
	}
	return source
}
