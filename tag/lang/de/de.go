// Package de registers the German part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/de"
//
//	tagger, err := tag.ForLanguage("de")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_German-GSD@ce54dbe (2026-05-06)
// (license: CC BY-SA 4.0), and scores 94.16% on that treebank's development
// section and 93.93% on its test section. See MODELS.md for how.
package de

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed de.bin
var model []byte

func init() { tag.RegisterLanguage("de", model) }
