package prose_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jdkato/prose/v3"
)

const sample = "Marie Curie studied in Paris. She won the Nobel Prize twice. " +
	"Her research on radioactivity changed modern physics."

func TestDocumentPipeline(t *testing.T) {
	doc, err := prose.NewDocument(sample)
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Tokens()) == 0 {
		t.Error("no tokens")
	}
	if len(doc.Sentences()) != 3 {
		t.Errorf("got %d sentences, want 3", len(doc.Sentences()))
	}
	if len(doc.Entities()) == 0 {
		t.Error("no entities")
	}
	for _, tok := range doc.Tokens() {
		if tok.Tag == "" {
			t.Errorf("token %q was not tagged", tok.Text)
			break
		}
	}
}

// TestDocumentPositions is the end-to-end position invariant: everything the
// pipeline hands back must locate itself in the source.
func TestDocumentPositions(t *testing.T) {
	texts := []string{
		sample,
		"He said “hi there” to me. They'll save $100 (roughly).",
		"It’s a test — with an em dash. And a second sentence!",
		"Dr. Smith went to Washington, D.C. on Jan. 5. He left.",
	}

	for _, src := range texts {
		doc, err := prose.NewDocument(src)
		if err != nil {
			t.Fatal(err)
		}
		for _, tok := range doc.Tokens() {
			if got := tok.In(src); got != tok.Text {
				t.Errorf("%q: token %q at [%d:%d] locates %q",
					src, tok.Text, tok.Start, tok.End(), got)
			}
		}
		for _, s := range doc.Sentences() {
			if got := s.In(src); got != s.Text {
				t.Errorf("%q: sentence %q at [%d:%d] locates %q",
					src, s.Text, s.Start, s.End(), got)
			}
		}
		for _, e := range doc.Entities() {
			if got := e.In(src); got != e.Text {
				t.Errorf("%q: entity %q at [%d:%d] locates %q",
					src, e.Text, e.Start, e.End(), got)
			}
		}
	}
}

func TestDocumentOptions(t *testing.T) {
	doc, err := prose.NewDocument(sample, prose.WithExtraction(false))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Entities()) != 0 {
		t.Errorf("extraction disabled but got %d entities", len(doc.Entities()))
	}
	if len(doc.Tokens()) == 0 {
		t.Error("tokens should still be produced")
	}

	doc, err = prose.NewDocument(sample, prose.WithSegmentation(false))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Sentences()) != 0 {
		t.Errorf("segmentation disabled but got %d sentences", len(doc.Sentences()))
	}

	// Extraction implies tagging, which implies tokenization.
	doc, err = prose.NewDocument(sample,
		prose.WithTokenization(false), prose.WithExtraction(true))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Tokens()) == 0 {
		t.Error("extraction should have forced tokenization back on")
	}
}

// TestContextCancellation checks that the context is actually honoured, not
// merely accepted.
func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := prose.NewDocumentContext(ctx, sample); !errors.Is(err, context.Canceled) {
		t.Errorf("got %v, want context.Canceled", err)
	}
}

func TestContextDeadline(t *testing.T) {
	// A large document so there is real work to interrupt.
	big := strings.Repeat(sample+" ", 5000)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	_, err := prose.NewDocumentContext(ctx, big)
	if err == nil {
		t.Skip("document finished before the deadline; nothing to assert")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("got %v, want context.DeadlineExceeded", err)
	}
}

// TestModelsAreShared guards the property the whole package depends on:
// creating many documents must not scale with model construction.
func TestModelsAreShared(t *testing.T) {
	if _, err := prose.NewDocument("warm up."); err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	for i := 0; i < 20; i++ {
		if _, err := prose.NewDocument(sample); err != nil {
			t.Fatal(err)
		}
	}
	elapsed := time.Since(start)

	// Rebuilding the models per document would put this in the seconds; if it
	// ever does, sharing has broken.
	if elapsed > time.Second {
		t.Errorf("20 documents took %v; models are probably not shared", elapsed)
	}
	t.Logf("20 documents in %v (%v each)", elapsed, elapsed/20)
}

func BenchmarkNewDocument(b *testing.B) {
	if _, err := prose.NewDocument("warm up."); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := prose.NewDocument(sample); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewDocumentNoNER(b *testing.B) {
	if _, err := prose.NewDocument("warm up."); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := prose.NewDocument(sample, prose.WithExtraction(false)); err != nil {
			b.Fatal(err)
		}
	}
}
