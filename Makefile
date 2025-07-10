.PHONY: build clean test run help

# Variables
BINARY_NAME=watch-fs
BUILD_DIR=bin

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run the application (example)
run:
	@echo "Running $(BINARY_NAME)..."
	go run main.go -path .

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Show help
help:
	@echo "Available commands:"
	@echo "  build  - Build the application"
	@echo "  clean  - Clean build artifacts"
	@echo "  test   - Run tests"
	@echo "  run    - Run the application (example with current directory)"
	@echo "  deps   - Install dependencies"
	@echo "  help   - Show this help message" 