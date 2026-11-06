// Package fr registers the French sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/fr"
//
//	seg, err := segment.ForLanguage("fr")
//
// The model is the French Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package fr

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed french.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"m", "mme", "mlle", "mm", "dr", "pr", "st", "ste", "av", "bd", "tél", "etc",
	"env", "cf", "ex", "réf", "vol", "chap", "fig", "p", "pp", "éd",
}

func init() { segment.RegisterLanguage("fr", model, abbreviations...) }
