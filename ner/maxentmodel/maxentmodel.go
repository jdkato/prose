// Package maxentmodel holds the embedded English named-entity model.
//
// Like tag/aptagmodel, it is a data leaf: it exports bytes and imports
// nothing. That is what keeps the ~4 MB out of binaries that only tag text —
// asserted by tag/isolation_test.go, which fails if the tagger ever reaches
// this package.
//
// The bytes are in the flat format defined by prose's internal/maxent package,
// which prose/v3/ner knows how to read.
package maxentmodel

import _ "embed"

// English is the built-in English entity model, in prose's flat format.
//
// It is read-only. Callers must not modify it: prose/v3/ner slices feature
// keys directly out of these bytes rather than copying them.
//
//go:embed en.bin
var English []byte
