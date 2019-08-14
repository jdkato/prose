package summarize

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jdkato/prose/internal/util"
	"github.com/stretchr/testify/assert"
)

var testdata = filepath.Join("..", "testdata")

func check(expected, observed float64) bool {
	return fmt.Sprintf("%0.2f", expected) == fmt.Sprintf("%0.2f", observed)
}

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
	cases := util.ReadDataFile(filepath.Join(testdata, "summarize.json"))

	util.CheckError(json.Unmarshal(cases, &tests))
	for _, test := range tests {
		d := NewDocument(test.Text)
		assert.Equal(t, test.Sentences, d.NumSentences)
		assert.Equal(t, test.Words, d.NumWords)
		assert.Equal(t, test.Characters, d.NumCharacters)
	}
}

func TestSummarize(t *testing.T) {
	data := util.ReadDataFile(filepath.Join(testdata, "article.txt"))
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
