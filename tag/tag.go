// Package tag assigns part-of-speech tags to tokenized text.
//
// The tagger is an averaged perceptron, after Honnibal's "A Good Part-of-Speech
// Tagger in about 200 Lines of Python". The English model uses the Penn
// Treebank tag set; the models for other languages, registered by the packages
// under prose/v3/tag/lang, use the Universal Dependencies tag set. See
// ForLanguage.
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
	"strconv"
	"strings"
	"sync"

	"github.com/jdkato/prose/v3/internal/emoticon"
	"github.com/jdkato/prose/v3/internal/flat"
	"github.com/jdkato/prose/v3/internal/perceptron"
	"github.com/jdkato/prose/v3/tag/aptagmodel"
	"github.com/jdkato/prose/v3/token"
)

// Tag sets a model can emit.
const (
	// Penn is the Penn Treebank tag set (NN, VBZ, JJ, ...), used by the
	// English model.
	Penn = "penn"

	// Universal is the Universal Dependencies tag set (NOUN, VERB, ADJ, ...),
	// used by every other language's model. It has 17 tags, the same in
	// every language, which is what lets one rule read across languages.
	Universal = "upos"
)

// Penn Treebank corpora mark traces and elisions with tokens like "*-1" or
// "*T*-2", and bracket certain literals as -LRB- / -RRB-. Both are passed
// through rather than scored.
var (
	noneRe = regexp.MustCompile(`^(?:0|\*[\w?]\*|\*\-\d{1,3}|\*[A-Z]+\*\-\d{1,3}|\*)$`)
	keepRe = regexp.MustCompile(`^\-[A-Z]{3}\-$`)
)

func isEmoticon(s string) bool { return emoticon.Is(s) }

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
	model    *flat.Model
	lexicon  Lexicon
	pool     sync.Pool
	template perceptron.Template
	tagset   string
	lang     string
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
	t := &Tagger{model: m, tagset: Penn, lang: "en"}

	// A model says what it is; one from before metadata existed is the
	// English Penn Treebank model, which is all there was.
	if ts := m.Meta(flat.MetaTagset); ts != "" {
		t.tagset = ts
	}
	if lang := m.Meta(flat.MetaLanguage); lang != "" {
		t.lang = lang
	}
	if n, err := strconv.Atoi(m.Meta(flat.MetaTemplate)); err == nil {
		t.template = perceptron.Template(n)
	}

	for _, opt := range opts {
		opt(t)
	}
	t.pool.New = func() any {
		return &scratch{
			scores: make([]float32, len(m.Classes())),
			feats:  perceptron.NewFeatures(t.template),
		}
	}
	return t
}

// scratch is the per-call working set, pooled so that tagging allocates
// nothing in steady state.
type scratch struct {
	scores  []float32
	feats   *perceptron.Features
	context []string
}

// Name reports which tagger this is.
func (t *Tagger) Name() string { return Default }

// Tagset reports the tag set the tagger emits: Penn or Universal.
func (t *Tagger) Tagset() string { return t.tagset }

// Metadata keys a model carries, for Meta.
const (
	// MetaSource describes the training data.
	MetaSource = flat.MetaSource
	// MetaLicense is the license of the training data.
	MetaLicense = flat.MetaLicense
)

// Meta returns what the model's metadata says under key, or "" when it says
// nothing: the English model predates metadata and says nothing at all.
func (t *Tagger) Meta(key string) string { return t.model.Meta(key) }

// Language reports the language the tagger's model was trained for, as a
// BCP 47 tag such as "en" or "de".
func (t *Tagger) Language() string { return t.lang }

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
	ctx[0] = perceptron.Start
	ctx[1] = perceptron.Start2
	for i := range tokens {
		ctx[i+2] = perceptron.Normalize(tokens[i].Text)
	}
	ctx[need-2] = perceptron.End
	ctx[need-1] = perceptron.End2

	p1, p2 := perceptron.Start, perceptron.Start2
	for i := range tokens {
		tag := t.tagOne(s, i, ctx, tokens[i].Text, p1, p2)
		tokens[i].Tag = tag
		p2, p1 = p1, tag
	}
}

// tagOne resolves a single token. The special cases below short-circuit the
// model for tokens it was never trained on.
func (t *Tagger) tagOne(s *scratch, i int, ctx []string, word, p1, p2 string) string {
	if tag, ok := t.fixed(word); ok {
		return tag
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

// fixed returns the tag of a token the model never decides.
//
// The Penn model has several: the corpus conventions for traces and brackets,
// and a few token shapes it was never trained on. A universal model decides
// everything but emoticons, which are symbols in every language.
func (t *Tagger) fixed(word string) (string, bool) {
	if t.tagset != Penn {
		if isEmoticon(word) {
			return "SYM", true
		}
		return "", false
	}
	switch {
	case word == "-":
		return "-", true
	case isEmoticon(word):
		return "SYM", true
	case strings.HasPrefix(word, "@"):
		return "NN", true
	case noneRe.MatchString(word):
		return "-NONE-", true
	case keepRe.MatchString(word):
		return word, true
	}
	return "", false
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
// feature key is assembled in a reusable buffer and looked up by value.
func (t *Tagger) predict(s *scratch, i int, ctx []string, word, p1, p2 string) string {
	for j := range s.scores {
		s.scores[j] = 0
	}

	i = min(len(ctx)-3, i+2)

	s.feats.Build(ctx, i, word, p1, p2)
	for k, n := 0, s.feats.Len(); k < n; k++ {
		t.model.AccumulateInto(s.feats.Key(k), 1, s.scores)
	}

	return t.model.Best(s.scores)
}
