// Package strcase converts strings between cases: title, sentence, camel,
// snake, and others.
//
// The title and sentence converters are language-aware — they consult a
// part-of-speech tagger to decide which words to leave lowercase — while the
// word-level converters in word.go are purely mechanical.
package strcase

import (
	"sort"
	"strings"

	"github.com/jdkato/prose/v3/internal/re"
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

// A vocabEntry is a vocab word and its compiled pattern, which is nil when the
// word isn't a valid pattern.
type vocabEntry struct {
	word    string
	pattern *re.Regexp
}

// compileVocab compiles each vocab word once, so a lookup costs no compiles.
func compileVocab(vocab []string) []vocabEntry {
	entries := make([]vocabEntry, len(vocab))
	for i, word := range vocab {
		entries[i].word = word
		entries[i].pattern, _ = re.Compile(word)
	}
	return entries
}

// lookupVocab returns the vocab word equal to s, ignoring case; s if a word's
// pattern matches it; or "" otherwise.
func lookupVocab(entries []vocabEntry, s string) string {
	for _, e := range entries {
		if strings.EqualFold(e.word, s) {
			return e.word
		} else if e.pattern != nil && e.pattern.MatchString(s) {
			return s
		}
	}
	return ""
}
