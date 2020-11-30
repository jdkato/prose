package prose

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenTester func(string) bool

type Tokenizer interface {
	Tokenize(string) []*Token
}

// iterTokenizer splits a sentence into words.
type iterTokenizer struct {
	specialRE      *regexp.Regexp
	sanitizer      *strings.Replacer
	contractions   []string
	suffixes       []string
	prefixes       []string
	emoticons      map[string]int
	isUnsplittable TokenTester
}

type TokenizerOptFunc func(*iterTokenizer)

// UsingIsUnsplittableFN gives a function that tests whether a token is splittable or not.
func UsingIsUnsplittable(x TokenTester) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.isUnsplittable = x
	}
}

// Use the provided special regex for unsplittable tokens.
func UsingSpecialRE(x *regexp.Regexp) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.specialRE = x
	}
}

// Use the provided sanitizer.
func UsingSanitizer(x *strings.Replacer) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.sanitizer = x
	}
}

// Use the provided suffixes.
func UsingSuffixes(x []string) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.suffixes = x
	}
}

// Use the provided prefixes.
func UsingPrefixes(x []string) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.prefixes = x
	}
}

// Use the provided map of emoticons.
func UsingEmoticons(x map[string]int) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.emoticons = x
	}
}

// Use the provided contractions.
func UsingContractions(x []string) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.contractions = x
	}
}

// Constructor for default iterTokenizer
func NewIterTokenizer(opts ...TokenizerOptFunc) *iterTokenizer {
	tok := new(iterTokenizer)

	// Set default parameters
	tok.contractions = contractions
	tok.emoticons = emoticons
	tok.isUnsplittable = func(_ string) bool { return false }
	tok.prefixes = prefixes
	tok.sanitizer = sanitizer
	tok.specialRE = internalRE
	tok.suffixes = suffixes

	// Apply options if provided
	for _, applyOpt := range opts {
		applyOpt(tok)
	}

	return tok
}

func addToken(toks []*Token, s string, from, to int) []*Token {
	if strings.TrimSpace(s) != "" {
		toks = append(toks, &Token{Text: s, Start: from, End: to})
	}
	return toks
}

func (t *iterTokenizer) isSpecial(token string) bool {
	_, found := t.emoticons[token]
	return found || t.specialRE.MatchString(token) || t.isUnsplittable(token)
}

func (t *iterTokenizer) doSplit(token string, offset int) []*Token {
	tokens := []*Token{}
	suffs := []*Token{}

	last := 0
	for token != "" && utf8.RuneCountInString(token) != last {
		if t.isSpecial(token) {
			// We've found a special case (e.g., an emoticon) -- so, we add it as a token without
			// any further processing.
			tokens = addToken(tokens, token, offset, offset+len(token))
			break
		}
		last = utf8.RuneCountInString(token)
		lower := strings.ToLower(token)
		if length := hasAnyPrefix(token, t.prefixes); length > 0 {
			// Remove prefixes -- e.g., $100 -> [$, 100].
			tokens = addToken(tokens, token[:length], offset, offset+length)
			token = token[length:]
			offset += length
		} else if idx := hasAnyIndex(lower, t.contractions); idx > -1 {
			// Handle "they'll", "I'll", "Don't", "won't", etc.
			//
			// they'll -> [they, 'll].
			// don't -> [do, n't].
			tokens = addToken(tokens, token[:idx], offset, offset+idx)
			token = token[idx:]
			offset += idx
		} else if length := hasAnySuffix(token, t.suffixes); length > 0 {
			// Remove suffixes -- e.g., Well) -> [Well, )].
			suffs = append([]*Token{
				{Text: string(token[len(token)-length]),
					Start: offset + len(token) - length,
					End:   offset + len(token)}},
				suffs...)
			token = token[:len(token)-1]
		} else {
			tokens = addToken(tokens, token, offset, offset+len(token))
		}
	}
	return append(tokens, suffs...)
}

type tokensOffset struct {
	tl     []*Token
	offset int
}

// tokenize splits a sentence into a slice of words.
func (t *iterTokenizer) Tokenize(text string) []*Token {
	tokens := []*Token{}

	clean, white := t.sanitizer.Replace(text), false
	length := len(clean)

	start, index := 0, 0
	cache := map[string]tokensOffset{}
	for index <= length {
		uc, size := utf8.DecodeRuneInString(clean[index:])
		if size == 0 {
			break
		} else if index == 0 {
			white = unicode.IsSpace(uc)
		}
		if unicode.IsSpace(uc) != white {
			if start < index {
				span := clean[start:index]
				if toks, found := cache[span]; found {
					for _, t := range toks.tl {
						tokens = append(tokens, &Token{
							Tag:   t.Tag,
							Text:  t.Text,
							Label: t.Label,
							Start: t.Start - toks.offset + start,
							End:   t.End - toks.offset + start,
						})
					}
				} else {
					toks := t.doSplit(span, start)
					cache[span] = tokensOffset{toks, start}
					tokens = append(tokens, toks...)
				}
			}
			if uc == ' ' {
				start = index + 1
			} else {
				start = index
			}
			white = !white
		}
		index += size
	}

	if start < index {
		tokens = append(tokens, t.doSplit(clean[start:index], start)...)
	}

	return tokens
}

var internalRE = regexp.MustCompile(`^(?:[A-Za-z]\.){2,}$|^[A-Z][a-z]{1,2}\.$`)
var sanitizer = strings.NewReplacer(
	"\u201c", `"`,
	"\u201d", `"`,
	"\u2018", "'",
	"\u2019", "'",
	"&rsquo;", "'")
var contractions = []string{"'ll", "'s", "'re", "'m", "n't"}
var suffixes = []string{",", ")", `"`, "]", "!", ";", ".", "?", ":", "'"}
var prefixes = []string{"$", "(", `"`, "["}
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
	"-__-":        1,
	"8-)":         1,
	"8-D":         1,
	"8D":          1,
	":(":          1,
	":((":         1,
	":(((":        1,
	":()":         1,
	":)))":        1,
	":-)":         1,
	":-))":        1,
	":-)))":       1,
	":-*":         1,
	":-/":         1,
	":-X":         1,
	":-]":         1,
	":-o":         1,
	":-p":         1,
	":-x":         1,
	":-|":         1,
	":-}":         1,
	":0":          1,
	":3":          1,
	":P":          1,
	":]":          1,
	":`(":         1,
	":`)":         1,
	":`-(":        1,
	":o":          1,
	":o)":         1,
	"=(":          1,
	"=)":          1,
	"=D":          1,
	"=|":          1,
	"@_@":         1,
	"O.o":         1,
	"O_o":         1,
	"V_V":         1,
	"XDD":         1,
	"[-:":         1,
	"^___^":       1,
	"o_0":         1,
	"o_O":         1,
	"o_o":         1,
	"v_v":         1,
	"xD":          1,
	"xDD":         1,
	"¯\\(ツ)/¯":    1,
}
