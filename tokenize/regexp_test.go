package tokenize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordPunctTokenizer(t *testing.T) {
	input, output := getWordData("word_punct.json")
	wordTokenizer := NewWordPunctTokenizer()
	for i, s := range input {
		assert.Equal(t, output[i], wordTokenizer.Tokenize(s))
	}
}

func TestNewRegexpTokenizer(t *testing.T) {
	input, _ := getWordData("word_punct.json")
	expected := NewWordPunctTokenizer()
	observed := NewRegexpTokenizer(`\w+|[^\w\s]+`, false, false)
	for _, s := range input {
		assert.Equal(t, expected.Tokenize(s), observed.Tokenize(s))
	}
}

func BenchmarkWordPunctTokenizer(b *testing.B) {
	word := NewWordPunctTokenizer()
	for n := 0; n < b.N; n++ {
		for _, s := range getWordBenchData() {
			word.Tokenize(s)
		}
	}
}
