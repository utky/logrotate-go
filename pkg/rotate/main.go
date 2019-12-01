package rotate

import (
	"fmt"
	"time"
)

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
	timeoutCh := time.After(source.config.ownerReleaseTimeout)
	intervalCh := time.Tick(source.config.ownerReleaseInterval)
	for wait := true; wait; {
		select {
		case <-intervalCh:
			released, rlsErr := AllOwnersReleased(source.owners, source.file)
			if released && rlsErr == nil {
				wait = false
				err = rlsErr
			}
		case <-timeoutCh:
			wait = false
			err = fmt.Errorf("Timedout to wait file handle released")
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
	rlsErr := WaitOwnerRelease(source)
	if rlsErr != nil {
		return temp, rlsErr
	}
	temp = &Temp{}
	return temp, nil
}

// Temp is a temporary state of rotation.
type Temp struct {
	*Base
}

// Compress archives current temp file.
func (temp *Temp) Compress() (*Archive, error) {
	archive := &Archive{}
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
func RunRotate(src *Source) error {
	temp, everr := src.Evacuate()
	if everr != nil {
		return everr
	}
	archive, cmerr := temp.Compress()
	if cmerr != nil {
		return cmerr
	}
	return archive.Finalize()
}
