/*
Package tokenize implements functions to split strings into slices of substrings.
*/
package tokenize

import (
	"github.com/jdkato/prose/internal/util"
	"gopkg.in/neurosnap/sentences.v1/english"
)

// TextToWords converts the string text into a slice of words.
//
// It does so by tokenizing text into sentences (using a port of NLTK's punkt
// tokenizer; see https://github.com/neurosnap/sentences) and then tokenizing
// the sentences into words via TreebankWordTokenizer.
func TextToWords(text string) []string {
	sentTokenizer, err := english.NewSentenceTokenizer(nil)
	util.CheckError(err)
	wordTokenizer := NewTreebankWordTokenizer()

	words := []string{}
	for _, s := range sentTokenizer.Tokenize(text) {
		words = append(words, wordTokenizer.Tokenize(s.Text)...)
	}

	return words
}
