package ner

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jdkato/prose/v3/tag"
	"github.com/jdkato/prose/v3/tokenize"
)

// The model was trained on exact feature strings, so any change to feature
// extraction or scoring must label this text identically. The golden file
// lists every non-O token as "start\tlabel\ttext".
func TestSherlockGolden(t *testing.T) {
	b, err := os.ReadFile("../testdata/sherlock.txt")
	if err != nil {
		t.Fatal(err)
	}
	src := string(b)

	toks := tokenize.New().Tokenize(src)
	tagger, err := tag.New()
	if err != nil {
		t.Fatal(err)
	}
	tagger.TagTokens(toks)
	ex, err := New()
	if err != nil {
		t.Fatal(err)
	}
	ex.Label(toks)

	var got []string
	for _, tok := range toks {
		if tok.Label != "O" {
			got = append(got, fmt.Sprintf("%d\t%s\t%s", tok.Start, tok.Label, tok.Text))
		}
	}

	f, err := os.Open("../testdata/sherlock.entities.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var want []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		want = append(want, sc.Text())
	}

	if len(got) != len(want) {
		t.Fatalf("labeled %d tokens, golden has %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("line %d: got %q, want %q", i+1, got[i], want[i])
		}
	}
}

func TestShape(t *testing.T) {
	cases := map[string]string{
		"42":       "number",
		"3.5":      "number",
		"...":      "punct",
		"hello":    "downcase",
		"Hello":    "upcase",
		"NASA":     "upcase",
		"McDonald": "upcase",
		"Jean-Luc": "upcase",
		"über":     "downcase",
	}
	for word, want := range cases {
		if got := shape(word); got != want {
			t.Errorf("shape(%q) = %q, want %q", word, got, want)
		}
	}
}

func TestEntities(t *testing.T) {
	src := "Marie Curie studied in Paris."
	toks := tokenize.New().Tokenize(src)
	tagger, _ := tag.New()
	tagger.TagTokens(toks)
	ex, _ := New()
	ex.Label(toks)

	var got []string
	for _, e := range ex.Chunk(src, toks) {
		if e.In(src) != e.Text {
			t.Errorf("%q is not at %d", e.Text, e.Start)
		}
		got = append(got, e.Text+" "+e.Label)
	}
	if want := "Marie Curie PERSON, Paris GPE"; strings.Join(got, ", ") != want {
		t.Errorf("got %q, want %q", strings.Join(got, ", "), want)
	}
}
