package prose

import (
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/mingrammer/commonregex"
	"github.com/willf/pad"
)

// TreebankWordTokenizer splits a sentence into words.
//
// This implementation is a port of the Sed script written by Robert McIntyre,
// which is available at https://gist.github.com/jdkato/fc8b8c4266dba22d45ac85042ae53b1e.
type treebankWordTokenizer struct {
}

// newTreebankWordTokenizer is a TreebankWordTokenizer constructor.
func newTreebankWordTokenizer() *treebankWordTokenizer {
	return new(treebankWordTokenizer)
}

var sanitizer = strings.NewReplacer(
	"\u201c", `"`,
	"\u201d", `"`,
	"\u2018", "'",
	"\u2019", "'",
	"&rsquo;", "'",
	"\r\n", "\n",
	"\r", "\n")

var misTok = regexp.MustCompile(`\b[a-z]{3,}\.$`)

var startingQuotes = map[string]*regexp.Regexp{
	"$1 `` ": regexp.MustCompile(`'([ (\[{<])"`),
	"``":     regexp.MustCompile(`^(")`),
	" ``":    regexp.MustCompile(`( ")`),
}
var startingQuotes2 = map[string]*regexp.Regexp{
	" $1 ": regexp.MustCompile("(``)"),
}
var punctuation = map[string]*regexp.Regexp{
	// NOTE: `-` was added ad-hoc -- see tag.go:183
	" $1 $2":   regexp.MustCompile(`([=:,-])([^\d])`),
	" ... ":    regexp.MustCompile(`\.\.\.`),
	"$1 $2$3 ": regexp.MustCompile(`([^\.])(\.)([\]\)}>"\']*)\s*$`),
	"$1 ' ":    regexp.MustCompile(`([^'])' `),
}
var punctuation2 = []*regexp.Regexp{
	regexp.MustCompile(`([:,])$`),
	regexp.MustCompile(`([;#$%&?!])`),
}
var brackets = map[string]*regexp.Regexp{
	" $1 ": regexp.MustCompile(`([\]\[\(\)\{\}\<\>])`),
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
var newlines = regexp.MustCompile(`(?:\n|\n\r|\r)`)
var spaces = regexp.MustCompile(`(?: {2,})`)

func prep(text string) (string, *strings.Replacer) {
	replacements := []string{}
	for i, link := range commonregex.Links(text) {
		key := pad.Right("URI"+strconv.Itoa(i), len(link), "X")
		if !stringInSlice(link, replacements) {
			text = strings.Replace(text, link, key, -1)
		}
		replacements = append(replacements, []string{key, link}...)
	}
	for i, emote := range emoticonRE.FindAllString(text, -1) {
		key := pad.Right("E"+strconv.Itoa(i), len(emote), "X")
		if !stringInSlice(emote, replacements) {
			text = strings.Replace(text, emote, key, -1)
		}
		replacements = append(replacements, []string{key, emote}...)
	}
	clean := html.UnescapeString(sanitizer.Replace(text))
	return clean, strings.NewReplacer(replacements...)
}

// tokenize splits a sentence into a slice of words.
//
// This tokenizer performs the following steps: (1) split on contractions (e.g.,
// "don't" -> [do n't]), (2) split on non-terminating punctuation, (3) split on
// single quotes when followed by whitespace, and (4) split on periods that
// appear at the end of lines.
func (t treebankWordTokenizer) tokenize(text string) []Token {
	clean, replacements := prep(text)

	for substitution, r := range startingQuotes {
		clean = r.ReplaceAllString(clean, substitution)
	}

	for substitution, r := range startingQuotes2 {
		clean = r.ReplaceAllString(clean, substitution)
	}

	for substitution, r := range punctuation {
		clean = r.ReplaceAllString(clean, substitution)
	}

	for _, r := range punctuation2 {
		clean = r.ReplaceAllString(clean, " $1 ")
	}

	for substitution, r := range brackets {
		clean = r.ReplaceAllString(clean, substitution)
	}

	clean = " " + clean + " "

	for substitution, r := range endingQuotes {
		clean = r.ReplaceAllString(clean, substitution)
	}

	for _, r := range endingQuotes2 {
		clean = r.ReplaceAllString(clean, "$1 $2 ")
	}

	for _, r := range contractions {
		clean = r.ReplaceAllString(clean, " $1 $2 ")
	}

	clean = newlines.ReplaceAllString(clean, " ")
	clean = strings.TrimSpace(spaces.ReplaceAllString(clean, " "))
	clean = replacements.Replace(clean)

	tokens := []Token{}
	for _, tok := range strings.SplitAfter(clean, " ") {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		} else if misTok.MatchString(tok) {
			tokens = append(tokens, Token{Text: strings.Trim(tok, ".")})
			tokens = append(tokens, Token{Text: "."})
		} else {
			tokens = append(tokens, Token{Text: strings.TrimSpace(tok)})
		}
	}
	return tokens
}

