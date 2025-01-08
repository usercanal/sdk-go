VERSION ?= dev
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -X grpc-simple/version.version=$(VERSION) \
          -X grpc-simple/version.commitHash=$(COMMIT_HASH) \
          -X grpc-simple/version.buildTime=$(BUILD_TIME)

.PHONY: build
build:
    go build -ldflags "$(LDFLAGS)" ./...

.PHONY: install
install:
    go install -ldflags "$(LDFLAGS)" ./...

.PHONY: release
release:
    GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/usercanal-linux-amd64 ./...
    GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/usercanal-darwin-amd64 ./...
    GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/usercanal-windows-amd64.exe ./...