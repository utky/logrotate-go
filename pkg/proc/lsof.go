package proc

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// LsofOut is structured output of lsof.
// Example:
//   COMMAND    PID      USER   FD   TYPE DEVICE SIZE/OFF    NODE NAM
type LsofOut struct {
	Command string
	Pid     int
	User    string
	Fd      string
	Typ     string
	Device  string
	Size    string
	Inode   int64
	Name    string
}

// ParseLsofOutput read lines and build structured output of lsof.
func ParseLsofOutput(lines []string) ([]LsofOut, error) {
	var outputList []LsofOut
	for i, line := range lines {

		// skip header
		if i == 0 {
			continue
		}

		fields := strings.Fields(line)
		numOfFields := len(fields)

		// skip empty line
		if numOfFields == 0 {
			continue
		}

		if numOfFields != 9 {
			return outputList, fmt.Errorf("Field number mismatch at line:%d expected %d actual %d : %s", i, 9, numOfFields, fields)
		}

		pid, pderr := strconv.ParseInt(fields[1], 10, 64)
		if pderr != nil {
			return outputList, fmt.Errorf("Cannot parse to int as pid at line:%d actual %s", i, fields[1])
		}

		inode, inerr := strconv.ParseInt(fields[7], 10, 64)
		if inerr != nil {
			return outputList, fmt.Errorf("Cannot parse to int64 as inode at line:%d actual %s", i, fields[7])
		}

		lsofOut := LsofOut{
			Command: fields[0],
			Pid:     int(pid),
			User:    fields[2],
			Fd:      fields[3],
			Typ:     fields[4],
			Device:  fields[5],
			Size:    fields[6],
			Inode:   inode,
			Name:    fields[8],
		}
		outputList = append(outputList, lsofOut)
	}
	return outputList, nil
}

// Lsof spawns lsof and parse output
func Lsof(path string, opts ...string) ([]LsofOut, error) {
	arguments := []string{"-w", path}
	arguments = append(arguments, opts...)
	cmd := exec.Command("lsof", arguments...)
	out, err := cmd.Output()
	if err != nil {
		return make([]LsofOut, 0), err
	}
	lines := strings.Split(string(out), "\n")
	return ParseLsofOutput(lines)
}
