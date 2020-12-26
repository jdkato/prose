package prose

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"
)

var testdata = "testdata"

func checkTokens(t *testing.T, tokens []*Token, expected []string, name string) {
	observed := []string{}
	for i := range tokens {
		observed = append(observed, tokens[i].Text)
	}
	if !reflect.DeepEqual(observed, expected) {
		t.Errorf("%v: unexpected tokens", name)
	}
}

func checkCase(t *testing.T, doc *Document, expected []string, name string) {
	tokens := getTokenText(doc)
	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("%v: unexpected tokens", name)
	}
}

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

	l := len(getTokenText(doc))
	if l != 0 {
		t.Errorf("TokenizationEmpty() expected = 0, got = %v", l)
	}
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
	checkCase(t, doc, expected, "TokenizationSimple()")
}

func TestTokenizationTreebank(t *testing.T) {
	input, output := getWordData("treebank_words.json")
	for i, s := range input {
		doc, _ := makeDoc(s)
		tokens := getTokenText(doc)
		if !reflect.DeepEqual(tokens, output[i]) {
			t.Errorf("TokenizationTreebank(): unexpected tokens")
		}
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
	checkCase(t, doc, expected, "TokenizationWeb()")
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
	checkCase(t, doc, expected, "TokenizationWebParagraph()")
}

func TestTokenizationTwitter(t *testing.T) {
	doc, _ := makeDoc("@twitter, what time does it start :-)")
	expected := []string{"@twitter", ",", "what", "time", "does", "it", "start", ":-)"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(1)")

	doc, _ = makeDoc("Mr. James plays basketball in the N.B.A., do you?")
	expected = []string{
		"Mr.", "James", "plays", "basketball", "in", "the", "N.B.A.", ",",
		"do", "you", "?"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(2)")

	doc, _ = makeDoc("ˌˌ kill the last letter")
	expected = []string{"ˌˌ", "kill", "the", "last", "letter"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(3)")

	doc, _ = makeDoc("ˌˌˌ kill the last letter")
	expected = []string{"ˌˌˌ", "kill", "the", "last", "letter"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(4)")

	doc, _ = makeDoc("March. July. March. June. January.")
	expected = []string{
		"March", ".", "July", ".", "March", ".", "June", ".", "January", "."}
	checkCase(t, doc, expected, "TokenizationWebParagraph(5)")
}

func TestTokenizationContractions(t *testing.T) {
	tokenizer := NewIterTokenizer()
	tokens := tokenizer.Tokenize("He's happy")
	expected := []string{"He", "'s", "happy"}
	checkTokens(t, tokens, expected, "TokenizationContraction(default-found)")

	tokens = tokenizer.Tokenize("I've been better")
	expected = []string{"I've", "been", "better"}
	checkTokens(t, tokens, expected, "TokenizationContraction(default-missing)")

	tokenizer = NewIterTokenizer(UsingContractions([]string{"'ve"}))
	tokens = tokenizer.Tokenize("I've been better")
	expected = []string{"I", "'ve", "been", "better"}
	checkTokens(t, tokens, expected, "TokenizationContraction(custom-found)")

	tokens = tokenizer.Tokenize("He's happy")
	expected = []string{"He's", "happy"}
	checkTokens(t, tokens, expected, "TokenizationContraction(custom-missing)")
}

func TestTokenizationMultipleContractions(t *testing.T) {
	tokenizer := NewIterTokenizer()
	tokens := tokenizer.Tokenize("He's're")
	for _, token := range tokens {
		t.Logf("token: %v\n", token)
	}
	expected := []string{"He", "'s", "'re"}
	checkTokens(t, tokens, expected, "TokenizationContraction(default-found)")
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
