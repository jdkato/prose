package tokenize

import (
	"regexp"
	"strings"
	"sync"
)

// ForLanguage returns a tokenizer for lang, a BCP 47 tag such as "fr" or
// "pt-BR". The tokenizer is built on first use and shared, and is safe for
// concurrent use.
//
// Every language splits punctuation the same way; what differs is what a
// language does with apostrophes and hyphens, and whether it uses spaces at
// all:
//
//   - English splits contractions: "don't" is "do" and "n't".
//   - French and Italian split elisions: "l'homme" is "l'" and "homme", as
//     the Universal Dependencies treebanks have it; French also
//     splits hyphenated clitics, "donne-moi" into "donne" and "-moi".
//   - Chinese and Japanese are split where the script changes and at every
//     punctuation mark, which is as far as a tokenizer without a dictionary
//     can go. See SplittingScripts.
//
// A language not listed gets the language-neutral tokenizer: punctuation is
// split off, and words are otherwise whatever lies between spaces. The
// tokens of every language carry byte offsets into the input.
func ForLanguage(lang string) *Tokenizer {
	lang = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(lang), "_", "-"))
	if lang == "" {
		lang = "en"
	}

	// Rules are per language, not per region: "fr" and "fr-CA" tokenize
	// alike, so they share one tokenizer.
	base, _, _ := strings.Cut(lang, "-")
	if t, ok := shared.Load(base); ok {
		return t.(*Tokenizer)
	}

	build, ok := languages[base]
	if !ok {
		build = neutral
	}
	t, _ := shared.LoadOrStore(base, build())
	return t.(*Tokenizer)
}

// Languages lists the languages with tokenization rules of their own, sorted.
// Any other language gets the neutral tokenizer.
func Languages() []string {
	out := make([]string, 0, len(languages))
	for lang := range languages {
		out = append(out, lang)
	}
	sortStrings(out)
	return out
}

var shared sync.Map // lang -> *Tokenizer

var languages = map[string]func() *Tokenizer{
	"en": func() *Tokenizer { return New() },
	"fr": func() *Tokenizer {
		return New(
			UsingSplitCases(nil),
			UsingElision(frenchElision, "aujourd'hui", "aujourd’hui", "prud'homme", "prud'hommes",
				"prud’homme", "prud’hommes"),
			UsingClitics(frenchClitics),
		)
	},
	"it": func() *Tokenizer { return New(UsingSplitCases(nil), UsingElision(italianElision)) },
	"ja": func() *Tokenizer { return New(UsingSplitCases(nil), SplittingScripts(true)) },
	"zh": func() *Tokenizer { return New(UsingSplitCases(nil), SplittingScripts(true)) },
}

// neutral is the tokenizer for a language with no rules of its own: the
// English one without its contractions, which are English.
func neutral() *Tokenizer { return New(UsingSplitCases(nil)) }

// Both apostrophe forms are matched: the input is never rewritten, so the
// typographic apostrophe has to be recognized as it stands.
var (
	frenchElision = regexp.MustCompile(
		`^(?i:[cdjlmnst]|qu|jusqu|lorsqu|puisqu|quoiqu|quelqu|presqu|entr)['’]`)
	frenchClitics = regexp.MustCompile(
		`-(?i:t-il|t-elle|t-on|t-ils|t-elles|je|tu|il|elle|on|nous|vous|ils|elles|moi|toi|lui|leur|le|la|les|y|en|ce|ci|là)$`)
	italianElision = regexp.MustCompile(
		`^(?i:l|un|d|c|s|m|t|n|v|gl|dell|dall|nell|sull|all|coll|quell|quest|bell|sant|anch|com|tutt|senz|dov|cos|quand|finch|perch|poich|sicch)['’]`)
)

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
