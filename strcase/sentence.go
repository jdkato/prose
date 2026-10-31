package strcase

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/jdkato/prose/v3/internal/re"
	"strings"

	"github.com/jdkato/prose/v3/internal/strutil"
)

var reNumberList = re.MustCompile(`\d+\.`)
var defaultSentOpts = CaseOpts{
	vocab: []string{},
	indicator: func(word string, idx int) bool {
		if strings.HasSuffix(word, ":") {
			return true
		} else if idx == 0 && reNumberList.MatchString(word) {
			return true
		}
		return false
	},
}

// A SentenceConverter converts a string to sentence case.
type SentenceConverter struct {
	CaseOpts
	entries  []vocabEntry
	tokenRe  *re.Regexp
	prefixRe *regexp.Regexp
}

// NewSentenceConverter returns a new SentenceConverter.
func NewSentenceConverter(opts ...CaseOptFunc) *SentenceConverter {
	sent := new(SentenceConverter)

	base := defaultSentOpts
	for _, opt := range opts {
		opt(&base)
	}

	sent.vocab = base.vocab
	sent.indicator = base.indicator
	sent.prefix = base.prefix

	if sent.prefix != "" {
		sent.prefixRe = regexp.MustCompile(sent.prefix)
	}

	ps := `[\p{N}\p{L}*]+[^\s]*`
	if len(sent.vocab) > 0 {
		ps = fmt.Sprintf(`\b(?:%s)\b|%s`, strings.Join(sent.vocab, "|"), ps)
	}
	sent.tokenRe = re.MustCompile(`(?i)` + ps)
	sent.entries = compileVocab(sent.vocab)

	return sent
}

// Convert returns a copy of the string s in sentence case format.
func (sc *SentenceConverter) Convert(s string) string {
	var made []string

	prefix := ""
	if sc.prefixRe != nil && sc.prefixRe.MatchString(s) {
		prefix = sc.prefixRe.FindString(s)
		s = strings.TrimPrefix(s, prefix)
	}

	// Tokenize before lowering anything, so vocabulary entries keep the case
	// they were written with.
	tokens := sc.tokenRe.FindAllString(s, -1)

	for i, token := range tokens {
		prev := ""
		if i-1 >= 0 {
			prev = tokens[i-1]
		}

		if entry := lookupVocab(sc.entries, token); entry != "" {
			made = append(made, entry)
		} else if isAcronym(token) {
			made = append(made, token)
		} else if i == 0 || sc.indicator(prev, i-1) {
			made = append(made, strutil.ToTitle(token, true))
		} else {
			made = append(made, strings.ToLower(token))
		}
	}

	return prefix + strings.Join(made, " ")
}

// isAcronym reports whether token is written in capitals, like TAP or CLI,
// so sentence case keeps it as written. A plural or possessive (APIs, API's)
// and trailing punctuation don't count against it.
func isAcronym(token string) bool {
	token = strings.TrimRightFunc(token, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	for _, suffix := range []string{"'s", "’s", "s"} {
		if trimmed, ok := strings.CutSuffix(token, suffix); ok {
			token = trimmed
			break
		}
	}

	letters := 0
	for _, r := range token {
		if unicode.IsLower(r) {
			return false
		} else if unicode.IsLetter(r) {
			letters++
		}
	}
	return letters >= 2
}
