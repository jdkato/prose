package strcase

import (
	"fmt"
	"regexp"

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

	return sent
}

// Convert returns a copy of the string s in sentence case format.
func (sc *SentenceConverter) Convert(s string) string {
	var made []string

	prefix := ""
	if sc.prefix != "" {
		if prefixRe := regexp.MustCompile(sc.prefix); prefixRe.MatchString(s) {
			prefix = prefixRe.FindString(s)
			s = strings.TrimPrefix(s, prefix)
		}
	}

	ps := `[\p{N}\p{L}*]+[^\s]*`
	if len(sc.vocab) > 0 {
		ps = fmt.Sprintf(`\b(?:%s)\b|%s`, strings.Join(sc.vocab, "|"), ps)
	}
	re := re.MustCompile(`(?i)` + ps)

	// Tokenize before lowering anything, so vocabulary entries keep the case
	// they were written with.
	tokens := re.FindAllString(s, -1)

	for i, token := range tokens {
		prev := ""
		if i-1 >= 0 {
			prev = tokens[i-1]
		}

		if entry := sc.inVocab(token); entry != "" {
			made = append(made, entry)
		} else if i == 0 || sc.indicator(prev, i-1) {
			made = append(made, strutil.ToTitle(token, true))
		} else {
			made = append(made, strings.ToLower(token))
		}
	}

	return prefix + strings.Join(made, " ")
}

func (sc *SentenceConverter) inVocab(s string) string {
	for _, token := range sc.vocab {
		matched := re.MatchString(token, s)
		if strings.EqualFold(token, s) {
			return token
		} else if matched {
			return s
		}
	}
	return ""
}
