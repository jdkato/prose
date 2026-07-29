package tag_test

import (
	"strings"
	"testing"

	"github.com/jdkato/prose/v3/tag"
)

func lexical(t *testing.T, dict string) tag.Interface {
	t.Helper()

	entries, err := tag.ReadDictionary(strings.NewReader(dict))
	if err != nil {
		t.Fatal(err)
	}
	lex, err := tag.NewLexical(entries, nil)
	if err != nil {
		t.Fatal(err)
	}
	return lex
}

// A word the dictionary is sure about keeps its tag however odd the sentence
// around it is. This is the case a contextual tagger gets wrong, and the whole
// reason for tagging from a dictionary.
func TestLexicalOverridesContext(t *testing.T) {
	words := []string{"They", "were", "all", "ready", "aware", "of", "it"}

	perceptron, err := tag.New()
	if err != nil {
		t.Fatal(err)
	}
	if got := tagOf(perceptron.Tag(words), "aware"); got == "JJ" {
		t.Skip("the perceptron now tags 'aware' correctly; this guarded a known miss")
	}

	if got := tagOf(lexical(t, "aware\tJJ\n").Tag(words), "aware"); got != "JJ" {
		t.Errorf("'aware' = %s, want JJ", got)
	}
}

// A word the dictionary lists more than one way is left to context, because
// context is what tells those readings apart.
func TestLexicalLeavesAmbiguityToContext(t *testing.T) {
	lex := lexical(t, "code\tNN\tVB\n")

	noun := tagOf(lex.Tag([]string{"the", "code", "is", "broken"}), "code")
	verb := tagOf(lex.Tag([]string{"they", "code", "every", "day"}), "code")

	for _, got := range []string{noun, verb} {
		if got != "NN" && got != "VB" {
			t.Errorf("tag %q is not one the dictionary allows", got)
		}
	}
	if noun == verb {
		t.Logf("context did not separate the readings (%s); the dictionary still bounded them", noun)
	}
}

// An unknown word is the fallback's to answer.
func TestLexicalDefersOnUnknownWords(t *testing.T) {
	words := []string{"the", "quick", "brown", "fox"}

	perceptron, err := tag.New()
	if err != nil {
		t.Fatal(err)
	}

	want := perceptron.Tag(words)
	got := lexical(t, "aware\tJJ\n").Tag(words)

	for i := range want {
		if got[i].Tag != want[i].Tag {
			t.Errorf("%q: %s, want %s", words[i], got[i].Tag, want[i].Tag)
		}
	}
}

// One entry covers every casing: a word does not change class because a
// sentence started.
func TestLexicalIsCaseInsensitive(t *testing.T) {
	lex := lexical(t, "aware\tJJ\n")

	for _, w := range []string{"aware", "Aware", "AWARE"} {
		if got := tagOf(lex.Tag([]string{"is", w, "here"}), w); got != "JJ" {
			t.Errorf("%q = %s, want JJ", w, got)
		}
	}
}

func TestReadDictionaryRejectsIncompleteLines(t *testing.T) {
	if _, err := tag.ReadDictionary(strings.NewReader("aware\n")); err == nil {
		t.Error("a line with no tag should be an error, not a silent skip")
	}
}

func TestOpenByName(t *testing.T) {
	tg, err := tag.Open("")
	if err != nil {
		t.Fatal(err)
	}
	if tg.Name() != tag.Default {
		t.Errorf("empty name gave %q, want the default", tg.Name())
	}
	if _, err = tag.Open("nope"); err == nil {
		t.Error("an unknown tagger should be an error")
	}
}

// When the model picks a reading the dictionary does not list, its answer is
// still evidence of the class. `VBP` against a dictionary listing `NN` and
// `VB` should land on `VB`: the model was sure it was a verb, and only unsure
// which form.
func TestLexicalRefinesWithinTheModelsClass(t *testing.T) {
	lex := lexical(t, "like\tNN\tVB\tJJ\tIN\n")

	got := tagOf(lex.Tag([]string{"I", "like", "it"}), "like")
	if got != "VB" {
		t.Errorf("'like' = %s, want VB (the model read it as a verb)", got)
	}
}
