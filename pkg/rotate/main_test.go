package rotate

import (
	"fmt"
	"testing"
	"time"

	"github.com/utky/logproc-go/pkg/log"
)

type DummyFile struct {
	p string
}

func (f *DummyFile) Basename() string {
	return "test"
}
func (f *DummyFile) AbsolutePath() string {
	return "/tmp/rotate/test"
}
func (f *DummyFile) MoveTo(dest string) error {
	f.p = dest
	return nil
}

type OwnerReleaseImmediately struct {
}

func (this *OwnerReleaseImmediately) NotifyRelease(file File) error {
	return nil
}
func (this *OwnerReleaseImmediately) Released(file File) (bool, error) {
	return true, nil
}

func testConfig() *Config {
	config := DefaultConfig()
	config.OwnerReleaseTimeout = 2 * time.Millisecond
	config.OwnerReleaseInterval = 1 * time.Millisecond
	return config
}

func TestRotate(t *testing.T) {
	releaseNow := []Owner{&OwnerReleaseImmediately{}}
	logFile := &DummyFile{}
	source := &Source{
		Base:   &Base{testConfig(), logFile, log.New()},
		owners: releaseNow,
	}
	err := RunRotate(source)
	if err != nil {
		t.Errorf("Failed to rotate: %s", err)
	}
}

type OwnerFailNotify struct {
}

func (this *OwnerFailNotify) NotifyRelease(file File) error {
	return fmt.Errorf("Failed to send SIGHUP")
}
func (this *OwnerFailNotify) Released(file File) (bool, error) {
	return true, nil
}

func TestPreventByFailNotify(t *testing.T) {
	prevention := []Owner{&OwnerFailNotify{}}
	logFile := &DummyFile{}
	source := &Source{
		Base:   &Base{testConfig(), logFile, log.New()},
		owners: prevention,
	}
	err := RunRotate(source)
	if err == nil {
		t.Errorf("Accidentally succeeded")
	}
	if err.Error() != "Failed to send SIGHUP" {
		t.Errorf("Failed by other reason: %s", err)
	}
}

type OwnerDoesNotRelease struct {
}

func (this *OwnerDoesNotRelease) NotifyRelease(file File) error {
	return nil
}
func (this *OwnerDoesNotRelease) Released(file File) (bool, error) {
	return false, nil
}

func TestPreventByReleaseTimeout(t *testing.T) {
	prevention := []Owner{&OwnerDoesNotRelease{}}
	logFile := &DummyFile{}
	source := &Source{
		Base:   &Base{testConfig(), logFile, log.New()},
		owners: prevention,
	}
	err := RunRotate(source)
	if err == nil {
		t.Errorf("Accidentally succeeded")
	}
}

func TestWaitOwnerRelease(t *testing.T) {
	prevention := []Owner{&OwnerDoesNotRelease{}}
	logFile := &DummyFile{}
	source := &Source{
		Base:   &Base{testConfig(), logFile, log.New()},
		owners: prevention,
	}
	err := WaitOwnerRelease(source)
	if err == nil {
		t.Errorf("Timeout was not happened")
	}
}
