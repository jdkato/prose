package tokenize_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jdkato/prose/v3/tokenize"
)

// TestOffsetsLocateTheirText is the property the whole positions feature rests
// on: for every token, the offsets must locate exactly the token's text in the
// original input.
//
// This is the test that the tsawler/prose fork fails. It computes offsets
// against a sanitized copy of the input, so `He said “hi there” to me.` yields
// tokens whose offsets point into the middle of multi-byte runes — 7 of 9
// tokens misplaced. The failure is invisible on pure-ASCII input, which is why
// the cases below deliberately include curly quotes, entities, emoji, and
// other multi-byte text.
func TestOffsetsLocateTheirText(t *testing.T) {
	cases := []string{
		"",
		" ",
		"Hello.",
		`He said "hi there" to me.`,
		"He said “hi there” to me.",
		"It&rsquo;s a test of offsets.",
		"It’s a test — with an em dash.",
		"They'll say don't, won't, and I'm.",
		"Dr. Smith went to Washington, D.C. on Jan. 5.",
		"$100 (roughly) [see note] — 50% off!",
		"emoji 🎉 and 😀 mixed with ASCII",
		"Ünïcödé wörds with áccents and ñ.",
		"multiple    spaces\tand\ttabs\nand\nnewlines",
		"trailing whitespace   ",
		"   leading whitespace",
		"a",
		"...",
		"(-8 :) 8-D",
		"naïve café résumé",
		"“Nested ‘quotes’ here,” she said.",
		"Repeated repeated repeated words words.",
		"http://example.com/path?q=1 and user@example.com",
	}

	tk := tokenize.New()
	for _, src := range cases {
		for _, tok := range tk.Tokenize(src) {
			if got := tok.In(src); got != tok.Text {
				t.Errorf("input %q: token %q at [%d:%d] locates %q",
					src, tok.Text, tok.Start, tok.End(), got)
			}
			if tok.Start < 0 || tok.End() > len(src) {
				t.Errorf("input %q: token %q has out-of-range offsets [%d:%d], len=%d",
					src, tok.Text, tok.Start, tok.End(), len(src))
			}
		}
	}
}

// TestOffsetsAreOrderedAndDisjoint checks that tokens march forward through
// the source without overlapping. A splitter that re-derives offsets by
// searching for a token's text will violate this on repeated words.
func TestOffsetsAreOrderedAndDisjoint(t *testing.T) {
	src := "the cat sat on the mat, and then the cat sat again on the mat."
	prev := -1
	for _, tok := range tokenize.New().Tokenize(src) {
		if tok.Start < prev {
			t.Errorf("token %q starts at %d, before the previous token ended at %d",
				tok.Text, tok.Start, prev)
		}
		if got := tok.In(src); got != tok.Text {
			t.Errorf("token %q at [%d:%d] locates %q", tok.Text, tok.Start, tok.End(), got)
		}
		prev = tok.End()
	}
}

// TestOffsetsOverCorpus runs the invariant over every sentence available,
// including the UD treebanks when PROSE_CORPORA is set.
func TestOffsetsOverCorpus(t *testing.T) {
	tk := tokenize.New()
	var checked int

	for _, src := range corpusText(t) {
		for _, tok := range tk.Tokenize(src) {
			checked++
			if tok.In(src) != tok.Text {
				t.Fatalf("offset mismatch in %q: token %q at [%d:%d] locates %q",
					truncate(src), tok.Text, tok.Start, tok.End(), tok.In(src))
			}
		}
	}
	t.Logf("verified offsets for %d tokens", checked)
	if checked < 1000 {
		t.Logf("note: only %d tokens checked; set PROSE_CORPORA for a fuller run", checked)
	}
}

// TestTokenizationBehaviour pins the splitting rules themselves, so a future
// change to the affix sets cannot silently alter tokenization.
func TestTokenizationBehaviour(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"They'll save and invest more.", []string{"They", "'ll", "save", "and", "invest", "more", "."}},
		{"Don't stop.", []string{"Do", "n't", "stop", "."}},
		{"$100 well spent.", []string{"$", "100", "well", "spent", "."}},
		{"Well) done", []string{"Well", ")", "done"}},
		{"He said “hi” now.", []string{"He", "said", "“", "hi", "”", "now", "."}},
		{"It’s fine.", []string{"It", "’s", "fine", "."}},
	}

	tk := tokenize.New()
	for _, tc := range cases {
		toks := tk.Tokenize(tc.in)
		got := make([]string, len(toks))
		for i, tk := range toks {
			got[i] = tk.Text
		}
		if strings.Join(got, "|") != strings.Join(tc.want, "|") {
			t.Errorf("Tokenize(%q)\n  got  %q\n  want %q", tc.in, got, tc.want)
		}
	}
}

func truncate(s string) string {
	if len(s) > 60 {
		return s[:60] + "..."
	}
	return s
}

func corpusText(t testing.TB) []string {
	t.Helper()
	var out []string

	if b, err := os.ReadFile("../testdata/treebank_sents.json"); err == nil {
		var sents []string
		if err := json.Unmarshal(b, &sents); err == nil {
			out = append(out, sents...)
		}
	}

	dir := os.Getenv("PROSE_CORPORA")
	if dir == "" {
		return out
	}
	files, _ := filepath.Glob(filepath.Join(dir, "*.conllu"))
	for _, f := range files {
		out = append(out, conlluText(t, f)...)
	}
	return out
}

// conlluText pulls the "# text = ..." lines out of a UD treebank, giving real
// sentences with their original punctuation and Unicode intact.
func conlluText(t testing.TB, path string) []string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("opening %s: %v", path, err)
	}
	defer f.Close()

	var out []string
	scan := bufio.NewScanner(f)
	scan.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scan.Scan() {
		line := scan.Text()
		if s, ok := strings.CutPrefix(line, "# text = "); ok {
			out = append(out, s)
		}
	}
	return out
}
