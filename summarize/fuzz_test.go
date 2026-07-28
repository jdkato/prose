package summarize_test

import (
	"testing"

	"github.com/jdkato/prose/v3/summarize"
)

// FuzzReadability replaces the go-fuzz target prose used with fuzzit.dev,
// which no longer exists. Go's native fuzzing runs the seed corpus on every
// `go test`, so these inputs act as regression cases even without a fuzzing
// run.
//
// The readability metrics divide by word, sentence and syllable counts, so
// degenerate input — empty strings, punctuation only, no sentence terminator —
// is where they are most likely to panic or produce NaN.
func FuzzReadability(f *testing.F) {
	seeds := []string{
		"",
		" ",
		".",
		"...",
		"A",
		"One sentence.",
		"Two sentences. Here is the second one.",
		"No terminator",
		"Ünïcödé wörds with áccents.",
		"🎉 emoji only 😀",
		"a a a a a a a a a a.",
		"Hyphenated-words and contractions don't break it.",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(_ *testing.T, text string) {
		d := summarize.NewDocument(text)

		// Each of these must return without panicking. Values are not
		// asserted: the metrics are undefined for degenerate input, but they
		// must not crash.
		d.AutomatedReadability()
		d.ColemanLiau()
		d.DaleChall()
		d.FleschKincaid()
		d.FleschReadingEase()
		d.GunningFog()
		d.LIX()
		d.SMOG()
	})
}
