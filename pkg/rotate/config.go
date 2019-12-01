package rotate

import "time"

// Config stores config
type Config struct {
	OwnerProcName        string
	OwnerReleaseInterval time.Duration
	OwnerReleaseTimeout  time.Duration
	TempStorage          string
	ArchiveStorage       string
	FinalizeCommands     []string
}

// NewConfig build config with default parameters.
func DefaultConfig() *Config {
	config := &Config{
		OwnerProcName:        "",
		OwnerReleaseInterval: 1 * time.Second,
		OwnerReleaseTimeout:  5 * time.Minute,
		TempStorage:          "/tmp/rotate/tmp",
		ArchiveStorage:       "/tmp/rotate/archive",
		FinalizeCommands:     []string{},
	}
	return config
}

func NewConfig(
	owner string,
	ownerReleaseTimeout time.Duration,
	tempStorage string,
	archiveStorage string,
	finalizeCommands []string) *Config {
	config := DefaultConfig()
	config.OwnerProcName = owner
	config.OwnerReleaseTimeout = ownerReleaseTimeout
	config.TempStorage = tempStorage
	config.ArchiveStorage = archiveStorage
	config.FinalizeCommands = finalizeCommands
	return config
}
