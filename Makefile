# YOLO CLI Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=yolo
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date +%FT%T%z)

# Build directory
BUILD_DIR=bin
DIST_DIR=dist

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Supported platforms for cross-compilation
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

.PHONY: all build build-all clean test coverage deps lint install uninstall

all: clean deps build test

build:
	@echo "Building YOLO CLI..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/yolo/main.go

build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	$(foreach platform,$(PLATFORMS),\
		$(eval GOOS=$(word 1,$(subst /, ,$(platform)))) \
		$(eval GOARCH=$(word 2,$(subst /, ,$(platform)))) \
		$(eval BINARY=$(BINARY_NAME)_$(GOOS)_$(GOARCH)) \
		echo "Building $(BINARY)..." && \
		GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY) cmd/yolo/main.go || exit 1 && \
		tar -czf $(DIST_DIR)/$(BINARY).tar.gz -C $(DIST_DIR) $(BINARY) && \
		rm $(DIST_DIR)/$(BINARY) ; \
	)

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	$(GOCLEAN)

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

coverage:
	@echo "Generating test coverage report..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@rm coverage.out

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi

install: build
	@echo "Installing YOLO CLI..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "YOLO CLI installed successfully!"

uninstall:
	@echo "Uninstalling YOLO CLI..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "YOLO CLI uninstalled successfully!"

# Development helpers
dev: build
	@./$(BUILD_DIR)/$(BINARY_NAME)

watch:
	@if command -v air >/dev/null; then \
		air; \
	else \
		echo "air not installed. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Help target
help:
	@echo "YOLO CLI Makefile targets:"
	@echo "  make              - Clean, download dependencies, build, and test"
	@echo "  make build        - Build for current platform"
	@echo "  make build-all    - Build for all supported platforms"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make coverage     - Generate test coverage report"
	@echo "  make deps         - Download dependencies"
	@echo "  make lint         - Run linters"
	@echo "  make install      - Install YOLO CLI to /usr/local/bin"
	@echo "  make uninstall    - Remove YOLO CLI from /usr/local/bin"
	@echo "  make dev          - Build and run for development"
	@echo "  make watch        - Run with hot reload (requires air)"
	@echo "  make help         - Show this help message"
