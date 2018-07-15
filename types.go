package prose

// A Token represents an individual token of text such as a word or punctuation
// symbol.
//
// The tokenization process performs the following steps: (1) split on
// contractions (e.g., "don't" -> [do n't]), (2) split on non-terminating
// punctuation, (3) split on single quotes when followed by whitespace, and (4)
// split on periods that appear at the end of lines.
type Token struct {
	Tag   string // The token's part-of-speech tag.
	Text  string // The token's actual content.
	Label string // The token's IOB label.
}

// An Entity represents an individual named-entity.
type Entity struct {
	Text  string // The entity's actual content.
	Label string // The entity's label.
}

// A Sentence represents a segmented portion of text.
type Sentence struct {
	Text string // The sentence's text.
}

// A span represents an in-text location within a Document.
type span struct {
	attributes map[string]string
	begin      int
	end        int
}
