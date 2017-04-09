package summarize

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/jdkato/prose/internal/util"
	"github.com/jdkato/syllables"
	"github.com/stretchr/testify/assert"
)

func TestSyllables(t *testing.T) {
	cases := util.ReadDataFile(filepath.Join(testdata, "syllables.json"))
	tests := make(map[string]int)
	util.CheckError(json.Unmarshal(cases, &tests))

	for word, count := range tests {
		assert.Equal(t, count, Syllables(word), word)
	}
}

func BenchmarkSyllables(b *testing.B) {
	cases := util.ReadDataFile(filepath.Join(testdata, "syllables.json"))
	tests := make(map[string]int)
	util.CheckError(json.Unmarshal(cases, &tests))

	for n := 0; n < b.N; n++ {
		for word := range tests {
			Syllables(word)
		}
	}
}

func BenchmarkSyllablesIn(b *testing.B) {
	cases := util.ReadDataFile(filepath.Join(testdata, "syllables.json"))
	tests := make(map[string]int)
	util.CheckError(json.Unmarshal(cases, &tests))

	for n := 0; n < b.N; n++ {
		for word := range tests {
			syllables.In(word)
		}
	}
}
