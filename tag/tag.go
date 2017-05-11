/*
Package tag implements functions for tagging parts of speech.
*/
package tag

import "strings"

// Token represents a tagged section of text.
type Token struct {
	Text string
	Tag  string
}

// Tuple is a list of tuples in the form (words, tags).
type Tuple [][][]string

// Len returns the length of a Tuple.
func (t Tuple) Len() int { return len(t) }

// Swap switches the ith and jth elements in a Tuple.
func (t Tuple) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// ReadTagged converts pre-tagged input into a Tuple suitable for training.
func ReadTagged(text, sep string) Tuple {
	t := Tuple{}
	for _, sent := range strings.Split(text, "\n") {
		tokens := []string{}
		tags := []string{}
		for _, token := range strings.Split(sent, " ") {
			parts := strings.Split(token, sep)
			tokens = append(tokens, parts[0])
			tags = append(tags, parts[1])
		}
		t = append(t, [][]string{tokens, tags})
	}
	return t
}
