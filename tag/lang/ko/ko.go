// Package ko registers the Korean part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/ko"
//
//	tagger, err := tag.ForLanguage("ko")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Korean-GSD@7628157 (2026-05-06)
// (license: CC BY-SA 4.0), and scores 90.04% on that treebank's development
// section and 90.73% on its test section. See MODELS.md for how.
package ko

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed ko.bin
var model []byte

func init() { tag.RegisterLanguage("ko", model) }
