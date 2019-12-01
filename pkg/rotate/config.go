package rotate

import "time"

// Config stores config
type Config struct {
	ownerReleaseInterval time.Duration
	ownerReleaseTimeout  time.Duration
}

// NewConfig build config with default parameters.
func NewConfig() *Config {
	config := &Config{
		ownerReleaseInterval: 1 * time.Second,
		ownerReleaseTimeout:  5 * time.Minute,
	}
	return config
}
