package rotate

import "testing"

func TestPgrep(t *testing.T) {
	pids, err := Pgrep("/sbin/init")
	if err != nil {
		t.Errorf("Failed to run pgrep: %s", err)
	}
	if len(pids) == 0 {
		t.Errorf("No process matches with %s", "/sbin/init")
	}
	if len(pids) != 1 {
		t.Errorf("Other process matched than init %v", pids)
	}
}
