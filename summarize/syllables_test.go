package summarize

import (
	"bufio"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/jdkato/prose/v3/internal/strutil"
)

func TestSyllables(t *testing.T) {
	cases := strutil.ReadDataFile(filepath.Join(testdata, "syllables.json"))
	tests := make(map[string]int)

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		t.Error(err)
	}

	for word, count := range tests {
		if Syllables(word) != count {
			t.Errorf("Syllables(%s): got %d; expected %d", word, Syllables(word), count)
		}
	}

	total := 9462.0
	right := 0.0
	p := filepath.Join(testdata, "1-syllable-words.txt")
	right += testNSyllables(t, p, 1)

	p = filepath.Join(testdata, "2-syllable-words.txt")
	right += testNSyllables(t, p, 2)

	p = filepath.Join(testdata, "3-syllable-words.txt")
	right += testNSyllables(t, p, 3)

	p = filepath.Join(testdata, "4-syllable-words.txt")
	right += testNSyllables(t, p, 4)

	p = filepath.Join(testdata, "5-syllable-words.txt")
	right += testNSyllables(t, p, 5)

	p = filepath.Join(testdata, "6-syllable-words.txt")
	right += testNSyllables(t, p, 6)

	p = filepath.Join(testdata, "7-syllable-words.txt")
	right += testNSyllables(t, p, 7)

	ratio := math.Round(right / total)
	if ratio < 0.93 {
		t.Errorf("Less than 93%% accurate on NSyllables!")
	}
}

func testNSyllables(t *testing.T, fpath string, n int) float64 {
	file, err := os.Open(fpath)
	if err != nil {
		t.Error(err)
	}

	right := 0.0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if n == Syllables(word) {
			right++
		}
	}

	err = scanner.Err()
	if err != nil {
		t.Error(err)
	}

	err = file.Close()
	if err != nil {
		t.Error(err)
	}

	return right
}

func BenchmarkSyllables(b *testing.B) {
	cases := strutil.ReadDataFile(filepath.Join(testdata, "syllables.json"))
	tests := make(map[string]int)

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		b.Error(err)
	}

	for n := 0; n < b.N; n++ {
		for word := range tests {
			Syllables(word)
		}
	}
}

func BenchmarkSyllablesIn(b *testing.B) {
	cases := strutil.ReadDataFile(filepath.Join(testdata, "syllables.json"))
	tests := make(map[string]int)

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		b.Error(err)
	}

	for n := 0; n < b.N; n++ {
		for word := range tests {
			Syllables(word)
		}
	}
}
