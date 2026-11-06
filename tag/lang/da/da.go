// Package da registers the Danish part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/da"
//
//	tagger, err := tag.ForLanguage("da")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Danish-DDT@5880293 (2026-05-06)
// (license: CC BY-SA 4.0), and scores 94.22% on that treebank's development
// section and 94.23% on its test section. See MODELS.md for how.
package da

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed da.bin
var model []byte

func init() { tag.RegisterLanguage("da", model) }
