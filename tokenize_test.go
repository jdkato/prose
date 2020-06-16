package prose

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testdata = "testdata"

func makeDoc(text string) (*Document, error) {
	return NewDocument(
		text,
		WithSegmentation(false),
		WithTagging(false),
		WithExtraction(false))
}

func getTokenText(doc *Document) []string {
	observed := []string{}
	tokens := doc.Tokens()
	for i := range tokens {
		observed = append(observed, tokens[i].Text)
	}
	return observed
}

func getWordData(file string) ([]string, [][]string) {
	in := readDataFile(filepath.Join(testdata, "treebank_sents.json"))
	out := readDataFile(filepath.Join(testdata, file))

	input := []string{}
	output := [][]string{}

	checkError(json.Unmarshal(in, &input))
	checkError(json.Unmarshal(out, &output))

	return input, output
}

func getWordBenchData() []string {
	in := readDataFile(filepath.Join(testdata, "treebank_sents.json"))
	input := []string{}
	checkError(json.Unmarshal(in, &input))
	return input
}

func TestTokenizationEmpty(t *testing.T) {
	doc, _ := makeDoc("")
	assert.Len(t, getTokenText(doc), 0)
}

func TestTokenizationSimple(t *testing.T) {
	doc, _ := makeDoc("Vale is a natural language linter that supports plain text, markup (Markdown, reStructuredText, AsciiDoc, and HTML), and source code comments. Vale doesn't attempt to offer a one-size-fits-all collection of rules—instead, it strives to make customization as easy as possible.")
	expected := []string{
		"Vale", "is", "a", "natural", "language", "linter", "that", "supports",
		"plain", "text", ",", "markup", "(", "Markdown", ",", "reStructuredText",
		",", "AsciiDoc", ",", "and", "HTML", ")", ",", "and", "source",
		"code", "comments", ".", "Vale", "does", "n't", "attempt", "to",
		"offer", "a", "one-size-fits-all", "collection", "of", "rules—instead",
		",", "it", "strives", "to", "make", "customization", "as", "easy", "as",
		"possible", "."}
	assert.Equal(t, expected, getTokenText(doc))
}

func TestTokenizationTreebank(t *testing.T) {
	input, output := getWordData("treebank_words.json")
	for i, s := range input {
		doc, _ := makeDoc(s)
		assert.Equal(t, output[i], getTokenText(doc))
	}
}

func TestTokenizationWeb(t *testing.T) {
	web := `Independent of current body composition, IGF-I levels at 5 yr were significantly
            associated with rate of weight gain between 0-2 yr (beta=0.19; P&lt;0.0005);
            and children who showed postnatal catch-up growth (i.e. those who showed gains in
            weight or length between 0-2 yr by >0.67 SD score) had higher IGF-I levels than other
			children (P=0.02; http://univ.edu.es/study.html) [20-22].`
	expected := []string{"Independent", "of", "current", "body", "composition", ",", "IGF-I",
		"levels", "at", "5", "yr", "were", "significantly", "associated", "with", "rate", "of",
		"weight", "gain", "between", "0-2", "yr", "(", "beta=0.19", ";", "P&lt;0.0005", ")", ";",
		"and", "children", "who", "showed", "postnatal", "catch-up", "growth", "(", "i.e.", "those",
		"who", "showed", "gains", "in", "weight", "or", "length", "between", "0-2", "yr", "by",
		">0.67", "SD", "score", ")", "had", "higher", "IGF-I", "levels", "than", "other", "children",
		"(", "P=0.02", ";", "http://univ.edu.es/study.html", ")", "[", "20-22", "]", "."}
	doc, _ := makeDoc(web)
	assert.Equal(t, expected, getTokenText(doc))
}

