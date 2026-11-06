// Package perceptron holds what the averaged-perceptron tagger and its trainer
// must agree on: how a word is normalized, and which features describe it.
//
// The tagger scores features by string key, so a model is only as good as the
// match between the keys it was trained with and the keys built at tagging
// time. Both sides build them here, through one code path, which is what makes
// training a model in one package and loading it in another safe.
package perceptron

import (
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// Template numbers a feature set. A model records the template it was trained
// with, and the tagger builds the same one.
type Template int

const (
	// TemplateBytes is the original English template: a word's affixes are
	// its first byte and last three bytes. Adequate for ASCII, and what the
	// shipped English model was trained against, so it stays as it is.
	TemplateBytes Template = 0

	// TemplateRunes measures affixes in runes, so "straße" and "Москва" get
	// a real suffix rather than a torn byte sequence, and adds a shape
	// feature (capitalized, upper-case, digits) that is cheap and tells a
	// name from a noun in every script that has case.
	TemplateRunes Template = 1
)

// Sentinels pad the context so features can address two words to either
// side of the first and last without bounds checks.
const (
	Start  = "-START-"
	Start2 = "-START2-"
	End    = "-END-"
	End2   = "-END2-"
)

// Context returns words normalized and padded with the sentinels: the word at
// index i of words is at index i+2 of the result. dst is reused when large
// enough.
func Context(dst []string, words []string) []string {
	need := len(words) + 4
	if cap(dst) < need {
		dst = make([]string, need)
	}
	ctx := dst[:need]
	ctx[0] = Start
	ctx[1] = Start2
	for i, w := range words {
		ctx[i+2] = Normalize(w)
	}
	ctx[need-2] = End
	ctx[need-1] = End2
	return ctx
}

// Features builds the feature keys for one word into a reusable buffer.
//
// Keys are views into the buffer rather than strings of their own, so once the
// buffer has grown to fit a sentence, building the features for a word
// allocates nothing. A key is valid only until the next Build.
type Features struct {
	template Template
	buf      []byte
	ends     []int
}

// NewFeatures returns a builder for the given template.
func NewFeatures(t Template) *Features {
	return &Features{template: t, buf: make([]byte, 0, 512), ends: make([]int, 0, 24)}
}

// Template reports which feature set the builder produces.
func (f *Features) Template() Template { return f.template }

// Len is the number of keys built by the last Build.
func (f *Features) Len() int { return len(f.ends) }

// Key returns the k-th key as a string view into the buffer. It does not
// allocate; the caller must not retain it past the next Build.
func (f *Features) Key(k int) string {
	start := 0
	if k > 0 {
		start = f.ends[k-1]
	}
	b := f.buf[start:f.ends[k]]
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

// Strings returns copies of every key, for callers that keep them.
func (f *Features) Strings() []string {
	out := make([]string, f.Len())
	for k := range out {
		out[k] = strings.Clone(f.Key(k))
	}
	return out
}

// Build fills the builder with the features of the word at ctx[i].
//
// ctx is the output of Context, i already includes its offset of two, word is
// the raw (un-normalized) word, and p1 and p2 are the tags assigned to the
// previous two words.
func (f *Features) Build(ctx []string, i int, word, p1, p2 string) {
	f.buf = f.buf[:0]
	f.ends = f.ends[:0]

	var suf, pre, iMinus, iPlus string
	switch f.template {
	case TemplateRunes:
		suf = suffixRunes(word, 3)
		pre = prefixRunes(word, 1)
		iMinus = suffixRunes(ctx[i-1], 3)
		iPlus = suffixRunes(ctx[i+1], 3)
	default:
		suf = suffixBytes(word, 3)
		pre = prefixByte(word)
		iMinus = suffixBytes(ctx[i-1], 3)
		iPlus = suffixBytes(ctx[i+1], 3)
	}

	// Each feature is scored with weight 1. A feature that appears twice
	// contributes twice, which is correct because scoring is linear.
	//
	// The byte template spells its feature names out, because the English
	// model was trained with them spelled out. The rune template names each
	// feature with one byte: the names are a third of a model's size.
	if f.template == TemplateBytes {
		f.add("bias")
		f.add("i suffix ", suf)
		f.add("i pref1 ", pre)
		f.add("i-1 tag ", p1)
		f.add("i-2 tag ", p2)
		f.add("i tag+i-2 tag ", p1, " ", p2)
		f.add("i word ", ctx[i])
		f.add("i-1 tag+i word ", p1, " ", ctx[i])
		f.add("i-1 word ", ctx[i-1])
		f.add("i-1 suffix ", iMinus)
		f.add("i-2 word ", ctx[i-2])
		f.add("i+1 word ", ctx[i+1])
		f.add("i+1 suffix ", iPlus)
		f.add("i+2 word ", ctx[i+2])
		return
	}

	f.add("a")                       // bias
	f.add("b", suf)                  // i suffix
	f.add("c", pre)                  // i pref1
	f.add("d", p1)                   // i-1 tag
	f.add("e", p2)                   // i-2 tag
	f.add("f", p1, " ", p2)          // i-1 tag + i-2 tag
	f.add("g", ctx[i])               // i word
	f.add("h", p1, " ", ctx[i])      // i-1 tag + i word
	f.add("i", ctx[i-1])             // i-1 word
	f.add("j", iMinus)               // i-1 suffix
	f.add("k", ctx[i-2])             // i-2 word
	f.add("l", ctx[i+1])             // i+1 word
	f.add("m", iPlus)                // i+1 suffix
	f.add("n", ctx[i+2])             // i+2 word
	f.add("o", suffixRunes(word, 2)) // i suffix2
	f.add("p", suffixRunes(word, 4)) // i suffix4
	f.add("q", prefixRunes(word, 2)) // i pref2
	f.add("r", prefixRunes(word, 3)) // i pref3
	f.add("s", Shape(word))          // i shape
	f.add("t", p1, " ", suf)         // i-1 tag + i suffix
}

func (f *Features) add(parts ...string) {
	for _, p := range parts {
		f.buf = append(f.buf, p...)
	}
	f.ends = append(f.ends, len(f.buf))
}

// Normalize maps a word onto the form the model was trained against.
//
// Allocation-free on the common path. Note the year test: strconv.Atoi would
// be the obvious way to write it, but Atoi allocates a *NumError on failure —
// which is almost every word — so it scans digits directly instead.
func Normalize(word string) string {
	if word == "" {
		return word
	}
	first := word[0]
	if first != '-' && strings.IndexByte(word, '-') >= 0 {
		return "!HYPHEN"
	}
	if len(word) == 4 && isDecimalInt(word) {
		return "!YEAR"
	}
	if first >= '0' && first <= '9' {
		return "!DIGITS"
	}
	return toLower(word)
}

// Shape describes a word's casing and makeup in a few letters: "Xx" for a
// capitalized word, "XX" for one in capitals, "xx" for lower case, "d" for
// digits, "p" for punctuation and "o" for anything else. Mixed-case words
// such as "iPhone" are "xX".
func Shape(word string) string {
	var upper, lower, digit, punct, other int
	for _, r := range word {
		switch {
		case unicode.IsUpper(r):
			upper++
		case unicode.IsLower(r):
			lower++
		case unicode.IsDigit(r):
			digit++
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			punct++
		default:
			other++
		}
	}
	letters := upper + lower
	switch {
	case letters == 0 && digit > 0:
		return "d"
	case letters == 0 && punct > 0:
		return "p"
	case letters == 0:
		return "o"
	case lower == 0:
		return "XX"
	case upper == 0:
		return "xx"
	}
	r, _ := utf8.DecodeRuneInString(word)
	if unicode.IsUpper(r) && upper == 1 {
		return "Xx"
	}
	return "xX"
}

// isDecimalInt reports whether s is what strconv.Atoi would accept: an
// optional sign followed by at least one ASCII digit. Only ever called with
// len(s) == 4, so overflow is not a concern.
func isDecimalInt(s string) bool {
	i := 0
	if s[0] == '+' || s[0] == '-' {
		i = 1
	}
	if i == len(s) {
		return false
	}
	for ; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// toLower avoids allocating when a word is already lowercase ASCII, which is
// the common case once a sentence's first word is past. Anything with an
// uppercase letter or a non-ASCII byte falls through to strings.ToLower for
// correct Unicode folding.
func toLower(s string) string {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'A' && c <= 'Z') || c >= utf8.RuneSelf {
			return strings.ToLower(s)
		}
	}
	return s
}

// suffixBytes returns the last n bytes of w, or all of w if shorter.
func suffixBytes(w string, n int) string {
	if len(w) <= n {
		return w
	}
	return w[len(w)-n:]
}

// prefixByte returns the first byte of w as a string.
//
// Byte-oriented rather than rune-oriented, deliberately: the English model was
// trained with these keys, so switching to runes would silently stop matching
// learned features.
func prefixByte(w string) string {
	if w == "" {
		return ""
	}
	return w[:1]
}

// suffixRunes returns the last n runes of w, or all of w if shorter. It
// slices rather than copies, so it does not allocate.
func suffixRunes(w string, n int) string {
	i := len(w)
	for k := 0; k < n && i > 0; k++ {
		_, size := utf8.DecodeLastRuneInString(w[:i])
		i -= size
	}
	return w[i:]
}

// prefixRunes returns the first n runes of w, or all of w if shorter.
func prefixRunes(w string, n int) string {
	i := 0
	for k := 0; k < n && i < len(w); k++ {
		_, size := utf8.DecodeRuneInString(w[i:])
		i += size
	}
	return w[:i]
}
