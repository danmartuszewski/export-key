.PHONY: build test test-coverage lint fmt clean install deps run

BINARY_NAME=export-key
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-ldflags "-s -w \
	-X github.com/danmartuszewski/export-key/internal/cmd.Version=$(VERSION) \
	-X github.com/danmartuszewski/export-key/internal/cmd.Commit=$(COMMIT) \
	-X github.com/danmartuszewski/export-key/internal/cmd.BuildDate=$(BUILD_DATE)"

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/export-key

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

clean:
	rm -rf bin/ dist/ coverage.out coverage.html

install: build
	cp bin/$(BINARY_NAME) /usr/local/bin/
	codesign --sign - --force /usr/local/bin/$(BINARY_NAME)

deps:
	go mod download
	go mod tidy

run:
	go run ./cmd/export-key $(ARGS)
