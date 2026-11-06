// Package lt registers the Lithuanian part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/lt"
//
//	tagger, err := tag.ForLanguage("lt")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Lithuanian-ALKSNIS@2dd19c5 (2026-05-06)
// (license: CC BY-SA 4.0), and scores 92.02% on that treebank's development
// section and 89.81% on its test section. See MODELS.md for how.
package lt

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed lt.bin
var model []byte

func init() { tag.RegisterLanguage("lt", model) }
