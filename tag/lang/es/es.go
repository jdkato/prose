// Package es registers the Spanish part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/es"
//
//	tagger, err := tag.ForLanguage("es")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Spanish-AnCora@3ea1f74 (2026-11-10); UD_Spanish-GSD@1b694f5 (2026-11-10)
// (license: CC BY 4.0, CC BY-SA 4.0), and scores 97.74% on that treebank's development
// section and 97.68% on its test section. See MODELS.md for how.
package es

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed es.bin
var model []byte

func init() { tag.RegisterLanguage("es", model) }
