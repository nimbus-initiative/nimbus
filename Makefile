.PHONY: all build-nimbusd build-nimbusctl install clean test lint

# Version information
VERSION ?= 0.1.0
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Directories
BIN_DIR := bin
CMD_DIR := cmd
INTERNAL_DIR := internal

# Binaries
NIMBUSD_BIN := $(BIN_DIR)/nimbusd
NIMBUSCTL_BIN := $(BIN_DIR)/nimbusctl

# Build flags
GO_BUILD_FLAGS := -ldflags "$(LDFLAGS)"

# Default target
all: build-nimbusd build-nimbusctl

# Build nimbusd
build-nimbusd:
	@echo "Building nimbusd..."
	@mkdir -p $(BIN_DIR)
	@go build $(GO_BUILD_FLAGS) -o $(NIMBUSD_BIN) ./$(CMD_DIR)/nimbusd

# Build nimbusctl
build-nimbusctl:
	@echo "Building nimbusctl..."
	@mkdir -p $(BIN_DIR)
	@go build $(GO_BUILD_FLAGS) -o $(NIMBUSCTL_BIN) ./$(CMD_DIR)/nimbusctl

# Install nimbusd
install-nimbusd: build-nimbusd
	@echo "Installing nimbusd to /usr/local/bin/"
	@sudo install -m 755 $(NIMBUSD_BIN) /usr/local/bin/

# Install nimbusctl
install-nimbusctl: build-nimbusctl
	@echo "Installing nimbusctl to /usr/local/bin/"
	@sudo install -m 755 $(NIMBUSCTL_BIN) /usr/local/bin/

# Install both
install: install-nimbusd install-nimbusctl

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Generate code (mocks, protobufs, etc.)
generate:
	@echo "Generating code..."
	@go generate ./...

# Run in development mode
dev: build-nimbusd build-nimbusctl
	@echo "Starting nimbusd in development mode..."
	@sudo $(NIMBUSD_BIN) --config ./configs/nimbusd.toml

# Help
doc:
	@echo "Available targets:"
	@echo "  all              - Build all binaries (default)"
	@echo "  build-nimbusd    - Build nimbusd"
	@echo "  build-nimbusctl  - Build nimbusctl"
	@echo "  install-nimbusd  - Install nimbusd to /usr/local/bin/"
	@echo "  install-nimbusctl- Install nimbusctl to /usr/local/bin/"
	@echo "  install          - Install all binaries"
	@echo "  clean            - Remove build artifacts"
	@echo "  test             - Run tests"
	@echo "  lint             - Run linter"
	@echo "  generate         - Generate code"
	@echo "  dev              - Run in development mode"

# Default target
.DEFAULT_GOAL := all
