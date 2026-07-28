package tokenize_test

import (
	"testing"

	"github.com/jdkato/prose/v3/tokenize"
)

// FuzzOffsetInvariant fuzzes the property everything in prose depends on: a
// token's offsets must locate exactly its text in the source.
//
// This is the right thing to fuzz because the ways to break it are all
// input-dependent — multi-byte runes, normalization that changes byte length,
// repeated spans, degenerate whitespace. A hand-written table catches the
// cases you thought of; this catches the ones you did not.
func FuzzOffsetInvariant(f *testing.F) {
	seeds := []string{
		"",
		" ",
		"Hello.",
		`He said "hi there" to me.`,
		"He said “hi there” to me.",
		"It&rsquo;s a test of offsets.",
		"It’s a test — with an em dash.",
		"They'll say don't, won't, and I'm.",
		"$100 (roughly) [see note] — 50% off!",
		"emoji 🎉 and 😀 mixed with ASCII",
		"Ünïcödé wörds with áccents and ñ.",
		"multiple    spaces\tand\ttabs\nand\nnewlines",
		"the cat sat on the mat, the cat sat",
		"\u200b\u200bzero width",
		"...",
		"(-8 :) 8-D",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	tk := tokenize.New()

	f.Fuzz(func(t *testing.T, src string) {
		prev := 0
		for _, tok := range tk.Tokenize(src) {
			if tok.Start < 0 || tok.End() > len(src) {
				t.Fatalf("token %q has offsets [%d:%d] outside a %d-byte input",
					tok.Text, tok.Start, tok.End(), len(src))
			}
			if got := tok.In(src); got != tok.Text {
				t.Fatalf("token %q at [%d:%d] locates %q",
					tok.Text, tok.Start, tok.End(), got)
			}
			if tok.Start < prev {
				t.Fatalf("token %q starts at %d, before the previous token ended at %d",
					tok.Text, tok.Start, prev)
			}
			prev = tok.End()
		}
	})
}
