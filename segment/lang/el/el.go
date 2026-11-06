// Package el registers the Greek sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/el"
//
//	seg, err := segment.ForLanguage("el")
//
// The model is the Greek Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package el

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed greek.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"κ", "κα", "δρ", "σελ", "αρ", "τηλ", "βλ", "κεφ", "π.χ",
}

func init() { segment.RegisterLanguage("el", model, abbreviations...) }
