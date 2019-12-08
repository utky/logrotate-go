package core

import "time"

// Config stores config
type Config struct {
	OwnerProcName        string
	OwnerReleaseInterval time.Duration
	OwnerReleaseTimeout  time.Duration
	SourcePattern        string
	TempStorage          string
	ArchiveStorage       string
	FinalizeCommands     []string
}

// DefaultConfig build config with default parameters.
func DefaultConfig() *Config {
	config := &Config{
		OwnerProcName:        "",
		OwnerReleaseInterval: 1 * time.Second,
		OwnerReleaseTimeout:  5 * time.Minute,
		SourcePattern:        "/tmp/rotate/source",
		TempStorage:          "/tmp/rotate/tmp",
		ArchiveStorage:       "/tmp/rotate/archive",
		FinalizeCommands:     []string{},
	}
	return config
}

// NewConfig builds configuration.
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
