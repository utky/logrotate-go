package rotate

import (
	"testing"
)

func TestLsof(t *testing.T) {
	output, err := Lsof("/dev/null")
	if err != nil {
		t.Errorf("Failed to run and capture lsof against /dev/null: %s", err)
	}

	for _, o := range output {
		if o.Command == "" {
			t.Errorf("Command is empty: %v", o)
		}

		if o.Pid == 0 {
			t.Errorf("Pid is zero: %v", o)
		}

		if o.User == "" {
			t.Errorf("User is empty: %v", o)
		}

		if o.Fd == "" {
			t.Errorf("Fd is empty: %v", o)
		}

		if o.Typ == "" {
			t.Errorf("Typ is empty: %v", o)
		}

		if o.Device == "" {
			t.Errorf("Device is empty: %v", o)
		}

		if o.Size == "" {
			t.Errorf("Size is empty: %v", o)
		}

		if o.Inode < 1 {
			t.Errorf("Size is not positive number: %v", o)
		}

		if o.Name == "" {
			t.Errorf("Name is empty: %v", o)
		}
	}
}

func TestParseLsofOutput(t *testing.T) {
	lines := []string{
		"COMMAND    PID      USER   FD   TYPE DEVICE SIZE/OFF    NODE NAM",
		"code      455112 ilyaletre    0r   CHR    1,3      0t0 2051 /dev/null",
	}

	output, err := ParseLsofOutput(lines)
	if err != nil {
		t.Errorf("Failed to parse: %s", err)
	}
	if len(output) != 1 {
		t.Errorf("Number of parsed output was wrong, expected 1 actual %d : %v", len(output), output)
	}
	record := output[0]
	if record.Command != "code" {
		t.Errorf("Command: actual %s", record.Command)
	}
	if record.Pid != 455112 {
		t.Errorf("Pid: actual %d", record.Pid)
	}
	if record.User != "ilyaletre" {
		t.Errorf("User: actual %s", record.User)
	}
	if record.Fd != "0r" {
		t.Errorf("Fd: actual %s", record.Fd)
	}
	if record.Typ != "CHR" {
		t.Errorf("Typ: actual %s", record.Typ)
	}
	if record.Device != "1,3" {
		t.Errorf("Device: actual %s", record.Device)
	}
	if record.Size != "0t0" {
		t.Errorf("Size: actual %s", record.Size)
	}
	if record.Inode != 2051 {
		t.Errorf("Inode: actual %d", record.Inode)
	}
	if record.Name != "/dev/null" {
		t.Errorf("Name: actual %s", record.Name)
	}
}
