// Package da registers the Danish sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/da"
//
//	seg, err := segment.ForLanguage("da")
//
// The model is the Danish Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package da

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed danish.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"hr", "fru", "frk", "dr", "ca", "evt", "osv", "fx", "nr", "tlf", "s", "jf",
	"mv", "pga", "vedr",
}

func init() { segment.RegisterLanguage("da", model, abbreviations...) }
