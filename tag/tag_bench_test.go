package tag_test

import (
	"runtime"
	"testing"

	"github.com/jdkato/prose/v3/tag"
	"github.com/jdkato/prose/v3/tag/aptagmodel"
)

var benchWords = []string{
	"The", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog", ".",
	"Kubernetes", "orchestrates", "containerized", "workloads", "efficiently", ".",
	"I", "have", "read", "the", "documentation", "and", "it", "explains",
	"the", "configuration", "options", "clearly", ".",
}

// BenchmarkLoad measures one-time model construction.
func BenchmarkLoad(b *testing.B) {
	data := modelBytes(b)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := tag.FromBytes(data); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTag measures steady-state tagging, reported per token.
func BenchmarkTag(b *testing.B) {
	tg, err := tag.New()
	if err != nil {
		b.Fatal(err)
	}
	toks := make([]tag.Token, len(benchWords))
	for i, w := range benchWords {
		toks[i] = tag.Token{Text: w}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tg.TagTokens(toks)
	}
	b.ReportMetric(
		float64(b.Elapsed().Nanoseconds())/float64(b.N*len(toks)),
		"ns/token")
}

// BenchmarkTagAlloc is the same work through the allocating entry point, for
// callers that want a fresh slice each call.
func BenchmarkTagAlloc(b *testing.B) {
	tg, err := tag.New()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tg.Tag(benchWords)
	}
}

// BenchmarkTagParallel checks that the pooled scratch space actually scales.
func BenchmarkTagParallel(b *testing.B) {
	tg, err := tag.New()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		toks := make([]tag.Token, len(benchWords))
		for i, w := range benchWords {
			toks[i] = tag.Token{Text: w}
		}
		for pb.Next() {
			tg.TagTokens(toks)
		}
	})
}

// TestModelResidentMemory reports what the loaded model costs on the heap.
// The flat layout should keep it well under 6 MB.
func TestModelResidentMemory(t *testing.T) {
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)

	// FromBytes rather than New: New caches the decoded model for the process,
	// so it allocates nothing once another test has called it and this measures
	// whichever test happened to run first.
	tg, err := tag.FromBytes(aptagmodel.English)
	if err != nil {
		t.Fatal(err)
	}
	tg.Tag(benchWords) // force the pool to populate

	runtime.GC()
	runtime.ReadMemStats(&after)
	runtime.KeepAlive(tg)

	// Signed: HeapAlloc can fall across the two reads, and the unsigned
	// difference then wraps to a number that looks like terabytes.
	mb := float64(int64(after.HeapAlloc)-int64(before.HeapAlloc)) / (1 << 20)
	t.Logf("model resident: %.2f MB", mb)
	if mb > 6 {
		t.Errorf("model resident memory %.2f MB exceeds the 6 MB target", mb)
	}
}

func modelBytes(tb testing.TB) []byte {
	tb.Helper()
	// Round-trips the embedded model through the public API so the benchmark
	// measures the same path callers use.
	return embeddedModel
}
