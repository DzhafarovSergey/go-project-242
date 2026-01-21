.PHONY: build test lint clean all help

BINARY_NAME=hexlet-path-size
BINARY_DIR=bin

build:
	go build -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/hexlet-path-size

run:
	go run cmd/hexlet-path-size/main.go $(ARGS)

test:
	go test -v -cover ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

clean:
	rm -rf $(BINARY_DIR)/ coverage.out coverage.html

all: lint test build

help:
	@echo "Available commands:"
	@echo "  make build          - Build the binary"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make lint           - Run linter"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make all            - Run lint, test, and build"