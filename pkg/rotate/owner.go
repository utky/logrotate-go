package rotate

import (
	"os"
	"syscall"
)

// Owner handles request to release handler of target log file.
type Owner interface {
	NotifyRelease(file File) error
	// Released queries owner
	Released(file File) (bool, error)
}

// ProcessOwner is a owner represented by process.
type ProcessOwner struct {
	process *os.Process
}

// NotifyRelease sends SIGHUP to the process.
func (owner *ProcessOwner) NotifyRelease(file File) error {
	return owner.process.Signal(syscall.SIGHUP)
}

// Released queries if specified file is owned by the process
func (owner *ProcessOwner) Released(file File) (bool, error) {
	opts := []string{"-p", string(owner.process.Pid)}
	procs, err := Lsof(file.AbsolutePath(), opts...)
	if err != nil {
		return false, err
	}
	return len(procs) < 1, nil
}

// NotifyAll sends signal specified owners.
func NotifyAll(owners []Owner, file File) error {
	for _, owner := range owners {
		notifyErr := owner.NotifyRelease(file)
		if notifyErr != nil {
			return notifyErr
		}
	}
	return nil
}

// AllOwnersReleased queries current handle of file.
func AllOwnersReleased(owners []Owner, file File) (bool, error) {
	var released bool
	var err error
	released = true
	err = nil
	for _, owner := range owners {
		released, err = owner.Released(file)
		if !released || err != nil {
			break
		}
	}
	return released, err
}
