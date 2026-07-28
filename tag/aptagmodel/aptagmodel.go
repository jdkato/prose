// Package aptagmodel holds the embedded English averaged-perceptron model.
//
// It is a data leaf on purpose: it exports bytes and imports nothing. Keeping
// it dependency-free means a consumer that never tags text never links the
// model in, and the ~3 MB stays out of their binary entirely. Enforced by
// TestTaggerDoesNotLinkForeignModels in the parent package.
//
// The bytes are in the flat format defined by prose's internal/flat package,
// which prose/v3/tag knows how to read.
package aptagmodel

import _ "embed"

// English is the built-in English model, in prose's flat tagger format.
//
// It is read-only. Callers must not modify it: prose/v3/tag slices feature
// strings directly out of these bytes rather than copying them.
//
//go:embed en.bin
var English []byte
