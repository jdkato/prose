// Package fi registers the Finnish sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/fi"
//
//	seg, err := segment.ForLanguage("fi")
//
// The model is the Finnish Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package fi

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed finnish.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"esim", "ns", "yms", "jne", "mm", "n", "klo", "s", "ym", "ks", "vrt",
	"prof", "tri", "os", "puh",
}

func init() { segment.RegisterLanguage("fi", model, abbreviations...) }
