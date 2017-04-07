/*
Package summarize implements functions for analyzing readability and usage statistics of text.
*/
package summarize

import (
	"unicode"

	"github.com/jdkato/prose/internal/util"
	"github.com/jdkato/prose/tokenize"
)

// A Document represents a collection of text to be analyzed.
//
// A Document's calculations depend on its word and sentence tokenizers. You
// can use the defaults by invoking NewDocument, choose another implemention
// from the tokenize package, or use your own (as long as it implements the
// ProseTokenizer interface). For example,
//
//    d := Document{Content: ..., WordTokenizer: ..., SentenceTokenizer: ...}
//    d.Initialize()
//
// TODO: There should be a way to efficiently add or remove text from a the
// content of a Document (e.g., we should be able to build it incrementally).
// Perhaps we should look into using a rope as our underlying data structure?
type Document struct {
	Content         string           // Actual text
	NumCharacters   float64          // Number of Characters
	NumComplexWords float64          // PolysylWords without common suffixes
	NumPolysylWords float64          // Number of words with > 2 syllables
	NumSentences    float64          // Number of sentences
	NumSyllables    float64          // Number of syllables
	NumWords        float64          // Number of words
	Sentences       map[string]int   // {sentence: length}
	Words           map[string][]int // {word: [frequency, syllables]}

	SentenceTokenizer tokenize.ProseTokenizer
	WordTokenizer     tokenize.ProseTokenizer
}

// An Assessment provides comprehensive access to a Document's metrics.
type Assessment struct {
	AutomatedReadability float64
	FleschKincaid        float64
	ReadingEase          float64
	GunningFog           float64
	SMOG                 float64
}

// NewDocument is a Document constructor that takes a string as an argument. It
// then calculates the data necessary for computing readability and usage
// statistics.
//
// This is a convenience wrapper around the Document initialization process
// that defaults to using a WordBoundaryTokenizer and a PunktSentenceTokenizer
// as its word and sentence tokenizers, respectively.
func NewDocument(text string) *Document {
	wTok := tokenize.NewWordBoundaryTokenizer()
	sTok := tokenize.NewPunktSentenceTokenizer()
	doc := Document{Content: text, WordTokenizer: wTok, SentenceTokenizer: sTok}
	doc.Initialize()
	return &doc
}

// Initialize calculates the data necessary for computing readability and usage
// statistics.
func (d *Document) Initialize() {
	d.Words = make(map[string][]int)
	d.Sentences = make(map[string]int)
	for _, s := range d.SentenceTokenizer.Tokenize(d.Content) {
		wordCount := d.NumWords
		d.NumSentences++
		for _, word := range d.WordTokenizer.Tokenize(s) {
			d.NumCharacters += countChars(word)
			syllables := Syllables(word)
			if _, found := d.Words[word]; found {
				d.Words[word][0]++
			} else {
				d.Words[word] = []int{1, syllables}
			}
			d.NumSyllables += float64(syllables)
			if syllables > 2 {
				d.NumPolysylWords++
			}
			if isComplex(word, syllables) {
				d.NumComplexWords++
			}
			d.NumWords++
		}
		d.Sentences[s] = int(d.NumWords - wordCount)
	}
}

// Assess returns an Assessment for the Document d.
func (d *Document) Assess() *Assessment {
	return &Assessment{
		FleschKincaid: d.FleschKincaid(), ReadingEase: d.ReadingEase(),
		GunningFog: d.Gunningfog(), SMOG: d.SMOG(),
		AutomatedReadability: d.AutomatedReadability()}
}

// Syllables returns the number of syllables in the string word.
func Syllables(word string) int {
	vowels := []rune{'a', 'e', 'i', 'o', 'u', 'y'}
	vowelCount := 0
	ext := len(word)

	lastWasVowel := false
	for _, c := range word {
		found := false
		for _, v := range vowels {
			if v == c {
				found = true
				if !lastWasVowel {
					vowelCount++
				}
				break
			}
		}
		lastWasVowel = found
	}
	if (ext > 2 && word[ext-2:] == "es") || (ext > 1 && word[ext-1:] == "e") {
		vowelCount--
	}

	return vowelCount
}

func isComplex(word string, syllables int) bool {
	if util.HasAnySuffix(word, []string{"es", "ed", "ing"}) {
		syllables--
	}
	return syllables > 2
}

func countChars(word string) float64 {
	count := 0
	for _, c := range word {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			count++
		}
	}
	return float64(count)
}
