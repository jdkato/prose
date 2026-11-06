// Package nl registers the Dutch sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/nl"
//
//	seg, err := segment.ForLanguage("nl")
//
// The model is the Dutch Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package nl

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed dutch.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"dhr", "mevr", "mw", "dr", "drs", "prof", "ir", "mr", "nr", "bv", "bijv",
	"ca", "enz", "etc", "evt", "excl", "incl", "pag", "tel", "zgn",
}

func init() { segment.RegisterLanguage("nl", model, abbreviations...) }
