.PHONY: build test lint clean

build:
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size

test:
	go test -v .

test-coverage:
	go test -coverprofile=coverage.out .
	go tool cover -html=coverage.out

lint:
	golangci-lint run

run-example:
	./bin/hexlet-path-size data.csv

run-test-example:
	./bin/hexlet-path-size data.csv
	./bin/hexlet-path-size -h data.csv
	./bin/hexlet-path-size -a .

clean:
	rm -rf bin/ coverage.out

all: build test