var emoticons = []string{
	`:>`, `._.`, `[-:`, `:X`, `(-_-)`, `(^_^)`, `:-}`, `ಠ_ಠ`, `¯\\(ツ)/¯`, `;)`,
	`O.o`, `:-(`, `(╯°□°）╯︵┻━┻`, `:*`, `(-8`, `^__^`, `8-D`, `O.O`, `(-;`,
	`:-D`, `=D`, `v.v`, `:o)`, `=/`, `-__-`, `;D`, `@_@`, `8)`, `:o`, `</3`,
	`:-x`, `O_O`, `:-))`, `(-:`, `(._.)`, `V.V`, `:-|`, `0.o`, `:->`, `:p`,
	`8-)`, `:-0`, `xDD`, `>.>`, `:()`, `:1`, `<33`, `)-:`, `:-p`, `0.0`,
	":`-(", `<3`, `><(((*>`, `[:`, `:-*`, `-_-`, `=)`, `;_;`, `:((`, `ಠ︵ಠ`,
	`:P`, `:(`, `>:(`, `o.o`, `xD`, `):`, `(=`, `:}`, `:3`, `;-D`, `(¬_¬)`,
	`:-(((`, `(ಠ_ಠ)`, `:)`, `:0`, `:-((`, `v_v`, `:-)`, `o_o`, `:))`, `(:`,
	`0_0`, `:)))`, `0_o`, `o_O`, `o.O`, `:(((`, `(*_*)`, `O_o`, ":`-)", `:]`,
	`;-)`, `^___^`, `(>_<)`, `(o:`, `:-P`, `:-)))`, `:D`, `o_0`, `<333`, `XDD`,
	`=3`, `:-o`, `:-3`, `=(`, `:O`, ":`)", `o.0`, `:-X`, `:|`, `:-]`, `>:o`,
	`V_V`, `(;`, `8D`, `XD`, `:-/`, `^_^`, ":`(", `:-O`, `<.<`, `:/`, `>.<`,
	`:x`, `=|`}

var emoticonRE = regexp.MustCompile(`:>|\._\.|\[-:|:X|\(-_-\)|\(\^_\^\)|:-\}|ಠ_ಠ|¯\\\\\(ツ\)/¯|;\)|O\.o|:-\(|\(╯°□°）╯︵┻━┻|:\*|\(-8|\^__\^|8-D|O\.O|\(-;|:-D|=D|v\.v|:o\)|=/|-__-|;D|@_@|8\)|:o|</3|:-x|O_O|:-\)\)|\(-:|\(\._\.\)|V\.V|:-\||0\.o|:->|:p|8-\)|:-0|xDD|>\.>|:\(\)|:1|<33|\)-:|:-p|0\.0|<3|><\(\(\(\*>|\[:|:-\*|-_-|=\)|;_;|:\(\(|ಠ︵ಠ|:P|:\(|>:\(|o\.o|xD|\):|\(=|:\}|:3|;-D|\(¬_¬\)|:-\(\(\(|\(ಠ_ಠ\)|:\)|:0|:-\(\(|v_v|:-\)|o_o|:\)\)|\(:|0_0|:\)\)\)|0_o|o_O|o\.O|:\(\(\(|\(\*_\*\)|O_o|:\]|;-\)|\^___\^|\(>_<\)|\(o:|:-P|:-\)\)\)|:D|o_0|<333|XDD|=3|:-o|:-3|=\(|:O|o\.0|:-X|:\||:-\]|>:o|V_V|\(;|8D|XD|:-/|\^_\^|:-O|<\.<|:/|>\.<|:x|=\|`)
