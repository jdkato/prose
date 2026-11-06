.PHONY: all test race bench lint fmt vet tidy cover clean models

all: test

# Run the test suite.
test:
	go test -count=1 ./...

# The taggers, segmenters and extracters are meant to be shared across
# goroutines; this is what proves it.
race:
	go test -race -count=1 ./...

bench:
	go test -bench=. -benchmem -run=^$$ ./...

lint:
	golangci-lint run

fmt:
	gofmt -w .

vet:
	go vet ./...

tidy:
	go mod tidy

cover:
	go test -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

clean:
	rm -f coverage.out

# Retrain every language model from Universal Dependencies treebanks checked
# out under $(UD); see tag/lang/build.sh.
models:
	UD=$(UD) tag/lang/build.sh
