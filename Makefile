# ---------------------------------------------------------
# Variables
# ---------------------------------------------------------

# Version info for binaries
GIT_REVISION := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

# Compiler flags
VPREFIX := github.com/utky/logproc-go/version
GO_LDFLAGS   := -s -w -X $(VPREFIX).Branch=$(GIT_BRANCH) -X $(VPREFIX).Revision=$(GIT_REVISION)
GO_FLAGS     := -ldflags "-extldflags \"-static\" $(GO_LDFLAGS)" -tags netgo


# ---------------------------------------------------------
# mgld
# ---------------------------------------------------------
.PHONY: rotate
APP_ROTATE := cmd/rotate/rotate

$(APP_ROTATE): cmd/rotate/main.go pkg/**/*.go
	CGO_ENABLED=0 go build $(GO_FLAGS) -o $@ ./$(@D)

# ---------------------------------------------------------
# Gobal
# ---------------------------------------------------------
.PHONY: clean test
test:
	go test ./...

APPS := $(APP_ROTATE)
clean:
	rm -f $(APPS)