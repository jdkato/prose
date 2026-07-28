package segment_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/jdkato/prose/v3/segment"
)

func newSegmenter(tb testing.TB) *segment.Segmenter {
	tb.Helper()
	s, err := segment.New()
	if err != nil {
		tb.Fatalf("segment.New: %v", err)
	}
	return s
}

// TestSentenceOffsetsLocateTheirText is the counterpart to the tokenizer's
// offset test: because sentences are whitespace-trimmed, the trimming has to
// move the offsets too, or they point at whitespace the text no longer
// contains.
func TestSentenceOffsetsLocateTheirText(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"One sentence.",
		"One. Two. Three.",
		"  Leading space. And another.  ",
		"Dr. Smith went to Washington. He arrived on Jan. 5.",
		"She said “hello there.” Then she left.",
		"An ellipsis... and then more. Really.",
		"The F.B.I. investigated. Nothing came of it.",
		"Multi\nline\ntext. With a break.\n\nNew paragraph here.",
		"Unicode: naïve café résumé. Ünïcödé wörds!",
		"Emoji 🎉 sentence. Another 😀 one.",
		"Yahoo! is a company. So is Google.",
		"Repeated. Repeated. Repeated.",
	}

	seg := newSegmenter(t)
	for _, src := range cases {
		for _, s := range seg.Segment(src) {
			if got := s.In(src); got != s.Text {
				t.Errorf("input %q: sentence %q at [%d:%d] locates %q",
					src, s.Text, s.Start, s.End(), got)
			}
			if s.Start < 0 || s.End() > len(src) {
				t.Errorf("input %q: sentence %q has out-of-range offsets [%d:%d], len=%d",
					src, s.Text, s.Start, s.End(), len(src))
			}
			if s.Text != strings.TrimSpace(s.Text) {
				t.Errorf("input %q: sentence %q is not trimmed", src, s.Text)
			}
		}
	}
}

// TestSentencesAreOrderedAndDisjoint catches offsets that go backwards, which
// is what happens when a "Repeated. Repeated." style input is resolved by
// searching for the text instead of tracking position.
func TestSentencesAreOrderedAndDisjoint(t *testing.T) {
	src := "Repeated sentence. Something else. Repeated sentence. Done."
	prev := -1
	for _, s := range newSegmenter(t).Segment(src) {
		if s.Start < prev {
			t.Errorf("sentence %q starts at %d, before the previous ended at %d",
				s.Text, s.Start, prev)
		}
		if s.In(src) != s.Text {
			t.Errorf("sentence %q at [%d:%d] locates %q", s.Text, s.Start, s.End(), s.In(src))
		}
		prev = s.End()
	}
}

// TestOffsetsOverCorpus runs the invariant over real prose.
func TestOffsetsOverCorpus(t *testing.T) {
	seg := newSegmenter(t)
	var checked int
	for _, src := range corpusText(t) {
		for _, s := range seg.Segment(src) {
			checked++
			if s.In(src) != s.Text {
				t.Fatalf("offset mismatch: sentence %q at [%d:%d] locates %q in %q",
					s.Text, s.Start, s.End(), s.In(src), truncate(src))
			}
		}
	}
	t.Logf("verified offsets for %d sentences", checked)
}

// TestConcurrentSegmentation checks the claim that a Segmenter is safe to
// share. The punkt annotators hold a *sentences.Storage, and if any of them
// wrote to it during segmentation, sharing one Segmenter — which is the whole
// point of loading the model once — would be a data race.
//
// Run with -race for this to mean anything.
func TestConcurrentSegmentation(t *testing.T) {
	seg := newSegmenter(t)
	inputs := corpusText(t)
	if len(inputs) > 200 {
		inputs = inputs[:200]
	}

	want := make([][]segment.Sentence, len(inputs))
	for i, src := range inputs {
		want[i] = seg.Segment(src)
	}

	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i, src := range inputs {
				got := seg.Segment(src)
				if len(got) != len(want[i]) {
					t.Errorf("input %d: got %d sentences, want %d",
						i, len(got), len(want[i]))
					return
				}
				for j := range got {
					if got[j] != want[i][j] {
						t.Errorf("input %d sentence %d: got %+v, want %+v",
							i, j, got[j], want[i][j])
						return
					}
				}
			}
		}()
	}
	wg.Wait()
}

// TestSharedModelIsNotRebuilt checks that two Segmenters share one trained
// storage rather than each building their own.
func TestSharedModelIsNotRebuilt(t *testing.T) {
	a, err := segment.New()
	if err != nil {
		t.Fatal(err)
	}
	b, err := segment.New()
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatal("New returned the same *Segmenter twice; it should return distinct wrappers")
	}
	// Both must segment identically, which they only do if they share the
	// same trained storage.
	const src = "Dr. Smith went to Washington. He arrived on Jan. 5."
	if len(a.Segment(src)) != len(b.Segment(src)) {
		t.Error("two Segmenters disagree; they are not sharing a model")
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
		if json.Unmarshal(b, &sents) == nil {
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
		if s, ok := strings.CutPrefix(scan.Text(), "# text = "); ok {
			out = append(out, s)
		}
	}
	return out
}

func BenchmarkSegmenterNew(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := segment.New(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSegment(b *testing.B) {
	seg := newSegmenter(b)
	src := "Dr. Smith went to Washington. He arrived on Jan. 5. " +
		"She said “hello there.” Then she left. The F.B.I. investigated."
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		seg.Segment(src)
	}
}
