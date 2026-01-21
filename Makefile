.PHONY: build test lint

build:
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size

test:
	go test ./tests/...

test-coverage:
	go test -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out

lint:
	golangci-lint run


run-example:
	./bin/hexlet-path-size data.csv

clean:
	rm -rf bin/ coverage.out

all: build test