// Package sv registers the Swedish sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/sv"
//
//	seg, err := segment.ForLanguage("sv")
//
// The model is the Swedish Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package sv

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed swedish.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"hr", "fru", "dr", "ca", "osv", "nr", "tel", "s", "st", "kl", "ev", "tfn",
}

func init() { segment.RegisterLanguage("sv", model, abbreviations...) }
