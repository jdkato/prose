package tokenize_test

import (
	"testing"

	"github.com/jdkato/prose/v3/tag"
	"github.com/jdkato/prose/v3/tokenize"
)

// Positions originate in the tokenizer, but they are only useful if they
// survive everything downstream. The tagger writes Tag and must leave Start
// and Text alone.
func TestPositionsSurviveTagging(t *testing.T) {
	src := "He said “hi there” to me. They'll save $100 (roughly), " +
		"but Dr. Smith won't — it’s 50% off!"

	toks := tokenize.New().Tokenize(src)

	before := make([]tokenize.Token, len(toks))
	copy(before, toks)

	tg, err := tag.New()
	if err != nil {
		t.Fatal(err)
	}
	tg.TagTokens(toks)

	for i, tok := range toks {
		if tok.Start != before[i].Start || tok.Text != before[i].Text {
			t.Errorf("token %d changed across tagging: was %q@%d, now %q@%d",
				i, before[i].Text, before[i].Start, tok.Text, tok.Start)
		}
		if tok.In(src) != tok.Text {
			t.Errorf("token %d %q at [%d:%d] locates %q after tagging",
				i, tok.Text, tok.Start, tok.End(), tok.In(src))
		}
		if tok.Tag == "" {
			t.Errorf("token %d %q was not tagged", i, tok.Text)
		}
	}
}

// BenchmarkTokenizeAndTag measures the combined path a consumer like Vale
// actually uses.
func BenchmarkTokenizeAndTag(b *testing.B) {
	src := "He said “hi there” to me. They'll save $100 (roughly), " +
		"but Dr. Smith won't — it’s 50% off!"

	tk := tokenize.New()
	tg, err := tag.New()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		toks := tk.Tokenize(src)
		tg.TagTokens(toks)
	}
}
