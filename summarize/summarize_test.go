package summarize

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jdkato/prose/v3/internal/strutil"
)

var testdata = filepath.Join("..", "testdata")

type testCase struct {
	Text       string
	Sentences  float64
	Words      float64
	PolyWords  float64
	Characters float64

	AutomatedReadability float64
	ColemanLiau          float64
	FleschKincaid        float64
	GunningFog           float64
	SMOG                 float64
	LIX                  float64

	MeanGrade   float64
	StdDevGrade float64

	DaleChall   float64
	ReadingEase float64
}

func TestSummarizePrep(t *testing.T) {
	tests := make([]testCase, 0)
	cases := strutil.ReadDataFile(filepath.Join(testdata, "summarize.json"))

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		d := NewDocument(test.Text)

		if test.Sentences != d.NumSentences {
			t.Errorf("Sentences: got %0.2f; expected %0.2f", d.NumSentences, test.Sentences)
		}

		if test.Words != d.NumWords {
			t.Errorf("Words: got %0.2f; expected %0.2f", d.NumWords, test.Words)
		}

		if test.Characters != d.NumCharacters {
			t.Errorf("Characters: got %0.2f; expected %0.2f", d.NumCharacters, test.Characters)
		}
	}
}

func TestSummarize(_ *testing.T) {
	data := strutil.ReadDataFile(filepath.Join(testdata, "article.txt"))
	d := NewDocument(string(data))

	text := ""
	for _, paragraph := range d.Summary(7) {
		for _, s := range paragraph.Sentences {
			text += (s.Text + " ")
		}
		text += "\n\n"
	}
	fmt.Print(text)
}
