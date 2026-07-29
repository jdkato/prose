// Package tag assigns part-of-speech tags to tokenized text.
//
// The tagger is an averaged perceptron, after Honnibal's "A Good Part-of-Speech
// Tagger in about 200 Lines of Python". It uses the Penn Treebank tag set.
//
// # Performance
//
// The English model is loaded once per process, behind sync.OnceValues, and
// costs about 1.2 ms and 2 MB. Tagging runs at roughly 0.9 µs per token and,
// once the scratch pool is warm, allocates only for words that need Unicode
// case folding.
//
// That comes from three choices: weights are stored in a flat, sparse layout
// rather than nested maps; scores accumulate into pooled scratch space rather
// than a fresh map per token; and feature keys are assembled in a reusable
// buffer rather than joined into new strings.
//
// # Concurrency
//
// A Tagger is safe for concurrent use. The shared model is immutable after
// load, and per-call scratch space comes from a pool.
package tag

import (
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"
	"unsafe"

	"github.com/jdkato/prose/v3/internal/emoticon"
	"github.com/jdkato/prose/v3/internal/flat"
	"github.com/jdkato/prose/v3/tag/aptagmodel"
	"github.com/jdkato/prose/v3/token"
)

// Penn Treebank corpora mark traces and elisions with tokens like "*-1" or
// "*T*-2", and bracket certain literals as -LRB- / -RRB-. Both are passed
// through rather than scored.
var (
	noneRe = regexp.MustCompile(`^(?:0|\*[\w?]\*|\*\-\d{1,3}|\*[A-Z]+\*\-\d{1,3}|\*)$`)
	keepRe = regexp.MustCompile(`^\-[A-Z]{3}\-$`)
)

func isEmoticon(s string) bool { return emoticon.Is(s) }

// bytesToString views b as a string without copying. The result is used only
// for a map/blob lookup that does not retain it.
func bytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

// Token is prose's shared token type, re-exported so callers that only tag
// text need not import the token package separately.
//
// Tagging never touches Start, so a token's position survives the tagger
// untouched.
type Token = token.Token

// Tagger assigns part-of-speech tags with an averaged perceptron.
//
// It reads the words around the one it is tagging, which makes it accurate on
// ordinary prose and unsteady on text that is not: a grammar rule fires on a
// sentence that is already wrong, and the context there is not evidence for
// anything. See Lexical for a tagger that leans on a dictionary instead.
//
// It satisfies Interface, as every tagger here does.
type Tagger struct {
	model   *flat.Model
	lexicon Lexicon
	pool    sync.Pool
}

// A Lexicon fixes the tag for words that only ever carry one.
//
// The model predicts from context, which is what makes it general -- and what
// makes it wrong on text that is not idiomatic. "aware" is an adjective in
// every English sentence, but in "all ready aware of the risk" the surrounding
// words are not evidence for anything, and the model settles on NN. A caller
// that knows a word is unambiguous can say so.
//
// Lookups try the word as written, then its lower-case form, so one entry
// covers "aware", "Aware" and "AWARE".
//
// Tags are not validated against the model's own set: a caller may want a tag
// the model never emits.
type Lexicon map[string]string

// An Option configures a Tagger.
type Option func(*Tagger)

// WithLexicon fixes the tag for the given words.
//
// The lexicon takes precedence over the model, including over the model's own
// list of unambiguous words. It is not copied, and must not be modified once
// the Tagger is in use -- taggers are safe for concurrent use, and this is the
// one thing that would break that.
//
// This is the small case: a handful of words whose tag you want to fix. When a
// dictionary is what drives the tagging, use Lexical, which also carries the
// words that have more than one reading.
func WithLexicon(l Lexicon) Option {
	return func(t *Tagger) {
		t.lexicon = l
	}
}

// defaultModel loads the built-in English model on first use.
//
// Decoding is the tagger's only meaningful cost, so it happens once per
// process rather than once per caller.
var defaultModel = sync.OnceValues(func() (*flat.Model, error) {
	return flat.Unmarshal(aptagmodel.English)
})

// New returns a Tagger backed by the built-in English model.
//
// The model is loaded on the first call and shared thereafter, so calling New
// repeatedly is cheap.
func New(opts ...Option) (*Tagger, error) {
	m, err := defaultModel()
	if err != nil {
		return nil, err
	}
	return newWithModel(m, opts...), nil
}

// FromBytes returns a Tagger backed by a model in prose's flat format.
//
// data is retained and must not be modified afterwards; feature strings are
// sliced from it rather than copied.
func FromBytes(data []byte, opts ...Option) (*Tagger, error) {
	m, err := flat.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return newWithModel(m, opts...), nil
}

func newWithModel(m *flat.Model, opts ...Option) *Tagger {
	t := &Tagger{model: m}
	for _, opt := range opts {
		opt(t)
	}
	t.pool.New = func() any {
		return &scratch{
			scores: make([]float32, len(m.Classes())),
			key:    make([]byte, 0, 64),
		}
	}
	return t
}

// scratch is the per-call working set, pooled so that tagging allocates
// nothing in steady state.
type scratch struct {
	scores  []float32
	key     []byte
	context []string
}

// Name reports which tagger this is.
func (t *Tagger) Name() string { return Default }

