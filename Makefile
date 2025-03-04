BINARY_NAME=hetznerdns
VERSION=0.1.0
MODULE=github.com/shotgundd/hetznerdns
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"
BUILD_DIR=build

.PHONY: all build clean install uninstall test test-unit test-integration test-coverage

all: clean build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/hetznerdns

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out

install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) ./cmd/hetznerdns

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(shell go env GOPATH)/bin/hetznerdns

# Test targets
test: test-unit test-integration

test-unit:
	@echo "Running unit tests..."
	@go test -v ./pkg/...

test-integration:
	@echo "Running integration tests..."
	@go test -v ./cmd/...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./pkg/... ./cmd/...
	@go tool cover -html=coverage.out

# Cross-compilation targets
.PHONY: build-all build-linux build-windows build-macos

build-all: build-linux build-windows build-macos

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/hetznerdns

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/hetznerdns

build-macos:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/hetznerdns 