
# Makefile for Ethereum Transaction Parser

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=build/eth-tx-parser

# Directories
CMD_DIR=./cmd
INTERNAL_DIR=./internal

# Build target
build:
	@echo "Building the application..."
	$(GOBUILD) -o $(BINARY_NAME) $(CMD_DIR)/*.go
	@echo "Build completed successfully"

# Test target
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...
	@echo "Testing completed"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@if [ -d build ]; then rm -rf build; fi
	rm -f $(BINARY_NAME)
	$(GOCMD) clean -cache
	@echo "Clean up completed"

# Install dependencies
deps:
	@echo "Fetching dependencies..."
	$(GOGET) -v ./...
	@echo "Dependencies installed"
	
fmt:
	@echo "Formatting Go code..."
	@gofmt -w .
	@echo "Formatting completed"

# Install the application
install:
	@echo "Installing application..."
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "Application installed to $(GOPATH)/bin"

# Help target
help:
	@echo "Available targets:"
	@echo "  build   - Compile the application"
	@echo "  test    - Run tests"
	@echo "  clean   - Remove build artifacts"
	@echo "  deps    - Install dependencies"
	@echo "  install - Install the application globally"
	@echo "  help    - Show this help message"

# Default target
.DEFAULT_GOAL := help
