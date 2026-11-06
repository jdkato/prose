// Package sr registers the Serbian part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/sr"
//
//	tagger, err := tag.ForLanguage("sr")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Serbian-SET@17547a1 (2026-05-06)
// (license: CC BY-SA 4.0), and scores 96.69% on that treebank's development
// section and 97.34% on its test section. See MODELS.md for how.
package sr

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed sr.bin
var model []byte

func init() { tag.RegisterLanguage("sr", model) }
