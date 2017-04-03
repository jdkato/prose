package tokenize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var treebankTests = []struct {
	in  string
	out []string
}{
	{"They'll save and invest more.", []string{"They", "'ll", "save", "and", "invest", "more", "."}},
	{"How's it going?", []string{"How", "'s", "it", "going", "?"}},
	{"abbreviations like M.D. and initials containing periods, they", []string{"abbreviations", "like", "M.D.", "and", "initials", "containing", "periods", ",", "they"}},
}

func TestTreebankWordTokenizer(t *testing.T) {
	tokenizer := NewTreebankWordTokenizer()
	for _, tt := range treebankTests {
		assert.Equal(t, tt.out, tokenizer.Tokenize(tt.in))
	}
}
