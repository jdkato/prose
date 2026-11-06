// Package es registers the Spanish sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/es"
//
//	seg, err := segment.ForLanguage("es")
//
// The model is the Spanish Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package es

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed spanish.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"sr", "sra", "srta", "dr", "dra", "ud", "uds", "vd", "vds", "núm", "pág",
	"págs", "tel", "av", "avda", "aprox", "dpto", "fig", "cap", "vol", "ej",
}

func init() { segment.RegisterLanguage("es", model, abbreviations...) }