func TestTokenizationWebParagraph(t *testing.T) {
	web := `Independent of current body composition, IGF-I levels at 5 yr were significantly
            associated with rate of weight gain between 0-2 yr (beta=0.19; P&lt;0.0005);
            and children who showed postnatal catch-up growth (i.e. those who showed gains in
            weight or length between 0-2 yr by >0.67 SD score) had higher IGF-I levels than other
			children (P=0.02; http://univ.edu.es/study.html) [20-22].

			Independent of current body composition, IGF-I levels at 5 yr were significantly
            associated with rate of weight gain between 0-2 yr (beta=0.19; P&lt;0.0005);
            and children who showed postnatal catch-up growth (i.e. those who showed gains in
            weight or length between 0-2 yr by >0.67 SD score) had higher IGF-I levels than other
			children (P=0.02; http://univ.edu.es/study.html) [20-22].

			Independent of current body composition, IGF-I levels at 5 yr were significantly
            associated with rate of weight gain between 0-2 yr (beta=0.19; P&lt;0.0005);
            and children who showed postnatal catch-up growth (i.e. those who showed gains in
            weight or length between 0-2 yr by >0.67 SD score) had higher IGF-I levels than other
			children (P=0.02; http://univ.edu.es/study.html) [20-22].`

	expected := []string{"Independent", "of", "current", "body", "composition", ",", "IGF-I",
		"levels", "at", "5", "yr", "were", "significantly", "associated", "with", "rate", "of",
		"weight", "gain", "between", "0-2", "yr", "(", "beta=0.19", ";", "P&lt;0.0005", ")", ";",
		"and", "children", "who", "showed", "postnatal", "catch-up", "growth", "(", "i.e.", "those",
		"who", "showed", "gains", "in", "weight", "or", "length", "between", "0-2", "yr", "by",
		">0.67", "SD", "score", ")", "had", "higher", "IGF-I", "levels", "than", "other", "children",
		"(", "P=0.02", ";", "http://univ.edu.es/study.html", ")", "[", "20-22", "]", ".", "Independent", "of", "current", "body", "composition", ",", "IGF-I",
		"levels", "at", "5", "yr", "were", "significantly", "associated", "with", "rate", "of",
		"weight", "gain", "between", "0-2", "yr", "(", "beta=0.19", ";", "P&lt;0.0005", ")", ";",
		"and", "children", "who", "showed", "postnatal", "catch-up", "growth", "(", "i.e.", "those",
		"who", "showed", "gains", "in", "weight", "or", "length", "between", "0-2", "yr", "by",
		">0.67", "SD", "score", ")", "had", "higher", "IGF-I", "levels", "than", "other", "children",
		"(", "P=0.02", ";", "http://univ.edu.es/study.html", ")", "[", "20-22", "]", ".", "Independent", "of", "current", "body", "composition", ",", "IGF-I",
		"levels", "at", "5", "yr", "were", "significantly", "associated", "with", "rate", "of",
		"weight", "gain", "between", "0-2", "yr", "(", "beta=0.19", ";", "P&lt;0.0005", ")", ";",
		"and", "children", "who", "showed", "postnatal", "catch-up", "growth", "(", "i.e.", "those",
		"who", "showed", "gains", "in", "weight", "or", "length", "between", "0-2", "yr", "by",
		">0.67", "SD", "score", ")", "had", "higher", "IGF-I", "levels", "than", "other", "children",
		"(", "P=0.02", ";", "http://univ.edu.es/study.html", ")", "[", "20-22", "]", "."}

	doc, _ := makeDoc(web)
	assert.Equal(t, expected, getTokenText(doc))
}

func TestTokenizationTwitter(t *testing.T) {
	doc, _ := makeDoc("@twitter, what time does it start :-)")
	expected := []string{"@twitter", ",", "what", "time", "does", "it", "start", ":-)"}
	assert.Equal(t, expected, getTokenText(doc))

	doc, _ = makeDoc("Mr. James plays basketball in the N.B.A., do you?")
	expected = []string{
		"Mr.", "James", "plays", "basketball", "in", "the", "N.B.A.", ",",
		"do", "you", "?"}
	assert.Equal(t, expected, getTokenText(doc))

	doc, _ = makeDoc("ˌˌ kill the last letter")
	expected = []string{"ˌˌ", "kill", "the", "last", "letter"}
	assert.Equal(t, expected, getTokenText(doc))

	doc, _ = makeDoc("ˌˌˌ kill the last letter")
	expected = []string{"ˌˌˌ", "kill", "the", "last", "letter"}
	assert.Equal(t, expected, getTokenText(doc))

	doc, _ = makeDoc("March. July. March. June. January.")
	expected = []string{
		"March", ".", "July", ".", "March", ".", "June", ".", "January", "."}
	assert.Equal(t, expected, getTokenText(doc))
}

func BenchmarkTokenization(b *testing.B) {
	in := readDataFile(filepath.Join(testdata, "sherlock.txt"))
	text := string(in)
	for n := 0; n < b.N; n++ {
		_, err := makeDoc(text)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkTokenizationSimple(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, s := range getWordBenchData() {
			_, err := makeDoc(s)
			if err != nil {
				panic(err)
			}
		}
	}
}
