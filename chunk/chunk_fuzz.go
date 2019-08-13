// +build gofuzz

package chunk

import (
	"github.com/jdkato/prose/tag"
	"github.com/jdkato/prose/tokenize"
)

func Fuzz(data []byte) int {
	words := tokenize.TextToWords(string(data))
	if len(words) == 0 {
		return 0
	}

	tagger := tag.NewPerceptronTagger()
	tagged := tagger.Tag(words)
	if len(tagged) == 0 {
		return 0
	}

	chunks := Chunk(tagged, TreebankNamedEntities)
	if len(chunks) == 0 {
		return 0
	}

	return 1
}
