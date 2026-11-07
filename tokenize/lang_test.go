package tokenize_test

import (
	"strings"
	"testing"

	"github.com/jdkato/prose/v3/tokenize"
)

func texts(toks []tokenize.Token) string {
	out := make([]string, len(toks))
	for i, t := range toks {
		out[i] = t.Text
	}
	return strings.Join(out, "|")
}

// TestLanguageRules pins what each language does with apostrophes, hyphens,
// quotation styles and scripts.
func TestLanguageRules(t *testing.T) {
	cases := []struct {
		lang string
		in   string
		want string
	}{
		// French elides articles and pronouns, keeps the few words that own
		// their apostrophe, and hyphenates clitics onto the verb.
		{"fr", "L'homme qu'il a vu.", "L'|homme|qu'|il|a|vu|."},
		{"fr", "C’est aujourd'hui.", "C’|est|aujourd'hui|."},
		{"fr", "Donne-moi le livre, dit-il.", "Donne|-moi|le|livre|,|dit|-il|."},
		{"fr", "Y a-t-il un problème ?", "Y|a|-t-il|un|problème|?"},
		{"fr", "« Bonjour », dit-elle.", "«|Bonjour|»|,|dit|-elle|."},
		// English contractions do not apply elsewhere.
		{"fr", "Il n'est pas là.", "Il|n'|est|pas|là|."},
		{"de", "Peter's Haus", "Peter's|Haus"},
		{"en", "Peter's house", "Peter|'s|house"},

		{"it", "L'acqua dell'amore un'ora.", "L'|acqua|dell'|amore|un'|ora|."},

		// German and Russian quotation marks.
		{"de", "Er sagte „Hallo“ zu mir.", "Er|sagte|„|Hallo|“|zu|mir|."},
		{"ru", "Он сказал «привет» мне.", "Он|сказал|«|привет|»|мне|."},
		{"es", "¿Qué pasa? ¡Nada!", "¿|Qué|pasa|?|¡|Nada|!"},

		// Chinese and Japanese: script runs and punctuation.
		{"ja", "iPhoneを買った。", "iPhone|を|買|った|。"},
		{"ja", "東京タワーに行きました！", "東京|タワー|に|行|きました|！"},
		{"zh", "我爱北京，天安门。", "我爱北京|，|天安门|。"},

		// A language with no rules of its own is the neutral tokenizer.
		{"xx", "Don't stop.", "Don't|stop|."},
		{"pt-BR", "Olá, mundo…", "Olá|,|mundo|…"},
	}

	for _, tc := range cases {
		got := texts(tokenize.ForLanguage(tc.lang).Tokenize(tc.in))
		if got != tc.want {
			t.Errorf("[%s] Tokenize(%q)\n  got  %s\n  want %s", tc.lang, tc.in, got, tc.want)
		}
	}
}

// Offsets must hold in every language, elision and script splitting included.
func TestLanguageOffsets(t *testing.T) {
	cases := map[string]string{
		"fr": "« L'homme », dit-elle, n'est pas là aujourd'hui.",
		"ja": "東京タワーにiPhoneを持って行きました！「はい」と言った。",
		"zh": "我爱北京（天安门），你呢？",
		"it": "Dell'acqua, un'ora, l'amore.",
		"de": "„Zum Beispiel“, sagte er – ‚wirklich‘.",
	}
	for lang, src := range cases {
		for _, tok := range tokenize.ForLanguage(lang).Tokenize(src) {
			if tok.In(src) != tok.Text {
				t.Errorf("[%s] token %q at %d: source has %q", lang, tok.Text, tok.Start, tok.In(src))
			}
			if strings.TrimSpace(tok.Text) == "" {
				t.Errorf("[%s] blank token at %d", lang, tok.Start)
			}
		}
	}
}

func TestForLanguageIsShared(t *testing.T) {
	if tokenize.ForLanguage("FR") != tokenize.ForLanguage("fr-ca") {
		t.Error("the French tokenizer should be built once, however the tag is spelled")
	}
	if tokenize.ForLanguage("") != tokenize.ForLanguage("en") {
		t.Error("an empty language should mean English")
	}
	langs := tokenize.Languages()
	if len(langs) != 5 || langs[0] != "en" {
		t.Errorf("Languages() = %v", langs)
	}
}
