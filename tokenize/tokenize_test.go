package tokenize

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jdkato/prose/internal/util"
	"github.com/stretchr/testify/assert"
)

var testdata = filepath.Join("..", "testdata")

func getWordData(file string) ([]string, [][]string) {
	in := util.ReadDataFile(filepath.Join(testdata, "treebank_sents.json"))
	out := util.ReadDataFile(filepath.Join(testdata, file))

	input := []string{}
	output := [][]string{}

	util.CheckError(json.Unmarshal(in, &input))
	util.CheckError(json.Unmarshal(out, &output))

	return input, output
}

func getWordBenchData() []string {
	in := util.ReadDataFile(filepath.Join(testdata, "treebank_sents.json"))
	input := []string{}
	util.CheckError(json.Unmarshal(in, &input))
	return input
}

func ExampleNewWordBoundaryTokenizer() {
	t := NewWordBoundaryTokenizer()
	fmt.Println(t.Tokenize("They'll save and invest more."))
	// Output: [They'll save and invest more]
}

func ExampleNewWordPunctTokenizer() {
	t := NewWordPunctTokenizer()
	fmt.Println(t.Tokenize("They'll save and invest more."))
	// Output: [They ' ll save and invest more .]
}

func ExampleNewTreebankWordTokenizer() {
	t := NewTreebankWordTokenizer()
	fmt.Println(t.Tokenize("They'll save and invest more."))
	// Output: [They 'll save and invest more .]
}

func ExampleNewBlanklineTokenizer() {
	t := NewBlanklineTokenizer()
	fmt.Println(t.Tokenize("They'll save and invest more.\n\nThanks!"))
	// Output: [They'll save and invest more. Thanks!]
}

func TestTextToWords(t *testing.T) {
	text := "Vale is a natural language linter that supports plain text, markup (Markdown, reStructuredText, AsciiDoc, and HTML), and source code comments. Vale doesn't attempt to offer a one-size-fits-all collection of rules—instead, it strives to make customization as easy as possible."
	expected := []string{
		"Vale", "is", "a", "natural", "language", "linter", "that", "supports",
		"plain", "text", ",", "markup", "(", "Markdown", ",", "reStructuredText",
		",", "AsciiDoc", ",", "and", "HTML", ")", ",", "and", "source", "code",
		"comments", ".", "Vale", "does", "n't", "attempt", "to", "offer", "a",
		"one-size-fits-all", "collection", "of", "rules—instead", ",", "it",
		"strives", "to", "make", "customization", "as", "easy", "as", "possible",
		"."}
	assert.Equal(t, expected, TextToWords(text))
}
