BASE_DIR=$(shell echo $$GOPATH)/src/github.com/jdkato/prose
BUILD_DIR=./builds

LDFLAGS=-ldflags "-s -w"

.PHONY: clean test lint ci cross install bump model setup vendor

all: build

build:
	go build ${LDFLAGS} -o bin/prose ./cmd/prose

build-win:
	go build ${LDFLAGS} -o bin/prose.exe ./cmd/prose

bench:
	go test -bench=. ./tokenize ./transform ./summarize ./tag ./chunk

test-tokenize:
	go test -v ./tokenize

test-transform:
	go test -v ./transform

test-summarize:
	go test -v ./summarize

test-chunk:
	go test -v ./chunk

test-tag:
	go test -v ./tag

test: test-tokenize test-transform test-summarize test-chunk test-tag

ci: test lint

lint:
	golangci-lint run

setup:
	go get -u github.com/shogo82148/go-shuffle
	go get -u github.com/jdkato/syllables
	go get -u github.com/montanaflynn/stats
	go get -u gopkg.in/neurosnap/sentences.v1/english
	go get -u github.com/stretchr/testify/assert
	go get -u github.com/urfave/cli
	go get -u github.com/jteeuwen/go-bindata/...
	go-bindata -ignore=\\.DS_Store -pkg="model" -o internal/model/model.go internal/model/
	wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.19.1

model:
	go-bindata -ignore=\\.DS_Store -pkg="model" -o internal/model/model.go internal/model/*.gob

vendor:
	go mod tidy && go mod vendor
