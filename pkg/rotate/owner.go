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

type ProcessOwner struct {
	process *os.Process
}

func (owner *ProcessOwner) NotifyRelease(file File) error {
	return owner.process.Signal(syscall.SIGHUP)
}

func (owner *ProcessOwner) Released(file File) (bool, error) {
	procs, err := Lsof(file.AbsolutePath())
	if err != nil {
		return false, err
	}
	return len(procs) < 1, nil
}
