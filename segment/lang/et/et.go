// Package et registers the Estonian sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/et"
//
//	seg, err := segment.ForLanguage("et")
//
// The model is the Estonian Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package et

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed estonian.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"nt", "jne", "vt", "lk", "nr", "dr", "prof", "hr", "pr", "tel", "vrd",
}

func init() { segment.RegisterLanguage("et", model, abbreviations...) }
