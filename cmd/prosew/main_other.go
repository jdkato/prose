//go:build !(js && wasm)

package main

import (
	"fmt"
	"os"
)

// The module only means something in a JavaScript host; this keeps the
// package building natively, where its tests run.
func main() {
	fmt.Fprintln(os.Stderr, "prosew is a WebAssembly module: build it with GOOS=js GOARCH=wasm")
	os.Exit(1)
}
