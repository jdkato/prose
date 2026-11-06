// Package fi registers the Finnish part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/fi"
//
//	tagger, err := tag.ForLanguage("fi")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Finnish-TDT@bfaae13 (2026-05-06)
// (license: CC BY-SA 4.0), and scores 92.61% on that treebank's development
// section and 92.71% on its test section. See MODELS.md for how.
package fi

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed fi.bin
var model []byte

func init() { tag.RegisterLanguage("fi", model) }
