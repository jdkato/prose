package strcase_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/jdkato/prose/v3/internal/strutil"
	"github.com/jdkato/prose/v3/strcase"
)

var testdata = filepath.Join("..", "testdata")

type testCase struct {
	Input  string
	Expect string
}

var vocabTitles = []testCase{
	{"Getting started with iOS 15", "Getting Started With iOS 15"},
	{"Understanding json and yaml", "Understanding JSON and YAML"},
	{"Develop File-Proxy Plugin", "Develop File-Proxy Plugin"},
	{"b. Next title text", "b. Next Title Text"},
	{"vale ale", "Vale ale"},
	{"New Repository and Project", "New Repository and Project"},
}

func TestTitleVocab(t *testing.T) {
	tc := strcase.NewTitleConverter(strcase.APStyle,
		strcase.UsingVocab([]string{
			"iOS",
			"JSON",
			"YAML",
			`\bale\b`,
			`[Rr]epo`,
		}),
		strcase.UsingPrefix(`^[a-z]\.\s`))

	for _, test := range vocabTitles {
		sent := tc.Convert(test.Input)
		if test.Expect != sent {
			t.Fatalf("Got '%s'; expected '%s'", sent, test.Expect)
		}
	}
}

func TestAP(t *testing.T) {
	tests := make([]testCase, 0)
	cases := strutil.ReadDataFile(filepath.Join(testdata, "AP.json"))

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		t.Error(err)
	}

	tc := strcase.NewTitleConverter(strcase.APStyle)
	for _, test := range tests {
		title := tc.Convert(test.Input)
		if test.Expect != title {
			t.Fatalf("Got '%s'; expected '%s'", title, test.Expect)
		}
	}
}

func TestChicago(t *testing.T) {
	tests := make([]testCase, 0)
	cases := strutil.ReadDataFile(filepath.Join(testdata, "Chicago.json"))

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		t.Error(err)
	}

	tc := strcase.NewTitleConverter(strcase.ChicagoStyle)
	for _, test := range tests {
		title := tc.Convert(test.Input)
		if test.Expect != title {
			t.Fatalf("Got '%s'; expected '%s'", title, test.Expect)
		}
	}
}

func BenchmarkTitle(b *testing.B) {
	tests := make([]testCase, 0)
	cases := strutil.ReadDataFile(filepath.Join(testdata, "title.json"))

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		b.Error(err)
	}

	tc := strcase.NewTitleConverter(strcase.APStyle)
	for n := 0; n < b.N; n++ {
		for _, test := range tests {
			_ = tc.Convert(test.Input)
		}
	}
}
