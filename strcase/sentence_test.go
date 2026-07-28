package strcase_test

import (
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
	{"README files", "Readme files"},
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
