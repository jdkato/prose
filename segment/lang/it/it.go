// Package it registers the Italian sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/it"
//
//	seg, err := segment.ForLanguage("it")
//
// The model is the Italian Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package it

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed italian.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"sig", "sigg", "dott", "prof", "avv", "ing", "geom", "rag", "n", "pag",
	"tel", "es", "ecc", "cfr", "vol", "art", "fig", "cap",
}

func init() { segment.RegisterLanguage("it", model, abbreviations...) }
