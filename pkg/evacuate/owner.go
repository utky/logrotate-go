package evacuate

import (
	"os"
	"syscall"

	"github.com/utky/logproc-go/pkg/core"
	"github.com/utky/logproc-go/pkg/proc"
)

// Owner handles request to release handler of target log file.
type Owner interface {
	NotifyRelease(file core.File) error
	// Released queries owner
	Released(file core.File) (bool, error)
}

// ProcessOwner is a owner represented by process.
type ProcessOwner struct {
	process *os.Process
}

// NotifyRelease sends SIGHUP to the process.
func (owner *ProcessOwner) NotifyRelease(file core.File) error {
	return owner.process.Signal(syscall.SIGHUP)
}

// Released queries if specified file is owned by the process
func (owner *ProcessOwner) Released(file core.File) (bool, error) {
	opts := []string{"-p", string(owner.process.Pid)}
	procs, err := proc.Lsof(file.Base(), opts...)
	if err != nil {
		return false, err
	}
	return len(procs) < 1, nil
}

// NotifyAll sends signal specified owners.
func NotifyAll(owners []Owner, file core.File) error {
	for _, owner := range owners {
		notifyErr := owner.NotifyRelease(file)
		if notifyErr != nil {
			return notifyErr
		}
	}
	return nil
}

// AllOwnersReleased queries current handle of file.
func AllOwnersReleased(owners []Owner, file core.File) (bool, error) {
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

// NewProcOwnerList transform process list to owner list.
func NewProcOwnerList(procs []*os.Process) []Owner {
	owners := make([]Owner, 0)
	for _, proc := range procs {
		owners = append(owners, &ProcessOwner{
			process: proc,
		})
	}
	return owners
}
