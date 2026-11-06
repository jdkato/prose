// Package it registers the Italian part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/it"
//
//	tagger, err := tag.ForLanguage("it")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Italian-MarkIT@e210079 (2026-05-06); UD_Italian-PUD@71d932a (2026-05-06); UD_Italian-Valico@8d60b2c (2026-05-06)
// (license: CC BY 4.0, CC BY-SA 3.0, CC BY-SA 4.0), and scores 92.57% on that treebank's development
// section and 92.42% on its test section. See MODELS.md for how.
package it

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed it.bin
var model []byte

func init() { tag.RegisterLanguage("it", model) }