// Tag assigns a tag to each word and returns the tagged tokens.
//
// The input is treated as one sentence; the tagger conditions on the previous
// two tags, so sentence boundaries matter.
func (t *Tagger) Tag(words []string) []Token {
	out := make([]Token, len(words))
	for i, w := range words {
		out[i].Text = w
	}
	t.TagTokens(out)
	return out
}

// TagTokens fills in the Tag field of each token, using its Text.
//
// This is the low-allocation entry point: it writes into the caller's slice
// rather than returning a new one. Once the pool is warm it allocates only
// for words needing Unicode case folding.
func (t *Tagger) TagTokens(tokens []Token) {
	if len(tokens) == 0 {
		return
	}

	s := t.pool.Get().(*scratch)
	defer t.pool.Put(s)

	// context is the normalized words padded with two sentinels at each end,
	// so featurize can index i-2 .. i+2 without bounds checks.
	need := len(tokens) + 4
	if cap(s.context) < need {
		s.context = make([]string, need)
	}
	ctx := s.context[:need]
	ctx[0] = "-START-"
	ctx[1] = "-START2-"
	for i := range tokens {
		ctx[i+2] = normalize(tokens[i].Text)
	}
	ctx[need-2] = "-END-"
	ctx[need-1] = "-END2-"

	p1, p2 := "-START-", "-START2-"
	for i := range tokens {
		tag := t.tagOne(s, i, ctx, tokens[i].Text, p1, p2)
		tokens[i].Tag = tag
		p2, p1 = p1, tag
	}
}

// tagOne resolves a single token. The special cases below short-circuit the
// model for tokens it was never trained on.
func (t *Tagger) tagOne(s *scratch, i int, ctx []string, word, p1, p2 string) string {
	switch {
	case word == "-":
		return "-"
	case isEmoticon(word):
		return "SYM"
	case strings.HasPrefix(word, "@"):
		return "NN"
	case noneRe.MatchString(word):
		return "-NONE-"
	case keepRe.MatchString(word):
		return word
	}
	// Ahead of the model's own lexicon: a caller supplying one is making a
	// more specific claim than the training data.
	if tag, ok := t.lookup(word); ok {
		return tag
	}

	if tag, ok := t.model.TagFor(word); ok {
		return tag
	}
	return t.predict(s, i, ctx, word, p1, p2)
}

// lookup consults the caller's lexicon, trying the word as written before its
// lower-case form so that one entry covers every casing.
func (t *Tagger) lookup(word string) (string, bool) {
	if t.lexicon == nil {
		return "", false
	}

	if tag, ok := t.lexicon[word]; ok {
		return tag, true
	}

	lower := strings.ToLower(word)
	if lower == word {
		return "", false
	}

	tag, ok := t.lexicon[lower]

	return tag, ok
}

// predict scores the model's features for one token and returns the best tag.
//
// This does not allocate: scores live in pooled scratch space, and each
// feature key is assembled in a reusable buffer.
func (t *Tagger) predict(s *scratch, i int, ctx []string, word, p1, p2 string) string {
	for j := range s.scores {
		s.scores[j] = 0
	}

	i = min(len(ctx)-3, i+2)

	suf := suffix(word, 3)
	pre := prefix1(word)
	iMinus := suffix(ctx[i-1], 3)
	iPlus := suffix(ctx[i+1], 3)

	// Each feature is scored with weight 1. A feature that appears twice
	// contributes twice, which is correct because scoring is linear.
	t.feat(s, "bias")
	t.feat(s, "i suffix ", suf)
	t.feat(s, "i pref1 ", pre)
	t.feat(s, "i-1 tag ", p1)
	t.feat(s, "i-2 tag ", p2)
	t.feat(s, "i tag+i-2 tag ", p1, " ", p2)
	t.feat(s, "i word ", ctx[i])
	t.feat(s, "i-1 tag+i word ", p1, " ", ctx[i])
	t.feat(s, "i-1 word ", ctx[i-1])
	t.feat(s, "i-1 suffix ", iMinus)
	t.feat(s, "i-2 word ", ctx[i-2])
	t.feat(s, "i+1 word ", ctx[i+1])
	t.feat(s, "i+1 suffix ", iPlus)
	t.feat(s, "i+2 word ", ctx[i+2])

	return t.model.Best(s.scores)
}

// feat assembles a feature key from parts and accumulates its weights.
//
// The key is built in scratch space; the model looks it up by value, so no
// string is allocated.
func (t *Tagger) feat(s *scratch, parts ...string) {
	s.key = s.key[:0]
	for _, p := range parts {
		s.key = append(s.key, p...)
	}
	t.model.AccumulateInto(bytesToString(s.key), 1, s.scores)
}

// suffix returns the last n bytes of w, or all of w if shorter.
func suffix(w string, n int) string {
	if len(w) <= n {
		return w
	}
	return w[len(w)-n:]
}

// prefix1 returns the first byte of w as a string.
//
// This is byte-oriented rather than rune-oriented, deliberately: the model was
// trained with these keys, so switching to runes would silently stop matching
// learned features.
func prefix1(w string) string {
	if w == "" {
		return ""
	}
	return w[:1]
}

// normalize maps a word onto the form the model was trained against.
//
// Allocation-free on the common path. Note the year test: strconv.Atoi would
// be the obvious way to write it, but Atoi allocates a *NumError on failure —
// which is almost every word — so it scans digits directly instead.
func normalize(word string) string {
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
