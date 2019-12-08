package evacuate

import (
	"strconv"
	"strings"
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
	namePattern string
}

// NotifyRelease sends SIGHUP to the process.
func (owner *ProcessOwner) NotifyRelease(file core.File) error {
	procs, grepErr := proc.GrepProcs(owner.namePattern)
	if grepErr != nil {
		return grepErr
	}
	for _, p := range procs {
		if sigErr := p.Signal(syscall.SIGHUP); sigErr != nil {
			return sigErr
		}
	}
	return nil
}

// Released queries if specified file is owned by the process
func (owner *ProcessOwner) Released(file core.File) (bool, error) {
	owners, grepErr := proc.GrepProcs(owner.namePattern)
	if grepErr != nil {
		return false, grepErr
	}

	ownerPids := make([]string, len(owners))
	for i, o := range owners {
		ownerPids[i] = strconv.Itoa(o.Pid)
	}
	cvPids := strings.Join(ownerPids, ",")

	opts := []string{"-p", cvPids}
	procs, err := proc.Lsof(file.Abs(), opts...)
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

// NewProcOwner transform process list to owner list.
func NewProcOwner(config *core.Config) Owner {
	return &ProcessOwner{
		namePattern: config.OwnerProcName,
	}
}
