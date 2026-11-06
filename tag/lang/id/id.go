// Package id registers the Indonesian part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/id"
//
//	tagger, err := tag.ForLanguage("id")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Indonesian-GSD@bc23d75 (2026-05-06)
// (license: CC BY-SA 4.0), and scores 93.38% on that treebank's development
// section and 93.59% on its test section. See MODELS.md for how.
package id

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed id.bin
var model []byte

func init() { tag.RegisterLanguage("id", model) }
