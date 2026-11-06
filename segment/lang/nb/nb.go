// Package nb registers the Norwegian sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/nb"
//
//	seg, err := segment.ForLanguage("nb")
//
// The model is the Norwegian Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package nb

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed norwegian.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"hr", "fru", "dr", "ca", "osv", "nr", "tlf", "s", "kl", "jf", "mv", "evt",
	"bl.a", "f.eks",
}

func init() { segment.RegisterLanguage("nb", model, abbreviations...) }
