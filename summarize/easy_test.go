package summarize

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/jdkato/prose/internal/util"
)

func BenchmarkEasyWordsLookupSlice(b *testing.B) {
	cases := util.ReadDataFile(filepath.Join(testdata, "syllables.json"))
	tests := make(map[string]int)
	util.CheckError(json.Unmarshal(cases, &tests))

	for n := 0; n < b.N; n++ {
		for word := range tests {
			util.StringInSlice(word, easyWords)
		}
	}
}

func BenchmarkEasyWordsLookupMap(b *testing.B) {
	cases := util.ReadDataFile(filepath.Join(testdata, "syllables.json"))
	tests := make(map[string]int)
	util.CheckError(json.Unmarshal(cases, &tests))

	easyWordsMap := util.SliceToMap(easyWords)
	for n := 0; n < b.N; n++ {
		for word := range tests {
			if _, ok := easyWordsMap[word]; ok {

			}
		}
	}
}
