.PHONY: build test lint clean all help

BINARY_NAME=hexlet-path-size
BINARY_DIR=bin

build:
	@echo "Building binary..."
	go build -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/hexlet-path-size
	@echo "Binary built: $(BINARY_DIR)/$(BINARY_NAME)"

run:
	@go run cmd/hexlet-path-size/main.go $(ARGS)

test:
	@echo "Running tests..."
	go test -v -cover ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	@echo "Running linter..."
	golangci-lint run

clean:
	@echo "Cleaning up..."
	rm -rf $(BINARY_DIR)/ coverage.out coverage.html

all: lint test build
	@echo "All tasks completed successfully!"

help:
	@echo "Available commands:"
	@echo "  make build          - Build the binary"
	@echo "  make run            - Run the application (use ARGS='...' for arguments)"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with HTML coverage report"
	@echo "  make lint           - Run linter"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make all            - Run lint, test, and build"
	@echo ""
	@echo "Examples:"
	@echo "  make run ARGS=\"-H file.txt\""
	@echo "  make run ARGS=\"-raH ./directory\""