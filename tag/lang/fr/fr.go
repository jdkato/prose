// Package fr registers the French part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/fr"
//
//	tagger, err := tag.ForLanguage("fr")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_French-GSD@94d5b68 (2026-05-06)
// (license: CC BY-SA 4.0), and scores 97.00% on that treebank's development
// section and 96.67% on its test section. See MODELS.md for how.
package fr

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed fr.bin
var model []byte

func init() { tag.RegisterLanguage("fr", model) }
