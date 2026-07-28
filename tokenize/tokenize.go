// Package tokenize splits text into tokens that know where they came from.
//
// # Positions
//
// Every token carries a byte offset into the input, and the invariant
//
//	src[tok.Start:tok.End()] == tok.Text
//
// holds for all of them. Two design choices make that true by construction
// rather than by care:
//
//  1. The input is never rewritten before tokenizing. Normalizing curly quotes
//     to straight ones, or "&rsquo;" to "'", changes byte lengths, so every
//     offset measured against the rewritten copy is wrong for the original.
//     Characters that would otherwise be normalized are handled directly by
//     the prefix and suffix sets below.
//
//  2. Splitting carries offsets rather than recomputing them. Each split step
//     slices both the text and its base offset, so a token's position is
//     always a running sum, never the result of searching for the token's text
//     in the source, which finds the wrong instance for any repeated word.
//
// # Performance
//
// Tokenizing allocates one slice for the result and nothing per token: tokens
// are substrings of the input, not copies.
package tokenize

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jdkato/prose/v3/internal/emoticon"
	"github.com/jdkato/prose/v3/token"
)

// Token is prose's shared token type, re-exported for convenience.
type Token = token.Token

// Tokenizer splits text into tokens.
type Tokenizer struct {
	prefixes   []string
	suffixes   []string
	splitCases []string
	specialRE  *regexp.Regexp
	isSpecial  func(string) bool
}

// Option configures a Tokenizer.
type Option func(*Tokenizer)

// UsingPrefixes replaces the set of characters split off the front of a token.
func UsingPrefixes(p []string) Option {
	return func(t *Tokenizer) { t.prefixes = p }
}

// UsingSuffixes replaces the set of characters split off the end of a token.
func UsingSuffixes(s []string) Option {
	return func(t *Tokenizer) { t.suffixes = s }
}

// UsingSplitCases replaces the set of infixes that trigger a split, such as
// contraction endings.
func UsingSplitCases(c []string) Option {
	return func(t *Tokenizer) { t.splitCases = c }
}

// UsingSpecialTest supplies an extra predicate for spans that must not be
// split at all.
func UsingSpecialTest(fn func(string) bool) Option {
	return func(t *Tokenizer) { t.isSpecial = fn }
}

