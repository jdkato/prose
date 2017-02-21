package tokenize

// WordTokenizer is the interface implemented by an object that takes a string
// and returns a slice of strings representing words.
//
// Implementations include:
// * TreebankWordTokenizer
type WordTokenizer interface {
	Tokenize(text string) []string
}

// SentenceTokenizer is the interface implemented by an object that takes a
// string and returns a slice of representing sentences.
//
// Implementations include:
// * PunktSentenceTokenizer
type SentenceTokenizer interface {
	Tokenize(text string) []string
}
