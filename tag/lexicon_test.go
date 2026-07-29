package tag_test

import (
	"testing"

	"github.com/jdkato/prose/v3/tag"
)

// A lexicon fixes words the model gets wrong from context. The sentences a
// grammar rule fires on are not idiomatic, which is exactly where the model's
// context features stop being evidence.
func TestLexiconOverridesPrediction(t *testing.T) {
	words := []string{"They", "were", "all", "ready", "aware", "of", "it"}

	plain, err := tag.New()
	if err != nil {
		t.Fatal(err)
	}
	if got := tagOf(plain.Tag(words), "aware"); got != "NN" {
		t.Skipf("model now tags 'aware' as %s; this test guarded a known miss", got)
	}

	fixed, err := tag.New(tag.WithLexicon(tag.Lexicon{"aware": "JJ"}))
	if err != nil {
		t.Fatal(err)
	}
	if got := tagOf(fixed.Tag(words), "aware"); got != "JJ" {
		t.Errorf("with lexicon, 'aware' = %s, want JJ", got)
	}
}

// One entry should cover every casing, since a word does not change class
// because a sentence started or someone shouted.
func TestLexiconIsCaseInsensitive(t *testing.T) {
	tg, err := tag.New(tag.WithLexicon(tag.Lexicon{"available": "JJ"}))
	if err != nil {
		t.Fatal(err)
	}

	for _, w := range []string{"available", "Available", "AVAILABLE"} {
		if got := tagOf(tg.Tag([]string{"data", "is", w, "now"}), w); got != "JJ" {
			t.Errorf("%q tagged %s, want JJ", w, got)
		}
	}
}

// An exact-case entry should win over the folded one.
func TestLexiconPrefersExactCase(t *testing.T) {
	tg, err := tag.New(tag.WithLexicon(tag.Lexicon{
		"us": "PRP",
		"US": "NNP",
	}))
	if err != nil {
		t.Fatal(err)
	}

	if got := tagOf(tg.Tag([]string{"send", "US", "mail"}), "US"); got != "NNP" {
		t.Errorf("'US' tagged %s, want NNP", got)
	}
	if got := tagOf(tg.Tag([]string{"send", "us", "mail"}), "us"); got != "PRP" {
		t.Errorf("'us' tagged %s, want PRP", got)
	}
}

// No lexicon must behave exactly as before.
func TestNoLexiconUnchanged(t *testing.T) {
	words := []string{"The", "quick", "brown", "fox", "jumps"}

	a, err := tag.New()
	if err != nil {
		t.Fatal(err)
	}
	b, err := tag.New(tag.WithLexicon(nil))
	if err != nil {
		t.Fatal(err)
	}

	for i, tok := range a.Tag(words) {
		if other := b.Tag(words)[i]; other.Tag != tok.Tag {
			t.Errorf("%q: %s vs %s", tok.Text, tok.Tag, other.Tag)
		}
	}
}

func tagOf(tokens []tag.Token, word string) string {
	for _, t := range tokens {
		if t.Text == word {
			return t.Tag
		}
	}
	return ""
}
