package tokenize

import (
	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/english"
)

var sTokenizer, _ = english.NewSentenceTokenizer(nil)

// SentenceTokenizer splits text into sentences.
func SentenceTokenizer(text string) []*sentences.Sentence {
	return sTokenizer.Tokenize(text)
}
