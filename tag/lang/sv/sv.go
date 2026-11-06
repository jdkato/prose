// Package sv registers the Swedish part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/sv"
//
//	tagger, err := tag.ForLanguage("sv")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Swedish-Talbanken@c434778 (2026-05-06)
// (license: CC BY-SA 4.0), and scores 93.68% on that treebank's development
// section and 94.86% on its test section. See MODELS.md for how.
package sv

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed sv.bin
var model []byte

func init() { tag.RegisterLanguage("sv", model) }
