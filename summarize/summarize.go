/*
Package summarize implements functions for analyzing readability and usage statistics of text.
*/
package summarize

import (
	"unicode"

	"github.com/jdkato/prose/internal/util"
	"github.com/jdkato/prose/tokenize"
	"gopkg.in/neurosnap/sentences.v1/english"
)

// A Document represents a collection of text to be analyzed.
type Document struct {
	NumCharacters   float64
	NumComplexWords float64
	NumPolysylWords float64
	NumSentences    float64
	Sentences       []string
	NumSyllables    float64
	Content         string
	NumWords        float64
	Words           []string
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
func NewDocument(text string) *Document {
	doc := Document{Content: text}

	wordTokenizer := tokenize.NewWordBoundaryTokenizer()
	sentTokenizer, err := english.NewSentenceTokenizer(nil)
	util.CheckError(err)

	for _, s := range sentTokenizer.Tokenize(text) {
		doc.Sentences = append(doc.Sentences, s.Text)
		doc.NumSentences++
		for _, word := range wordTokenizer.Tokenize(s.Text) {
			doc.NumCharacters += countChars(word)
			syllables := Syllables(word)
			doc.NumSyllables += float64(syllables)
			if syllables > 2 {
				doc.NumPolysylWords++
			}
			if isComplex(word, syllables) {
				doc.NumComplexWords++
			}
			doc.Words = append(doc.Words, word)
			doc.NumWords++
		}
	}
	return &doc
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
