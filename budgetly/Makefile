# Project variables
APP_NAME = budgetly
BUILD_DIR = ./bin
SRC_DIR = ./cmd/main.go
BASE_DIR = $(shell pwd)

# Set the environment file
ENV_FILE = ./bin/.env

# Set the default Go commands
GO_CMD = go
GO_BUILD = $(GO_CMD) build
GO_RUN = $(GO_CMD) run
GO_TEST = $(GO_CMD) test
GO_FMT = $(GO_CMD) fmt
GO_LINT = golangci-lint run
GO_CLEAN = $(GO_CMD) clean

# Set the Go output binary name
OUTPUT_BINARY = $(BUILD_DIR)/$(APP_NAME)

# Targets
.PHONY: all build run test clean fmt lint

# Default target: Build the project
all: build

# Build the project
build:
	@mkdir -p $(BUILD_DIR)
	@$(GO_BUILD) -o $(OUTPUT_BINARY) $(SRC_DIR)

# Run the project
run: build
	@$(OUTPUT_BINARY)

# Test the project
test:
	@$(GO_TEST) -v ./...

# Format the Go files
fmt:
	@$(GO_FMT) ./...

# Lint the Go files (assumes golangci-lint is installed)
lint:
	@$(GO_LINT)

# Clean up the build artifacts
clean:
	@rm $(BUILD_DIR)/$(APP_NAME)
	@$(GO_CLEAN)
