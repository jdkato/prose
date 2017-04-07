package summarize

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/jdkato/prose/internal/util"
	"github.com/stretchr/testify/assert"
)

var syllableTests = []struct {
	in  string
	out int
}{
	{"what", 1},
	{"take", 1},
	{"taking", 2},
	{"dusted", 2},
	{"redo", 2},
	{"super", 2},
	// FIXME: {"worrying", 3},
	{"Maryland", 3},
	{"American", 3},
	{"disenfranchized", 5},
	{"Sophia", 2},
	{"misbehaving", 4},
}

func TestSyllables(t *testing.T) {
	for _, tt := range syllableTests {
		assert.Equal(t, tt.out, Syllables(tt.in), tt.in)
	}
}

type testCase struct {
	Text       string
	Sentences  float64
	Words      float64
	PolyWords  float64
	Characters float64
}

func TestSummarize(t *testing.T) {
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
