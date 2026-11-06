// Package tr registers the Turkish sentence-segmentation model with
// prose/v3/segment.
//
// Import it for its side effect, then ask for the segmenter by language:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/tr"
//
//	seg, err := segment.ForLanguage("tr")
//
// The model is the Turkish Punkt model distributed with
// gopkg.in/neurosnap/sentences (MIT), stored compressed and decoded on first
// use, plus the common abbreviations below that it did not learn from its
// corpus.
package tr

import (
	_ "embed"

	"github.com/jdkato/prose/v3/segment"
)

//go:embed turkish.json.gz
var model []byte

// abbreviations end in a period without ending a sentence.
var abbreviations = []string{
	"dr", "prof", "doç", "av", "no", "tel", "vb", "vs", "örn", "bkz", "sn",
	"yy", "hz", "mah", "cad", "sok", "apt",
}

func init() { segment.RegisterLanguage("tr", model, abbreviations...) }
