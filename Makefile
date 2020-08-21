BASE_DIR=$(shell echo $$GOPATH)/src/github.com/jdkato/prose
BUILD_DIR=./builds

LDFLAGS=-ldflags "-s -w"

.PHONY: clean test lint ci cross install bump model setup

all: build

build:
	go build ${LDFLAGS} -o bin/prose ./cmd/prose

build-win:
	go build ${LDFLAGS} -o bin/prose.exe ./cmd/prose

bench:
	go test -bench=. -run=^$$ -benchmem

test:
	go test -v

ci: lint test

lint:
	./bin/golangci-lint run

model:
	go-bindata -ignore=\\.DS_Store -pkg="prose" -o data.go model/**/*.gob
