package core

import "github.com/utky/logproc-go/pkg/log"

// Entry is common structure which all stages of log should have.
type Entry struct {
	Config *Config
	File   File
	Logger *log.Logger
}
