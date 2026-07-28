// Package prose is a natural language processing library for Go.
//
// It provides tokenization, sentence segmentation, part-of-speech tagging, and
// named-entity recognition. The pieces are usable on their own — see
// prose/v3/tokenize, prose/v3/segment, prose/v3/tag, and prose/v3/ner — and
// Document ties them together for the common case.
//
// # Positions
//
// Everything prose returns knows where it came from. For any token, sentence,
// or entity x produced from text:
//
//	text[x.Start:x.End()] == x.Text
//
// # Models are shared
//
// Every model is loaded once per process and shared, so constructing a
// Document is cheap and constructing many is cheaper still.
//
// # Choosing what you pay for
//
// Importing this package links every model prose ships (~8 MB). A program that
// only needs part-of-speech tags should import prose/v3/tag directly, which
// links only the tagger and leaves the entity model out of the binary
// entirely.
package prose

import (
	"context"
	"sync"

	"github.com/jdkato/prose/v3/ner"
	"github.com/jdkato/prose/v3/segment"
	"github.com/jdkato/prose/v3/tag"
	"github.com/jdkato/prose/v3/token"
	"github.com/jdkato/prose/v3/tokenize"
)

// Token is a piece of text with its position, part-of-speech tag, and entity
// label.
type Token = token.Token

// Sentence is a segmented span of text with its position.
type Sentence = segment.Sentence

// Entity is a named entity with its position.
type Entity = ner.Entity

// Document is text plus whatever prose was asked to compute about it.
type Document struct {
	// Text is the document's source text.
	Text string

	tokens    []Token
	sentences []Sentence
	entities  []Entity
}

// Tokens returns the document's tokens, or nil if tokenization was disabled.
func (d *Document) Tokens() []Token { return d.tokens }

// Sentences returns the document's sentences, or nil if segmentation was
// disabled.
func (d *Document) Sentences() []Sentence { return d.sentences }

// Entities returns the document's named entities, or nil if extraction was
// disabled.
func (d *Document) Entities() []Entity { return d.entities }

// config holds the pipeline settings for one NewDocument call.
type config struct {
	tokenize  bool
	segment   bool
	tag       bool
	extract   bool
	tokenizer *tokenize.Tokenizer
}

// Option configures a Document.
type Option func(*config)

// WithTokenization enables or disables tokenization. Enabled by default.
//
// Tagging and entity extraction both need tokens, so disabling this disables
// them too.
func WithTokenization(include bool) Option {
	return func(c *config) { c.tokenize = include }
}

// WithSegmentation enables or disables sentence segmentation. Enabled by
// default.
func WithSegmentation(include bool) Option {
	return func(c *config) { c.segment = include }
}

// WithTagging enables or disables part-of-speech tagging. Enabled by default.
func WithTagging(include bool) Option {
	return func(c *config) { c.tag = include }
}

// WithExtraction enables or disables named-entity extraction. Enabled by
// default.
//
// Extraction needs part-of-speech tags, so enabling it enables tagging.
func WithExtraction(include bool) Option {
	return func(c *config) { c.extract = include }
}

// defaultTokenizer returns the shared default tokenizer.
//
// A Tokenizer holds only its affix sets and is safe to share, so one is built
// per process rather than per document.
var defaultTokenizer = sync.OnceValue(func() *tokenize.Tokenizer {
	return tokenize.New()
})

// UsingTokenizer supplies a custom tokenizer.
func UsingTokenizer(t *tokenize.Tokenizer) Option {
	return func(c *config) { c.tokenizer = t }
}

// NewDocument creates a Document from text.
//
// By default it tokenizes, segments, tags, and extracts entities. Use the
// With* options to skip stages you do not need — skipping entity extraction
// in particular saves most of the work.
func NewDocument(text string, opts ...Option) (*Document, error) {
	return NewDocumentContext(context.Background(), text, opts...)
}

// NewDocumentContext is NewDocument with cancellation.
//
// Work is proportional to the length of text — around 1.4 s for a million
// words — so a caller processing untrusted or unbounded input has something
// real to cancel. The context is polled between pipeline stages.
func NewDocumentContext(ctx context.Context, text string, opts ...Option) (*Document, error) {
	cfg := config{tokenize: true, segment: true, tag: true, extract: true}
	for _, opt := range opts {
		opt(&cfg)
	}
	// Extraction is defined in terms of tagged tokens; asking for it without
	// them is a configuration error rather than something to silently ignore.
	if cfg.extract {
		cfg.tag = true
	}
	if cfg.tag {
		cfg.tokenize = true
	}

	doc := &Document{Text: text}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if cfg.segment {
		seg, err := segment.New()
		if err != nil {
			return nil, err
		}
		doc.sentences = seg.Segment(text)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if cfg.tokenize {
		tk := cfg.tokenizer
		if tk == nil {
			tk = defaultTokenizer()
		}
		doc.tokens = tk.Tokenize(text)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if cfg.tag {
		tg, err := tag.New()
		if err != nil {
			return nil, err
		}
		tg.TagTokens(doc.tokens)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if cfg.extract {
		ex, err := ner.New()
		if err != nil {
			return nil, err
		}
		ex.Label(doc.tokens)
		doc.entities = ex.Chunk(text, doc.tokens)
	}

	return doc, ctx.Err()
}
