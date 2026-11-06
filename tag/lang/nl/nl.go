// Package nl registers the Dutch part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/nl"
//
//	tagger, err := tag.ForLanguage("nl")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Dutch-Alpino@7f0413c (2026-05-06)
// (license: CC BY-SA 4.0), and scores 95.29% on that treebank's development
// section and 93.86% on its test section. See MODELS.md for how.
package nl

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed nl.bin
var model []byte

func init() { tag.RegisterLanguage("nl", model) }
