package evacuate

import (
	"fmt"
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

func TestRun(t *testing.T) {
	releaseNow := []Owner{&OwnerReleaseImmediately{}}
	logFile := &DummyFile{}
	source := &Source{
		Entry:  &core.Entry{testConfig(), logFile, log.New()},
		owners: releaseNow,
	}
	err := Run(source)
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
	err := Run(source)
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
	err := Run(source)
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
