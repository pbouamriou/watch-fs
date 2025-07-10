# Makefile for watch-fs
.PHONY: help build clean deps test lint release install

# Variables
BINARY_NAME=watch-fs
BUILD_DIR=bin
DIST_DIR=dist
VERSION=$(shell git describe --tags --always --dirty)

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/watch-fs
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "Clean complete"

deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

test: ## Run tests
	@echo "Running tests..."
	go test -v ./test/...
	@echo "Tests complete"

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out -covermode=atomic ./test/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run
	@echo "Lint complete"

lint-install: ## Install golangci-lint
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2
	@echo "golangci-lint installed"

install: build ## Install the application
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

uninstall: ## Uninstall the application
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstallation complete"

release-build: ## Build for multiple platforms
	@echo "Building for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/watch-fs
	GOOS=linux GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/watch-fs
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/watch-fs
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/watch-fs
	
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/watch-fs
	GOOS=windows GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o $(DIST_DIR)/$(BINARY_NAME)-windows-arm64.exe ./cmd/watch-fs
	
	@echo "Multi-platform build complete"

release: test lint release-build ## Prepare release (test, lint, build)
	@echo "Release preparation complete"
	@echo "Binaries available in $(DIST_DIR)/"

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) -path .

dev: ## Run in development mode
	@echo "Running in development mode..."
	go run ./cmd/watch-fs -path .

# Default target
.DEFAULT_GOAL := help 