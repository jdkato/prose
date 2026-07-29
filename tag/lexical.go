package tag

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Lexical assigns tags from a dictionary, asking a fallback tagger only about
// words the dictionary does not know.
//
// This is how a grammar checker usually tags: `aware` is an adjective in every
// English sentence, and a dictionary says so without reading the words around
// it. A contextual tagger has to infer it, and infers badly when the sentence
// is ungrammatical -- which is exactly the sentence a grammar rule is looking
// at. In "they were all ready aware of the risk" the perceptron reads `aware`
// as a noun, because the words beside it are not evidence for anything.
//
// The cost is the other direction: a dictionary cannot tell `code` the noun
// from `code` the verb. Words the dictionary lists more than one way are left
// to the fallback, so context still decides where context is what decides.
type Lexical struct {
	// entries maps a word to the tags it can take. One tag means the
	// dictionary is sure; several mean the fallback picks between them.
	entries map[string][]string

	fallback Interface
	name     string
}

// NewLexical builds a dictionary-driven tagger.
//
// fallback handles unknown and ambiguous words; passing nil uses the
// perceptron.
func NewLexical(entries map[string][]string, fallback Interface) (*Lexical, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("tag: a lexical tagger needs a dictionary")
	}

	if fallback == nil {
		p, err := New()
		if err != nil {
			return nil, err
		}
		fallback = p
	}

	return &Lexical{entries: entries, fallback: fallback, name: "lexical"}, nil
}

// Name reports which tagger this is.
func (l *Lexical) Name() string { return l.name }

// Tag assigns a tag to each word and returns the tagged tokens.
func (l *Lexical) Tag(words []string) []Token {
	out := make([]Token, len(words))
	for i, w := range words {
		out[i].Text = w
	}
	l.TagTokens(out)

	return out
}

// TagTokens fills in the Tag field of each token.
//
// The fallback runs over the whole sentence rather than word by word: it reads
// the tags around each word, so handing it only the words the dictionary could
// not place would have it reading a sentence with holes in it.
func (l *Lexical) TagTokens(tokens []Token) {
	if len(tokens) == 0 {
		return
	}

	l.fallback.TagTokens(tokens)

	for i := range tokens {
		tags := l.lookup(tokens[i].Text)
		if len(tags) == 0 {
			continue
		}

		// One possibility: the dictionary decides.
		if len(tags) == 1 {
			tokens[i].Tag = tags[0]
			continue
		}

		// Several: keep the fallback's answer when it is among them, since
		// context is what tells them apart. Otherwise fall back to the
		// dictionary's first, which is the most common reading.
		if !contains(tags, tokens[i].Tag) {
			tokens[i].Tag = tags[0]
		}
	}
}

// lookup finds a word, trying it as written before its lower-case form so one
// entry covers every casing.
func (l *Lexical) lookup(word string) []string {
	if tags, ok := l.entries[word]; ok {
		return tags
	}

	lower := strings.ToLower(word)
	if lower == word {
		return nil
	}

	return l.entries[lower]
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}

// ReadDictionary parses a dictionary of one word per line:
//
//	aware	JJ
//	code	NN	VB
//
// Fields are tab-separated: the word, then every tag it can take, most common
// first. Blank lines and lines beginning with # are ignored.
func ReadDictionary(r io.Reader) (map[string][]string, error) {
	entries := map[string][]string{}

	sc := bufio.NewScanner(r)
	for line := 1; sc.Scan(); line++ {
		text := strings.TrimSpace(sc.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}

		fields := strings.Split(text, "\t")
		if len(fields) < 2 {
			return nil, fmt.Errorf("tag: line %d: want a word and at least one tag", line)
		}

		word := strings.TrimSpace(fields[0])
		tags := make([]string, 0, len(fields)-1)
		for _, f := range fields[1:] {
			if f = strings.TrimSpace(f); f != "" {
				tags = append(tags, f)
			}
		}

		if word == "" || len(tags) == 0 {
			return nil, fmt.Errorf("tag: line %d: want a word and at least one tag", line)
		}
		entries[word] = tags
	}

	return entries, sc.Err()
}
