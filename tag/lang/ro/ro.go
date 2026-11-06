// Package ro registers the Romanian part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/ro"
//
//	tagger, err := tag.ForLanguage("ro")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Romanian-RRT@ee6055b (2026-05-06)
// (license: CC BY-SA 4.0), and scores 96.63% on that treebank's development
// section and 96.32% on its test section. See MODELS.md for how.
package ro

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed ro.bin
var model []byte

func init() { tag.RegisterLanguage("ro", model) }
