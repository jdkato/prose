// Package es registers the Spanish part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/es"
//
//	tagger, err := tag.ForLanguage("es")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Spanish-AnCora@20adddb (2026-05-06)
// (license: CC BY 4.0), and scores 98.03% on that treebank's development
// section and 97.98% on its test section. See MODELS.md for how.
package es

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed es.bin
var model []byte

func init() { tag.RegisterLanguage("es", model) }
