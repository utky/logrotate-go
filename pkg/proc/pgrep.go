package proc

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Pgrep runs pgrep and returns list of pid.
func Pgrep(pattern string) ([]int, error) {
	pids := make([]int, 0)
	cmd := exec.Command("pgrep", "-f", pattern)
	outBytes, cmdErr := cmd.Output()
	if cmdErr != nil {
		return pids, cmdErr
	}
	pidStrs := strings.Split(string(outBytes), "\n")
	for _, pidStr := range pidStrs {
		if pidStr == "" {
			continue
		}
		pid, convErr := strconv.Atoi(pidStr)
		if convErr != nil {
			return pids, convErr
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

// GrepProcs collects process structs with pgrep
func GrepProcs(pattern string) ([]*os.Process, error) {
	procs := make([]*os.Process, 0)
	pids, pgrepErr := Pgrep(pattern)
	if pgrepErr != nil {
		return procs, pgrepErr
	}
	for _, pid := range pids {
		proc, findErr := os.FindProcess(pid)
		if findErr != nil {
			return procs, findErr
		}
		procs = append(procs, proc)
	}
	return procs, nil
}
