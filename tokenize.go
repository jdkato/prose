package prose

import (
	"regexp"
	"strings"
)

// iterTokenizer splits a sentence into words.
type iterTokenizer struct {
}

// newIterTokenizer is a iterTokenizer constructor.
func newIterTokenizer() *iterTokenizer {
	return new(iterTokenizer)
}

func splitWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t' || r == '\r'
}

func addToken(s string, toks []Token) []Token {
	if s != "" {
		toks = append(toks, Token{Text: s})
	}
	return toks
}

// tokenize splits a sentence into a slice of words.
func (t *iterTokenizer) tokenize(text string) []Token {
	tokens := []Token{}

	clean := sanitizer.Replace(text)
	for _, token := range strings.FieldsFunc(clean, splitWhitespace) {
		lower := strings.ToLower(token)
		if _, found := emoticons[token]; found {
			// We've found an emoticon -- so, we add it as a token without any
			// further processing.
			tokens = append(tokens, Token{Text: token})
		} else {
			for hasAnyPrefix(token, prefixes) {
				// Remove prefixes -- e.g., $100 -> [$, 100].
				tokens = addToken(string(token[0]), tokens)
				token = token[1:]
			}

			if idx := hasAnyIndex(lower, []string{"'ll", "'s", "'re", "'m"}); idx > -1 {
				// Handle "they'll", "I'll", etc.
				//
				// they'll -> [they, 'll].
				tokens = addToken(token[:idx], tokens)
				token = token[idx:]
			} else if idx := hasAnyIndex(lower, []string{"n't"}); idx > -1 {
				// Handle "Don't", "won't", etc.
				//
				// don't -> [do, n't].
				tokens = addToken(token[:idx], tokens)
				token = token[idx:]
			}

			if internalRE.MatchString(token) {
				// We've found an instance of non-terminating punctuation --
				// e.g., "N.B.A." or "Mr."
				tokens = addToken(token, tokens)
			} else {
				tokens = addToken(strings.TrimRight(token, suffcset), tokens)
				found := []string{}
				for hasAnySuffix(token, suffixes) {
					// Remove suffixes -- e.g., Well) -> [Well, )].
					found = append([]string{string(token[len(token)-1])}, found...)
					token = token[:len(token)-1]
				}
				for _, token := range found {
					tokens = addToken(token, tokens)
				}

			}
		}
	}

	return tokens
}

var internalRE = regexp.MustCompile(`(?:[A-Za-z]\.){2,}$|^[A-Z][a-z]{1,3}\.$`)
var sanitizer = strings.NewReplacer(
	"\u201c", `"`,
	"\u201d", `"`,
	"\u2018", "'",
	"\u2019", "'",
	"&rsquo;", "'")
var suffixes = []string{",", ")", `"`, "]", "!", ";", ".", "?", ":", "'"}
var prefixes = []string{"$", "(", `"`, "["}
var suffcset = strings.Join(suffixes, "")
var emoticons = map[string]int{
	"(-8":         1,
	"(-;":         1,
	"(-_-)":       1,
	"(._.)":       1,
	"(:":          1,
	"(=":          1,
	"(o:":         1,
	"(¬_¬)":       1,
	"(ಠ_ಠ)":       1,
	"(╯°□°）╯︵┻━┻": 1,
	"-__-":     1,
	"8-)":      1,
	"8-D":      1,
	"8D":       1,
	":(":       1,
	":((":      1,
	":(((":     1,
	":()":      1,
	":)))":     1,
	":-)":      1,
	":-))":     1,
	":-)))":    1,
	":-*":      1,
	":-/":      1,
	":-X":      1,
	":-]":      1,
	":-o":      1,
	":-p":      1,
	":-x":      1,
	":-|":      1,
	":-}":      1,
	":0":       1,
	":3":       1,
	":P":       1,
	":]":       1,
	":`(":      1,
	":`)":      1,
	":`-(":     1,
	":o":       1,
	":o)":      1,
	"=(":       1,
	"=)":       1,
	"=D":       1,
	"=|":       1,
	"@_@":      1,
	"O.o":      1,
	"O_o":      1,
	"V_V":      1,
	"XDD":      1,
	"[-:":      1,
	"^___^":    1,
	"o_0":      1,
	"o_O":      1,
	"o_o":      1,
	"v_v":      1,
	"xD":       1,
	"xDD":      1,
	"¯\\(ツ)/¯": 1,
}
