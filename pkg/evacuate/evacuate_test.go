package evacuate

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/utky/logproc-go/pkg/core"
	"github.com/utky/logproc-go/pkg/log"
)

type DummyFile struct {
	p string
}

func (f *DummyFile) Base() string {
	return "test"
}
func (f *DummyFile) Abs() string {
	return "/tmp/rotate/test"
}
func (f *DummyFile) Move(dest string) error {
	f.p = dest
	return nil
}

type OwnerReleaseImmediately struct {
}

func (this *OwnerReleaseImmediately) NotifyRelease(file core.File) error {
	return nil
}
func (this *OwnerReleaseImmediately) Released(file core.File) (bool, error) {
	return true, nil
}

func testConfig() *core.Config {
	config := core.DefaultConfig()
	config.OwnerReleaseTimeout = 2 * time.Millisecond
	config.OwnerReleaseInterval = 1 * time.Millisecond
	return config
}

func TestRunWith(t *testing.T) {
	releaseNow := []Owner{&OwnerReleaseImmediately{}}
	logFile := &DummyFile{}
	source := &Source{
		Entry:  &core.Entry{testConfig(), logFile, log.New()},
		owners: releaseNow,
	}
	err := RunWith(source)
	if err != nil {
		t.Errorf("Failed to rotate: %s", err)
	}
}

type OwnerFailNotify struct {
}

func (this *OwnerFailNotify) NotifyRelease(file core.File) error {
	return fmt.Errorf("Failed to send SIGHUP")
}
func (this *OwnerFailNotify) Released(file core.File) (bool, error) {
	return true, nil
}

func TestPreventByFailNotify(t *testing.T) {
	prevention := []Owner{&OwnerFailNotify{}}
	logFile := &DummyFile{}
	source := &Source{
		Entry:  &core.Entry{testConfig(), logFile, log.New()},
		owners: prevention,
	}
	err := RunWith(source)
	if err == nil {
		t.Errorf("Accidentally succeeded")
	}
	if err.Error() != "Failed to send SIGHUP" {
		t.Errorf("Failed by other reason: %s", err)
	}
}

type OwnerDoesNotRelease struct {
}

func (this *OwnerDoesNotRelease) NotifyRelease(file core.File) error {
	return nil
}
func (this *OwnerDoesNotRelease) Released(file core.File) (bool, error) {
	return false, nil
}

func TestPreventByReleaseTimeout(t *testing.T) {
	prevention := []Owner{&OwnerDoesNotRelease{}}
	logFile := &DummyFile{}
	source := &Source{
		Entry:  &core.Entry{testConfig(), logFile, log.New()},
		owners: prevention,
	}
	err := RunWith(source)
	if err == nil {
		t.Errorf("Accidentally succeeded")
	}
}

func TestWaitOwnerRelease(t *testing.T) {
	prevention := []Owner{&OwnerDoesNotRelease{}}
	logFile := &DummyFile{}
	source := &Source{
		Entry:  &core.Entry{testConfig(), logFile, log.New()},
		owners: prevention,
	}
	err := WaitOwnerRelease(source)
	if err == nil {
		t.Errorf("Timeout was not happened")
	}
}

func TestRun(t *testing.T) {
	d := core.CreateTemp(t)
	defer os.RemoveAll(d)

	source := path.Join(d, "source")
	if err := os.MkdirAll(source, 0777); err != nil {
		t.Errorf("Failed to create source dir: %s", err)
	}
	temp := path.Join(d, "temp")
	if err := os.MkdirAll(temp, 0777); err != nil {
		t.Errorf("Failed to create temp dir: %s", err)
	}
	core.PrepareFile(t, source, "test.log", "TestRun")

	config := core.DefaultConfig()
	config.SourcePattern = path.Join(source, "*")
	config.TempStorage = temp

	failed, runErr := Run(config)
	if runErr != nil {
		t.Errorf("Failed to evacuate: %s", runErr)
	}
	if len(failed) > 0 {
		t.Errorf("Failed entry exists: %s", failed)
	}

	if _, absentErr := os.Stat(path.Join(temp, "test.log")); os.IsNotExist(absentErr) {
		t.Errorf("Cannot find moved file.: %s", absentErr)
	}
	//if _, absentErr := os.Stat(path.Join(source, "test.log")); !os.IsNotExist(absentErr) {
	//	t.Errorf("Cannot find ioved file.: %s", absentErr)
	//}
}
