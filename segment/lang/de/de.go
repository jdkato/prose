// Package de registers the German sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/de"
//
//	seg, err := segment.ForLanguage("de")
//
// The model is the German Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package de

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed german.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"bzw", "ca", "nr", "usw", "bzgl", "ggf", "evtl", "vgl", "dr", "prof", "str",
	"hr", "fr", "tel", "abs", "art", "inkl", "insb", "mind", "max", "min",
	"sog", "zzgl", "abb", "bd", "bsp", "jh", "hrsg",
}

func init() { segment.RegisterLanguage("de", model, abbreviations...) }
