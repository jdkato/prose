// Package ru registers the Russian part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/ru"
//
//	tagger, err := tag.ForLanguage("ru")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Russian-GSD@0f34b73 (2026-05-06); UD_Russian-Taiga@fcfd7db (2026-05-06)
// (license: CC BY-SA 4.0), and scores 96.38% on that treebank's development
// section and 95.65% on its test section. See MODELS.md for how.
package ru

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed ru.bin
var model []byte

func init() { tag.RegisterLanguage("ru", model) }
