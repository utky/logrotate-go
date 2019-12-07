package proc

import "testing"

func TestPgrep(t *testing.T) {
	pids, err := Pgrep("go")
	if err != nil {
		t.Errorf("Failed to run pgrep: %s", err)
	}
	if len(pids) == 0 {
		t.Errorf("No process matches with %s", "go")
	}
}
