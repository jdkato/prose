package summarize

import "math"

// FleschKincaid computes the Flesch–Kincaid grade level of the Document d.
func (d *Document) FleschKincaid() float64 {
	x := 0.39 * d.NumWords / d.NumSentences
	y := 11.8 * d.NumSyllables / d.NumWords
	return x + y - 15.59
}

// ReadingEase computes the Flesch reading-ease score of the Document d.
func (d *Document) ReadingEase() float64 {
	x := 1.015 * d.NumWords / d.NumSentences
	y := 84.6 * d.NumSyllables / d.NumWords
	return 206.835 - x - y
}

// Gunningfog computes the Gunning fog index score of the Document d.
func (d *Document) Gunningfog() float64 {
	x := d.NumWords / d.NumSentences
	y := d.NumComplexWords / d.NumWords
	return 0.4 * (x + 100.0*y)
}

// SMOG computes the SMOG grade of the Document d.
func (d *Document) SMOG() float64 {
	return 1.0430*math.Sqrt(d.NumPolysylWords*30.0/d.NumSentences) + 3.1291
}

// AutomatedReadability computes the automated readability index score of the
// Document d.
func (d *Document) AutomatedReadability() float64 {
	x := 4.71 * (d.NumCharacters / d.NumWords)
	y := 0.5 * (d.NumWords / d.NumSentences)
	return x + y - 21.43
}

// ColemanLiau computes the Coleman–Liau index score of the Document d.
func (d *Document) ColemanLiau() float64 {
	x := 0.0588 * (d.NumCharacters / d.NumWords) * 100
	y := 0.296 * (d.NumSentences / d.NumWords) * 100
	return x - y - 15.8
}
