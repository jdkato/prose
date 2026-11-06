// Package nb registers the Norwegian Bokmål part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/nb"
//
//	tagger, err := tag.ForLanguage("nb")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Norwegian-Bokmaal@396d11f (2026-05-06)
// (license: CC BY-SA 4.0), and scores 96.34% on that treebank's development
// section and 95.96% on its test section. See MODELS.md for how.
package nb

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed nb.bin
var model []byte

func init() { tag.RegisterLanguage("nb", model) }
