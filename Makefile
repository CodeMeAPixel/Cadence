## Makefile for Cadence (common developer tasks)
## Cross-platform build with automatic version injection from git tags

.PHONY: all build install test fmt tidy lint vet run clean help

BINARY := cadence
OUT := bin/$(BINARY)

# Detect OS - check for Windows first (more reliable on Windows systems)
ifdef OS
    # Windows (native make or mingw)
    OS_TYPE := windows
    BINARY := cadence.exe
    OUT := cadence.exe
else
    # Unix-like systems
    UNAME_S := $(shell uname -s 2>/dev/null)
    ifeq ($(UNAME_S),Darwin)
        OS_TYPE := darwin
    else
        OS_TYPE := unix
    endif
endif

# Unix/Linux/macOS build with version injection
ifneq ($(OS_TYPE),windows)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -ldflags="-X github.com/codemeapixel/cadence/internal/version.Version=$(VERSION) -X github.com/codemeapixel/cadence/internal/version.GitCommit=$(COMMIT) -X github.com/codemeapixel/cadence/internal/version.BuildTime=$(BUILD_TIME)"
endif

all: build

# Platform-specific build targets
build:
ifeq ($(OS_TYPE),windows)
	@powershell -ExecutionPolicy Bypass -File .\scripts\build.ps1
else
	@mkdir -p bin
	@echo "Building Cadence..."
	@echo "  Version: $(VERSION)"
	@echo "  Commit:  $(COMMIT)"
	@echo "  Time:    $(BUILD_TIME)"
	@go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/cadence
	@echo "Build complete: bin/$(BINARY)"
endif

install:
ifeq ($(OS_TYPE),windows)
	@powershell -ExecutionPolicy Bypass -File .\scripts\build.ps1 -Install
else
	go install $(LDFLAGS) ./cmd/cadence
endif

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

lint:
	@golangci-lint run || true

vet:
	go vet ./...

run:
	go run ./cmd/cadence

clean:
	@rm -rf bin

help:
	@echo "Cadence Makefile targets:"
	@echo "  make build   - Build binary with automatic version injection from git tags"
	@echo "  make install - Install binary to \$$GOPATH/bin"
	@echo "  make test    - Run all tests"
	@echo "  make fmt     - Format code"
	@echo "  make tidy    - Tidy dependencies"
	@echo "  make lint    - Run linter"
	@echo "  make vet     - Run go vet"
	@echo "  make run     - Run application"
	@echo "  make clean   - Clean build artifacts"

