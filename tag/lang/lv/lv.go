// Package lv registers the Latvian part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/lv"
//
//	tagger, err := tag.ForLanguage("lv")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Latvian-LVTB@f42ae0d (2026-05-06)
// (license: CC BY-SA 4.0), and scores 96.06% on that treebank's development
// section and 95.77% on its test section. See MODELS.md for how.
package lv

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed lv.bin
var model []byte

func init() { tag.RegisterLanguage("lv", model) }
