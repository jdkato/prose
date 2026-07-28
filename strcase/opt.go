// Package strcase converts strings between cases: title, sentence, camel,
// snake, and others.
//
// The title and sentence converters are language-aware — they consult a
// part-of-speech tagger to decide which words to leave lowercase — while the
// word-level converters in word.go are purely mechanical.
package strcase

import (
	"sort"
)

// A CaseConverter converts a string to a specific case.
type CaseConverter interface {
	Convert(string) string
}

// CaseOptFunc is a function that modifies a CaseConverter.
type CaseOptFunc func(opts *CaseOpts)

// An IndicatorFunc is a SentenceConverter callback that decides whether or not
// the string word should be capitalized.
type IndicatorFunc func(word string, idx int) bool

// CaseOpts is a struct that holds the options for a CaseConverter.
type CaseOpts struct {
	// A prefix to match before converting the string.
	//
	// Example:
	//
	// a. This is a sentence.
	prefix string

	// A list of words to specifying how to capitalize them.
	vocab []string

	// A function that determines whether or not a word should be capitalized.
	indicator IndicatorFunc
}

// UsingVocab sets the vocab for the CaseConverter.
func UsingVocab(vocab []string) CaseOptFunc {
	return func(opts *CaseOpts) {
		// NOTE: This is required to ensure that we have greedy alternation.
		sort.Slice(vocab, func(p, q int) bool {
			return len(vocab[p]) > len(vocab[q])
		})
		opts.vocab = vocab
	}
}

// UsingIndicator sets the indicator for the CaseConverter.
func UsingIndicator(indicator IndicatorFunc) CaseOptFunc {
	return func(opts *CaseOpts) {
		opts.indicator = indicator
	}
}

// UsingPrefix sets the prefix for the CaseConverter.
func UsingPrefix(prefix string) CaseOptFunc {
	return func(opts *CaseOpts) {
		opts.prefix = prefix
	}
}
