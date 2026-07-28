// Package token defines the token type shared across prose's pipeline.
//
// It is deliberately dependency-free so that tokenizing, tagging, and entity
// extraction can all speak the same type without any of them pulling in the
// others' models.
package token

// Token is a piece of text with its position in the source.
//
// # The position invariant
//
// Start is a byte offset into the text the token came from, and Text is
// exactly the bytes found there:
//
//	src[tok.Start:tok.End()] == tok.Text
//
// Everything in prose that produces tokens maintains this, and
// tokenize.TestOffsetsLocateTheirText enforces it over the full corpus.
//
// This is load-bearing and easy to lose. The common way to break it is to
// tokenize a *normalized* copy of the input — replacing curly quotes with
// straight ones, decoding HTML entities — and report offsets into that copy.
// Those replacements change byte lengths ("“" is three bytes, `"` is
// one), so every offset after the first replacement points somewhere wrong.
// prose therefore never rewrites the input before tokenizing; normalization
// that would shift bytes is the caller's business, done before prose sees the
// text.
type Token struct {
	// Text is the token's content, verbatim from the source.
	Text string

	// Tag is the part-of-speech tag, empty until tagged.
	Tag string

	// Label is the IOB entity label, empty until extracted.
	Label string

	// Start is the byte offset of Text within the source.
	Start int
}

// End returns the byte offset just past this token in the source.
//
// Derived rather than stored: it cannot drift out of sync with Text, and it
// keeps Token one word smaller, which matters when a document is a few
// hundred thousand tokens.
func (t Token) End() int { return t.Start + len(t.Text) }

// In returns the token's text as located in src.
//
// Intended for assertions and debugging: for any token produced by prose,
// tok.In(src) == tok.Text. It returns "" rather than panicking if the offsets
// do not fit src, so a corrupted token surfaces as a mismatch instead of a
// crash.
func (t Token) In(src string) string {
	if t.Start < 0 || t.End() > len(src) {
		return ""
	}
	return src[t.Start:t.End()]
}
