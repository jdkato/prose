// Package sl registers the Slovene sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/sl"
//
//	seg, err := segment.ForLanguage("sl")
//
// The model is the Slovene Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package sl

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed slovene.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"npr", "itd", "oz", "dr", "prof", "g", "ga", "št", "str", "tel",
}

func init() { segment.RegisterLanguage("sl", model, abbreviations...) }
