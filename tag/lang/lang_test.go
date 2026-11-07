package lang_test

import (
	"strings"
	"testing"

	"github.com/jdkato/prose/v3/tag"
	_ "github.com/jdkato/prose/v3/tag/lang/all"
	"github.com/jdkato/prose/v3/tokenize"
)

func tagsOf(t *testing.T, lang, text string) string {
	t.Helper()
	tagger, err := tag.ForLanguage(lang)
	if err != nil {
		t.Fatalf("[%s] %v", lang, err)
	}
	toks := tokenize.ForLanguage(lang).Tokenize(text)
	tagger.TagTokens(toks)
	out := make([]string, len(toks))
	for i, tok := range toks {
		out[i] = tok.Text + "/" + tok.Tag
	}
	return strings.Join(out, " ")
}

// Each model must get an ordinary sentence right. These are not corpus
// scores -- MODELS.md has those -- but a check that the model, the
// tokenizer and the tag set line up end to end.
func TestLanguages(t *testing.T) {
	cases := []struct{ lang, text, want string }{
		{"de", "Wir fahren morgen mit dem Zug nach Berlin.", "Wir/PRON fahren/VERB morgen/ADV mit/ADP dem/DET Zug/NOUN nach/ADP Berlin/PROPN ./PUNCT"},
		{"fr", "L'homme mange une pomme.", "L'/DET homme/NOUN mange/VERB une/DET pomme/NOUN ./PUNCT"},
		{"es", "¿Dónde está la casa de María?", "¿/PUNCT Dónde/ADV está/VERB la/DET casa/NOUN de/ADP María/PROPN ?/PUNCT"},
		{"it", "Il governo ha approvato ieri una nuova legge.", "Il/DET governo/NOUN ha/AUX approvato/VERB ieri/ADV una/DET nuova/ADJ legge/NOUN ./PUNCT"},
		{"pt", "O cachorro corre no parque.", "O/DET cachorro/NOUN corre/VERB no/ADP parque/NOUN ./PUNCT"},
		{"nl", "De kat slaapt op de bank.", "De/DET kat/NOUN slaapt/VERB op/ADP de/DET bank/NOUN ./PUNCT"},
		{"sv", "Hunden springer snabbt.", "Hunden/NOUN springer/VERB snabbt/ADV ./PUNCT"},
		{"ru", "Москва — большой город.", "Москва/PROPN —/PUNCT большой/ADJ город/NOUN ./PUNCT"},
		{"nb", "Regjeringen la fram et nytt forslag om skatt.", "Regjeringen/NOUN la/VERB fram/ADV et/DET nytt/ADJ forslag/NOUN om/ADP skatt/NOUN ./PUNCT"},
		{"da", "Hunden løber hurtigt.", "Hunden/NOUN løber/VERB hurtigt/ADV ./PUNCT"},
		{"fi", "Koira juoksee nopeasti.", "Koira/NOUN juoksee/VERB nopeasti/ADV ./PUNCT"},
		{"ca", "El gos corre ràpidament.", "El/DET gos/NOUN corre/VERB ràpidament/ADV ./PUNCT"},
		{"ro", "Guvernul a adoptat ieri o nouă lege.", "Guvernul/NOUN a/AUX adoptat/VERB ieri/ADV o/DET nouă/ADJ lege/NOUN ./PUNCT"},
		{"hr", "Djeca se igraju u parku.", "Djeca/NOUN se/PRON igraju/VERB u/ADP parku/NOUN ./PUNCT"},
		{"tr", "Çocuklar parkta top oynuyor.", "Çocuklar/NOUN parkta/NOUN top/NOUN oynuyor/VERB ./PUNCT"},
		{"id", "Anjing itu berlari cepat.", "Anjing/NOUN itu/DET berlari/VERB cepat/ADJ ./PUNCT"},
	}
	for _, tc := range cases {
		if got := tagsOf(t, tc.lang, tc.text); got != tc.want {
			t.Errorf("[%s] %q\n  got  %s\n  want %s", tc.lang, tc.text, got, tc.want)
		}
	}
}

func TestTagsetAndLanguage(t *testing.T) {
	de, err := tag.ForLanguage("de-AT")
	if err != nil {
		t.Fatal(err)
	}
	if de.Tagset() != tag.Universal || de.Language() != "de" {
		t.Errorf("German tagger reports %s/%s", de.Language(), de.Tagset())
	}
	en, err := tag.ForLanguage("en")
	if err != nil {
		t.Fatal(err)
	}
	if en.Tagset() != tag.Penn || en.Language() != "en" {
		t.Errorf("English tagger reports %s/%s", en.Language(), en.Tagset())
	}
	if no, err := tag.ForLanguage("no"); err != nil || no.Language() != "nb" {
		t.Errorf("Norwegian should resolve to the Bokmål model (err=%v)", err)
	}
	if _, err := tag.ForLanguage("ja"); err == nil {
		t.Error("Japanese has no model and must be an error, not English tags")
	}
	langs := tag.Languages()
	if len(langs) != 21 || langs[0] != "ca" || langs[len(langs)-1] != "tr" {
		t.Errorf("Languages() = %v", langs)
	}
}

// A universal model never emits a Penn tag, whatever the token looks like:
// the short-circuits for traces and hyphens are Penn conventions.
func TestUniversalModelEmitsOnlyUniversalTags(t *testing.T) {
	tagger, _ := tag.ForLanguage("de")
	known := map[string]bool{}
	for _, u := range tag.UniversalTags {
		known[u] = true
	}
	for _, tok := range tagger.Tag([]string{"-", "@jdkato", "*T*-1", "-LRB-", ":)", "2024", "e-mail"}) {
		if !known[tok.Tag] {
			t.Errorf("%q tagged %q, not a universal tag", tok.Text, tok.Tag)
		}
	}
}

// Taggers are built once per language and shared.
func TestForLanguageIsShared(t *testing.T) {
	a, _ := tag.ForLanguage("fr")
	b, _ := tag.ForLanguage("FR")
	if a != b {
		t.Error("ForLanguage should return the same tagger for one language")
	}
}

// A model taken out of the registry carries its provenance, which is what a
// program shipping it apart from the binary writes down beside it.
func TestLanguageModelMeta(t *testing.T) {
	data, ok := tag.LanguageModel("fr")
	if !ok {
		t.Fatal("French should be registered")
	}
	tg, err := tag.FromBytes(data)
	if err != nil {
		t.Fatal(err)
	}
	if tg.Language() != "fr" || tg.Tagset() != tag.Universal {
		t.Errorf("language %q, tagset %q", tg.Language(), tg.Tagset())
	}
	if !strings.Contains(tg.Meta(tag.MetaSource), "French") || tg.Meta(tag.MetaLicense) == "" {
		t.Errorf("source %q, license %q", tg.Meta(tag.MetaSource), tg.Meta(tag.MetaLicense))
	}
	if _, ok := tag.LanguageModel("xx"); ok {
		t.Error("an unregistered language should not have a model")
	}
}
