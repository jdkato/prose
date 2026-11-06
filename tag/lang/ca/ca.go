// Package ca registers the Catalan part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/ca"
//
//	tagger, err := tag.ForLanguage("ca")
//
// The model emits the Universal Dependencies tag set. It was trained on
// UD_Catalan-AnCora@7f3d23e (2026-05-06)
// (license: CC BY 4.0), and scores 98.03% on that treebank's development
// section and 97.82% on its test section. See MODELS.md for how.
package ca

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed ca.bin
var model []byte

func init() { tag.RegisterLanguage("ca", model) }
