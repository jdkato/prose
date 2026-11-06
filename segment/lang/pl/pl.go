// Package pl registers the Polish sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/pl"
//
//	seg, err := segment.ForLanguage("pl")
//
// The model is the Polish Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package pl

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed polish.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"np", "tzn", "itp", "itd", "tj", "prof", "dr", "ul", "al", "nr", "godz",
	"ok", "św", "ks", "inż", "mgr", "pt", "tel",
}

func init() { segment.RegisterLanguage("pl", model, abbreviations...) }
