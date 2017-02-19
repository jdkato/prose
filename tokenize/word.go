package tokenize

import (
	"regexp"
	"strings"
)

var startingQuotes = map[string]*regexp.Regexp{
	"``":     regexp.MustCompile(`^\"`),
	" $1 ":   regexp.MustCompile("(``)"),
	"$1 `` ": regexp.MustCompile(`'([ (\[{<])"`),
}
var punctuation = map[string]*regexp.Regexp{
	" $1$2":    regexp.MustCompile(`([:,])([^\d])`),
	" ... ":    regexp.MustCompile(`\.\.\.`),
	"$1 $2$3 ": regexp.MustCompile(`([^\.])(\.)([\]\)}>"\']*)\s*$`),
	"$1 ' ":    regexp.MustCompile(`([^'])' `),
}
var punctuation2 = []*regexp.Regexp{
	regexp.MustCompile(`([:,])$`),
	regexp.MustCompile(`([;@#$%&?!])`),
}
var brackets = map[string]*regexp.Regexp{
	" $1 ": regexp.MustCompile(`[\]\[\(\)\{\}\<\>]`),
	" -- ": regexp.MustCompile(`--`),
}
var endingQuotes = map[string]*regexp.Regexp{
	" '' ": regexp.MustCompile(`"`),
}
var endingQuotes2 = []*regexp.Regexp{
	regexp.MustCompile(`'(\S)(\'\')'`),
	regexp.MustCompile(`([^' ])('[sS]|'[mM]|'[dD]|') `),
	regexp.MustCompile(`([^' ])('ll|'LL|'re|'RE|'ve|'VE|n't|N'T) `),
}
var contractions = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(can)(not)\b`),
	regexp.MustCompile(`(?i)\b(d)('ye)\b`),
	regexp.MustCompile(`(?i)\b(gim)(me)\b`),
	regexp.MustCompile(`(?i)\b(gon)(na)\b`),
	regexp.MustCompile(`(?i)\b(got)(ta)\b`),
	regexp.MustCompile(`(?i)\b(lem)(me)\b`),
	regexp.MustCompile(`(?i)\b(mor)('n)\b`),
	regexp.MustCompile(`(?i)\b(wan)(na) `),
	regexp.MustCompile(`(?i) ('t)(is)\b`),
	regexp.MustCompile(`(?i) ('t)(was)\b`),
}

// WordTokenizer ...
type WordTokenizer func(sent string) []string

// WordTokenizerFn ...
func WordTokenizerFn(sent string) []string {
	for substitution, r := range startingQuotes {
		sent = r.ReplaceAllString(sent, substitution)
	}

	for substitution, r := range punctuation {
		sent = r.ReplaceAllString(sent, substitution)
	}

	for _, r := range punctuation2 {
		sent = r.ReplaceAllString(sent, " $1 ")
	}

	for substitution, r := range brackets {
		sent = r.ReplaceAllString(sent, substitution)
	}

	sent = " " + sent + " "

	for substitution, r := range endingQuotes {
		sent = r.ReplaceAllString(sent, substitution)
	}

	for _, r := range endingQuotes2 {
		sent = r.ReplaceAllString(sent, "$1 $2 ")
	}

	for _, r := range contractions {
		sent = r.ReplaceAllString(sent, " $1 $2 ")
	}

	return strings.Split(strings.TrimSpace(sent), " ")
}
