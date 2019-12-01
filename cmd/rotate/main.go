package main

import (
	"flag"
	"os"
	"time"

	"github.com/utky/logproc-go/pkg/log"
	"github.com/utky/logproc-go/pkg/rotate"
)

var logger = log.New()

type commands []string

func (cs *commands) String() string {
	return "hoge"
}

func (cs *commands) Set(value string) error {
	*cs = append(*cs, value)
	return nil
}

func parseOpts() (string, *rotate.Config) {
	sourceName := flag.String("src", "", "source file path to be processed")
	owner := flag.String("owner", "", "owner process name of src")
	ownerReleaseTimeout := flag.Duration("waitTimeout", 5*time.Minute, "timeout to wait release of owner")
	tempDirectory := flag.String("tempDir", "/tmp/rotate/tmp", "temp directory to rotate files")
	archiveDirectory := flag.String("archiveDir", "/tmp/rotate/archive", "archive directory to rotate files")
	cmds := make(commands, 0)
	flag.Var(&cmds, "finalize", "command to be passed to shell after rotation completed")
	flag.Parse()

	config := rotate.NewConfig(
		*owner,
		*ownerReleaseTimeout,
		*tempDirectory,
		*archiveDirectory,
		cmds,
	)

	return *sourceName, config
}

func abort(msg string, err error, exitCode int) {
	logger.Errorf(msg, log.Fields{"error": err})
	os.Exit(exitCode)
}

func main() {
	var config *rotate.Config
	var file rotate.File

	sourceName, config := parseOpts()
	if sourceName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	file, newFileErr := rotate.NewOsFile(sourceName)
	if newFileErr != nil {
		abort("Failed to resolve absolute path of src", newFileErr, 1)
	}
	procs, procErr := rotate.GrepProcs(config.OwnerProcName)
	if procErr != nil {
		abort("Failed to find pid", procErr, 1)
	}
	source := rotate.NewSource(config, file, rotate.NewProcOwnerList(procs))
	err := rotate.RunRotate(source)
	if err != nil {
		abort("Failed to rotate log", err, 1)
	}
}
