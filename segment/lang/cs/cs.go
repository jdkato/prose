// Package cs registers the Czech sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/cs"
//
//	seg, err := segment.ForLanguage("cs")
//
// The model is the Czech Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package cs

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed czech.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"např", "tzv", "atd", "apod", "tj", "č", "str", "prof", "dr", "ing", "mgr",
	"ul", "nám", "tel", "cca", "resp", "popř", "mj",
}

func init() { segment.RegisterLanguage("cs", model, abbreviations...) }
