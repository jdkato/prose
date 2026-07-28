package strcase

import (
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/jdkato/prose/v3/internal/re"
	"github.com/jdkato/prose/v3/internal/strutil"
	"github.com/jdkato/prose/v3/tag"
	"github.com/jdkato/prose/v3/tokenize"
)

// Both are shared and built on first use, so importing this package costs
// nothing until something is actually converted.
var tokenizer = sync.OnceValue(func() *tokenize.TreebankWordTokenizer {
	return tokenize.NewTreebankWordTokenizer()
})

var tagger = sync.OnceValue(func() *tag.Tagger {
	t, err := tag.New()
	if err != nil {
		panic("prose/strcase: loading the tagger model: " + err.Error())
	}
	return t
})
var smallWords = []string{
	"a", "an", "and", "as", "at", "but", "by", "en", "for", "if", "in", "nor",
	"of", "on", "or", "per", "the", "to", "vs", "vs.", "via", "v", "v."}
var prepositions = []string{
	"with", "from", "into", "during", "including", "until", "against", "among",
	"throughout", "despite", "towards", "upon", "concerning", "about", "over",
	"through", "before", "between", "after", "since", "without", "under",
	"within", "along", "following", "across", "beyond", "around", "down",
	"near", "above"}
var splitRE = regexp.MustCompile(`[\p{N}\p{L}]+[^\s-/]*`)

// sanitizer replaces a set of Unicode characters with ASCII equivalents.
var sanitizer = strings.NewReplacer(
	"\u201c", `"`,
	"\u201d", `"`,
	"\u2018", "'",
	"\u2019", "'",
	"\u2013", "-",
	"\u2014", "-",
	"\u2026", "...")

// An IgnoreFunc is a TitleConverter callback that decides whether or not the
// the string word should be capitalized. firstOrLast indicates whether or not
// word is the first or last word in the given string.
type IgnoreFunc func(word string, tags []tag.Token, idx int, firstOrLast bool) bool

// A TitleConverter converts a string to title case according to its style.
type TitleConverter struct {
	CaseOpts
	ignore IgnoreFunc
}

// Styles for TitleConverter.
var (
	// APStyle capitalizes according to the Associated Press Stylebook.
	APStyle IgnoreFunc = optionsAP

	// ChicagoStyle capitalizes according to the Chicago Manual of Style.
	ChicagoStyle IgnoreFunc = optionsChicago
)

var defaultTitleOpts = CaseOpts{
	vocab: []string{},
	indicator: func(_ string, _ int) bool {
		return false
	},
}

// NewTitleConverter returns a new TitleConverter set to enforce the specified
// style.
func NewTitleConverter(style IgnoreFunc, opts ...CaseOptFunc) *TitleConverter {
	title := &TitleConverter{ignore: style}

	base := defaultTitleOpts
	for _, opt := range opts {
		opt(&base)
	}

	title.vocab = base.vocab
	title.indicator = base.indicator
	title.prefix = base.prefix

	return title
}

// Convert returns a copy of the string s in title case format.
func (tc *TitleConverter) Convert(s string) string {
	prefix := ""
	if tc.prefix != "" {
		if prefixRe := regexp.MustCompile(tc.prefix); prefixRe.MatchString(s) {
			prefix = prefixRe.FindString(s)
			s = strings.TrimPrefix(s, prefix)
		}
	}

	idx, pos := 0, 0
	t := sanitizer.Replace(s)
	end := len(t)

	// NOTE: We do thos because the tagger is sensitive to trailing punctuation
	// AND the initial case of the input.
	forTagging := s
	if !strutil.HasAnySuffix(s, []string{".", "!", "?"}) {
		forTagging = s + "."
	}
	words := tokenizer().Tokenize(forTagging)

	tags := tagger().Tag(words)
	widx := -1

	return prefix + splitRE.ReplaceAllStringFunc(s, func(m string) string {
		widx++

		sm := strings.ToLower(m)
		pos = strings.Index(t[idx:], m) + idx
		prev := strutil.CharAt(t, pos-1)
		ext := utf8.RuneCountInString(m)

		idx = pos + ext
		if found := tc.inVocab(m); found != "" {
			return found
		} else if tc.ignore(sm, tags, widx, pos == 0 || idx == end) &&
			(prev == ' ' || prev == '-' || prev == '/') &&
			strutil.CharAt(t, pos-2) != ':' && strutil.CharAt(t, pos-2) != '-' &&
			(strutil.CharAt(t, pos+ext) != '-' || strutil.CharAt(t, pos-1) == '-') {
			return sm
		}

		return strutil.ToTitle(m, false)
	})
}

func (tc *TitleConverter) inVocab(s string) string {
	for _, token := range tc.vocab {
		matched := re.MatchString(token, s)
		if strings.EqualFold(token, s) {
			return token
		} else if matched {
			return s
		}
	}
	return ""
}

// optionsAP implements AP-style casing.
//
//   - Capitalize the first word and the last word of the title
//   - Capitalize "to" in infinitives
//   - Do not capitalize articles, conjunctions, and prepositions of three
//     letters or fewer
//
// See testdata/AP.json for examples.
func optionsAP(word string, tags []tag.Token, idx int, bounding bool) bool {
	if word == "to" && idx+1 < len(tags) {
		return strings.HasPrefix(tags[idx+1].Tag, "NN")
	}
	return !bounding && strutil.StringInSlice(word, smallWords)
}

// ChicagoStyle states to lowercase articles (a, an, the), coordinating
// conjunctions (and, but, or, for, nor), and prepositions, regardless of
// length, unless they are the first or last word of the title.
func optionsChicago(word string, _ []tag.Token, _ int, bounding bool) bool {
	return !bounding && (strutil.StringInSlice(word, smallWords) || strutil.StringInSlice(word, prepositions))
}
