.PHONY: build test lint clean install version

# Build variables
BINARY_NAME=zai
BUILD_DIR=bin
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0-dev")
BUILD_TIME?=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT?=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# Build ldflags
LDFLAGS=-X 'zai/internal/version.Version=$(VERSION)' \
		-X 'zai/internal/version.Build=$(BUILD_TIME)' \
		-X 'zai/internal/version.Commit=$(COMMIT)'

build: lint
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	@echo "Running tests..."
	@go test -v -race -cover ./...

lint:
	@echo "Running linters..."
	@golangci-lint run ./...

install: build
	@echo "Installing $(BINARY_NAME)..."
	@ln -sf $(PWD)/$(BUILD_DIR)/$(BINARY_NAME) $(HOME)/go/bin/$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@rm -f $(BINARY_NAME)

version:
	@echo "Version: $(VERSION)"
	@echo "Build:   $(BUILD_TIME)"
	@echo "Commit:  $(COMMIT)"
