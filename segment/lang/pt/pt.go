// Package pt registers the Portuguese sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/pt"
//
//	seg, err := segment.ForLanguage("pt")
//
// The model is the Portuguese Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package pt

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed portuguese.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"sr", "sra", "srta", "dr", "dra", "prof", "eng", "av", "tel", "pág", "ex",
	"etc", "cia", "fig", "cap", "vol",
}

func init() { segment.RegisterLanguage("pt", model, abbreviations...) }