// New returns a Tokenizer with prose's default English behaviour.
func New(opts ...Option) *Tokenizer {
	t := &Tokenizer{
		prefixes:   defaultPrefixes,
		suffixes:   defaultSuffixes,
		splitCases: defaultContractions,
		specialRE:  internalRE,
		isSpecial:  func(string) bool { return false },
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Tokenize splits text into tokens.
//
// The returned tokens hold substrings of text, so they stay valid only as
// long as text does — which is free, since Go strings are immutable.
func (t *Tokenizer) Tokenize(text string) []Token {
	// One allocation for the result. Punctuation splits push the token count
	// above a naive words-per-byte estimate, so this errs high: regrowing a
	// []Token costs more than the slack does.
	tokens := make([]Token, 0, len(text)/4+2)

	var (
		start, index int
		white        bool
		length       = len(text)
	)

	for index <= length {
		uc, size := utf8.DecodeRuneInString(text[index:])
		if size == 0 {
			break
		} else if index == 0 {
			white = unicode.IsSpace(uc)
		}
		if unicode.IsSpace(uc) != white {
			if start < index {
				tokens = t.split(text[start:index], start, tokens)
			}
			if uc == ' ' {
				start = index + 1
			} else {
				start = index
			}
			white = !white
		}
		index += size
	}

	if start < index {
		tokens = t.split(text[start:index], start, tokens)
	}

	return tokens
}

// split breaks one whitespace-delimited span into tokens.
//
// base is the span's byte offset in the original text; every offset produced
// here is base plus a distance within the span, so positions stay anchored to
// the source no matter how many times the span is subdivided.
//
// Note that this cannot be memoized by span text: two occurrences of the same
// span sit at different offsets. It would buy little anyway — the work is a
// handful of prefix and suffix tests.
func (t *Tokenizer) split(span string, base int, tokens []Token) []Token {
	// Suffixes peel off the end, so they are collected in reverse and
	// appended once the core of the span is resolved.
	//
	// The backing array is a local, so a span with the usual handful of
	// trailing punctuation (`."` , `)."`) costs no allocation at all; only a
	// pathological run of more than eight escapes to the heap.
	var sufArr [8]Token
	suffixes := sufArr[:0]

	for span != "" {
		if t.special(span) {
			tokens = append(tokens, Token{Text: span, Start: base})
			span = ""
			break
		}

		if p := matchPrefix(span, t.prefixes); p > 0 {
			tokens = append(tokens, Token{Text: span[:p], Start: base})
			span, base = span[p:], base+p
			continue
		}

		if idx := indexAnyCase(span, t.splitCases); idx > 0 {
			// "they'll" -> "they" + "'ll"; "don't" -> "do" + "n't".
			tokens = append(tokens, Token{Text: span[:idx], Start: base})
			span, base = span[idx:], base+idx
			continue
		}

		if s := matchSuffix(span, t.suffixes); s > 0 {
			cut := len(span) - s
			suffixes = append(suffixes, Token{Text: span[cut:], Start: base + cut})
			span = span[:cut]
			continue
		}

		tokens = append(tokens, Token{Text: span, Start: base})
		span = ""
	}

	// Emit collected suffixes in source order.
	for i := len(suffixes) - 1; i >= 0; i-- {
		tokens = append(tokens, suffixes[i])
	}
	return tokens
}

// special reports whether a span must be kept whole.
func (t *Tokenizer) special(span string) bool {
	return emoticon.Is(span) || t.specialRE.MatchString(span) || t.isSpecial(span)
}

// matchPrefix returns the byte length of the leading prefix, or 0.
func matchPrefix(s string, prefixes []string) int {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) && len(s) > len(p) {
			return len(p)
		}
	}
	return 0
}

// matchSuffix returns the byte length of the trailing suffix, or 0.
func matchSuffix(s string, suffixes []string) int {
	for _, suf := range suffixes {
		if strings.HasSuffix(s, suf) && len(s) > len(suf) {
			return len(suf)
		}
	}
	return 0
}

// indexAnyCase returns the byte offset of the earliest split case in s,
// compared case-insensitively, or -1.
//
// Only ASCII case is folded, which is all the default split cases need, and it
// avoids allocating a lowercased copy of every span.
func indexAnyCase(s string, cases []string) int {
	best := -1
	for _, c := range cases {
		if i := indexFoldASCII(s, c); i > 0 && (best == -1 || i < best) {
			best = i
		}
	}
	return best
}

func indexFoldASCII(s, substr string) int {
	if len(substr) == 0 || len(substr) > len(s) {
		return -1
	}
	last := len(s) - len(substr)
	for i := 0; i <= last; i++ {
		if equalFoldASCII(s[i:i+len(substr)], substr) {
			return i
		}
	}
	return -1
}

func equalFoldASCII(a, b string) bool {
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

var internalRE = regexp.MustCompile(`^(?:[A-Za-z]\.){2,}$|^[A-Z][a-z]{1,2}\.$`)

// Both apostrophe forms are listed. Since the input is never rewritten, the
// curly forms have to be recognized directly or "It’s" stays one token.
var defaultContractions = []string{
	"'ll", "'s", "'re", "'m", "n't",
	"’ll", "’s", "’re", "’m", "n’t",
}

// The curly variants are here rather than in a pre-pass replacer: rewriting
// them to ASCII before tokenizing would shift every subsequent offset.
// Handling them as ordinary affixes keeps the input intact.
var defaultSuffixes = []string{
	",", ")", `"`, "]", "!", ";", ".", "?", ":", "'",
	"”", // right double quotation mark
	"’", // right single quotation mark
}

var defaultPrefixes = []string{
	"$", "(", `"`, "[",
	"“", // left double quotation mark
	"‘", // left single quotation mark
}
