// Package tr registers the Turkish part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/tr"
//
//	tagger, err := tag.ForLanguage("tr")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Turkish-BOUN@f828b02 (2026-05-06)
// (license: CC BY-SA 4.0), and scores 89.72% on that treebank's development
// section and 90.93% on its test section. See MODELS.md for how.
package tr

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed tr.bin
var model []byte

func init() { tag.RegisterLanguage("tr", model) }
