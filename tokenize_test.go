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
		t.Errorf("%v: unexpected tokens: %#v", name, observed)
	}
}

func checkCase(t *testing.T, doc *Document, expected []string, name string) {
	tokens := getTokenText(doc)
	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("%v: unexpected tokens", name)
	}
}

func checkStartEnd(t *testing.T, token *Token, expectedText string, expectedStart, expectedEnd int) {
	if token.Text != expectedText {
		t.Errorf("got %v, want %v", token.Text, expectedText)
	}
	if token.Start != expectedStart {
		t.Errorf("got %v, want %v", token.Start, expectedStart)
	}
	if token.End != expectedEnd {
		t.Errorf("got %v, want %v", token.End, expectedEnd)
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
	text := "@twitter, what time does it start :-)"
	doc, _ := makeDoc(text)
	expected := []string{"@twitter", ",", "what", "time", "does", "it", "start", ":-)"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(1)")
	checkStartEnd(t, doc.tokens[0], "@twitter", 0, 8)
	checkStartEnd(t, doc.tokens[1], ",", 8, 9)
	checkStartEnd(t, doc.tokens[2], "what", 10, 14)
	checkStartEnd(t, doc.tokens[3], "time", 15, 19)
	checkStartEnd(t, doc.tokens[4], "does", 20, 24)
	checkStartEnd(t, doc.tokens[5], "it", 25, 27)
	checkStartEnd(t, doc.tokens[6], "start", 28, 33)
	checkStartEnd(t, doc.tokens[7], ":-)", 34, len(text))

	text = "Mr. James plays basketball in the N.B.A., do you?"
	doc, _ = makeDoc(text)
	expected = []string{
		"Mr.", "James", "plays", "basketball", "in", "the", "N.B.A.", ",",
		"do", "you", "?"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(2)")
	checkStartEnd(t, doc.tokens[0], "Mr.", 0, 3)
	checkStartEnd(t, doc.tokens[1], "James", 4, 9)
	checkStartEnd(t, doc.tokens[2], "plays", 10, 15)
	checkStartEnd(t, doc.tokens[3], "basketball", 16, 26)
	checkStartEnd(t, doc.tokens[4], "in", 27, 29)
	checkStartEnd(t, doc.tokens[5], "the", 30, 33)
	checkStartEnd(t, doc.tokens[6], "N.B.A.", 34, 40)
	checkStartEnd(t, doc.tokens[7], ",", 40, 41)
	checkStartEnd(t, doc.tokens[8], "do", 42, 44)
	checkStartEnd(t, doc.tokens[9], "you", 45, 48)
	checkStartEnd(t, doc.tokens[10], "?", 48, len(text))

	text = "ˌˌ kill the last letter"
	doc, _ = makeDoc(text)
	expected = []string{"ˌˌ", "kill", "the", "last", "letter"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(3)")
	checkStartEnd(t, doc.tokens[0], "ˌˌ", 0, 4)
	checkStartEnd(t, doc.tokens[1], "kill", 5, 9)
	checkStartEnd(t, doc.tokens[2], "the", 10, 13)
	checkStartEnd(t, doc.tokens[3], "last", 14, 18)
	checkStartEnd(t, doc.tokens[4], "letter", 19, len(text))

	text = "ˌˌˌ kill the last letter"
	doc, _ = makeDoc(text)
	expected = []string{"ˌˌˌ", "kill", "the", "last", "letter"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(4)")
	checkStartEnd(t, doc.tokens[0], "ˌˌˌ", 0, 6)
	checkStartEnd(t, doc.tokens[1], "kill", 7, 11)
	checkStartEnd(t, doc.tokens[2], "the", 12, 15)
	checkStartEnd(t, doc.tokens[3], "last", 16, 20)
	checkStartEnd(t, doc.tokens[4], "letter", 21, len(text))

	text = "March. July. March. June.  January."
	doc, _ = makeDoc(text)
	expected = []string{
		"March", ".", "July", ".", "March", ".", "June", ".", "January", "."}
	checkCase(t, doc, expected, "TokenizationWebParagraph(5)")
	checkStartEnd(t, doc.tokens[0], "March", 0, 5)
	checkStartEnd(t, doc.tokens[1], ".", 5, 6)
	checkStartEnd(t, doc.tokens[2], "July", 7, 11)
	checkStartEnd(t, doc.tokens[3], ".", 11, 12)
	checkStartEnd(t, doc.tokens[4], "March", 13, 18)
	checkStartEnd(t, doc.tokens[5], ".", 18, 19)
	checkStartEnd(t, doc.tokens[6], "June", 20, 24)
	checkStartEnd(t, doc.tokens[7], ".", 24, 25)
	checkStartEnd(t, doc.tokens[8], "January", 27, 34)
	checkStartEnd(t, doc.tokens[9], ".", 34, len(text))
}

func TestTokenizationContractions(t *testing.T) {
	tokenizer := NewIterTokenizer()
	tokens := tokenizer.Tokenize("He's happy")
	expected := []string{"He", "'s", "happy"}
	checkTokens(t, tokens, expected, "TokenizationContraction(default-found)")
	checkStartEnd(t, tokens[0], "He", 0, 2)
	checkStartEnd(t, tokens[1], "'s", 2, 4)
	checkStartEnd(t, tokens[2], "happy", 5, 10)

	tokens = tokenizer.Tokenize("I've been better")
	expected = []string{"I've", "been", "better"}
	checkTokens(t, tokens, expected, "TokenizationContraction(default-missing)")
	checkStartEnd(t, tokens[0], "I've", 0, 4)
	checkStartEnd(t, tokens[1], "been", 5, 9)
	checkStartEnd(t, tokens[2], "better", 10, 16)

	tokenizer = NewIterTokenizer(UsingContractions([]string{"'ve"}))
	tokens = tokenizer.Tokenize("I've been better")
	expected = []string{"I", "'ve", "been", "better"}
	checkTokens(t, tokens, expected, "TokenizationContraction(custom-found)")
	checkStartEnd(t, tokens[0], "I", 0, 1)
	checkStartEnd(t, tokens[1], "'ve", 1, 4)
	checkStartEnd(t, tokens[2], "been", 5, 9)
	checkStartEnd(t, tokens[3], "better", 10, 16)

	tokens = tokenizer.Tokenize("He's happy")
	expected = []string{"He's", "happy"}
	checkTokens(t, tokens, expected, "TokenizationContraction(custom-missing)")
	checkStartEnd(t, tokens[0], "He's", 0, 4)
	checkStartEnd(t, tokens[1], "happy", 5, 10)
}

func TestTokenizationSuffixes(t *testing.T) {
	tokenizer := NewIterTokenizer()
	tokens := tokenizer.Tokenize("(Well,\nthat\twasn't good).")
	expected := []string{"(", "Well", ",", "that", "was", "n't", "good", ")", "."}
	checkTokens(t, tokens, expected, "TestTokenizationSuffixes")
	checkStartEnd(t, tokens[0], "(", 0, 1)
	checkStartEnd(t, tokens[1], "Well", 1, 5)
	checkStartEnd(t, tokens[2], ",", 5, 6)
	checkStartEnd(t, tokens[3], "that", 7, 11)
	checkStartEnd(t, tokens[4], "was", 12, 15)
	checkStartEnd(t, tokens[5], "n't", 15, 18)
	checkStartEnd(t, tokens[6], "good", 19, 23)
	checkStartEnd(t, tokens[7], ")", 23, 24)
	checkStartEnd(t, tokens[8], ".", 24, 25)
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
