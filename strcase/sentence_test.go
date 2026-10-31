package strcase_test

import (
	"fmt"
	"testing"

	"github.com/jdkato/prose/v3/strcase"
)

var cases = []testCase{
	{"", ""},
	{"1. An important heading", "1. An important heading"},
	{"getting started with Vale server", "Getting started with vale server"},
	{"Lession 1: getting started with vale server", "Lession 1: Getting started with vale server"},
	{"Top-Level ideas", "Top-level ideas"},
	{"Intro to the top-level idEas", "Intro to the top-level ideas"},
	{"build-backend tools", "Build-backend tools"},
	{"README files", "README files"},
	{"TAP endpoint", "TAP endpoint"},
	{"Configure the CLI", "Configure the CLI"},
	{"Using the TAP Service", "Using the TAP service"},
	{"Querying the SIA Service", "Querying the SIA service"},
	{"Two APIs and the API's docs", "Two APIs and the API's docs"},
	{"Is it on", "Is it on"},
}

var vocabCases = []testCase{
	{"Getting started with vale server", "Getting started with Vale Server"},
	{"Issue triage", "Issue triage"},
	{"macOS 15: What's new", "macOS 15: What's new"},
	{"Configuration", "Configuration"},
	{"Six NASA Instruments Will Fly to Moon on Intuitive Machines Lander", "Six NASA instruments will fly to Moon on intuitive machines lander"},
	{"b. Next title text", "b. Next title text"},
	{"Axon Server connection", "Axon Server connection"},
	{"vale ale", "Vale ale"},
}

func TestSentence(t *testing.T) {
	tc := strcase.NewSentenceConverter()
	for _, test := range cases {
		sent := tc.Convert(test.Input)
		if test.Expect != sent {
			t.Fatalf("Got '%s'; expected '%s'", sent, test.Expect)
		}
	}
}

func TestVocab(t *testing.T) {
	tc := strcase.NewSentenceConverter(
		strcase.UsingVocab([]string{
			"Vale Server",
			`\bI\b`,
			"macOS",
			"[Cc]onfig",
			"NASA",
			"Moon",
			"Axon",
			"Axon Server",
			`\bale\b`,
		}),
		strcase.UsingPrefix(`^[a-z]\.\s`))

	for _, test := range vocabCases {
		sent := tc.Convert(test.Input)
		if test.Expect != sent {
			t.Fatalf("Got '%s'; expected '%s'", sent, test.Expect)
		}
	}
}

// Vocab patterns are compiled once, so a conversion shouldn't allocate more
// for a longer list of exceptions.
func TestVocabCostIsFlat(t *testing.T) {
	allocs := func(n int) float64 {
		vocab := make([]string, n)
		for i := range vocab {
			vocab[i] = fmt.Sprintf("Alpha%02d", i)
		}
		sc := strcase.NewSentenceConverter(strcase.UsingVocab(vocab))
		tc := strcase.NewTitleConverter(strcase.APStyle, strcase.UsingVocab(vocab))
		return testing.AllocsPerRun(20, func() {
			_ = sc.Convert("Get started with the API fast")
			_ = tc.Convert("Get started with the API fast")
		})
	}

	one, many := allocs(1), allocs(60)
	if many > 2*one {
		t.Fatalf("60 exceptions allocated %.0f times per run; 1 allocated %.0f", many, one)
	}
}
