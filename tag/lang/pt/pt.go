// Package pt registers the Portuguese part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/pt"
//
//	tagger, err := tag.ForLanguage("pt")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Portuguese-Bosque@8842885 (2026-05-06)
// (license: CC BY-SA 4.0), and scores 95.72% on that treebank's development
// section and 95.58% on its test section. See MODELS.md for how.
package pt

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed pt.bin
var model []byte

func init() { tag.RegisterLanguage("pt", model) }
