package tag

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// A language model is registered as bytes and built on first use, so that a
// program importing every language pays the decoding cost only for the ones
// it tags with.
type languageModel struct {
	data []byte
	once sync.Once
	tg   *Tagger
	err  error
}

var (
	languagesMu sync.RWMutex
	languages   = map[string]*languageModel{}
)

// RegisterLanguage makes a model available to ForLanguage under lang, a BCP 47
// tag such as "de". data is in prose's flat format and is retained.
//
// The packages under prose/v3/tag/lang call this from their init functions;
// import one to make its language available. Registering a language twice
// panics, for the same reason Register does.
func RegisterLanguage(lang string, data []byte) {
	lang = normalizeLang(lang)

	languagesMu.Lock()
	defer languagesMu.Unlock()

	if _, taken := languages[lang]; taken {
		panic("tag: language " + lang + " is already registered")
	}
	languages[lang] = &languageModel{data: data}
}

// ForLanguage returns the tagger for a language.
//
// lang is a BCP 47 tag; it is matched case-insensitively, and a tag with a
// region or script ("pt-BR", "zh-Hant") falls back to its base language when
// nothing more specific is registered. English is built in and uses the Penn
// Treebank tag set; every other language's model must have been registered,
// normally by importing its package under prose/v3/tag/lang, and uses the
// Universal tag set. A language with no model is an error rather than a
// silent fallback to English: tags from the wrong model are worse than none.
func ForLanguage(lang string) (*Tagger, error) {
	lang = normalizeLang(lang)
	if lang == "" || lang == "en" || strings.HasPrefix(lang, "en-") {
		return New()
	}

	languagesMu.RLock()
	m, ok := languages[lang]
	if !ok {
		base, _, _ := strings.Cut(lang, "-")
		if a, aliased := aliases[base]; aliased {
			base = a
		}
		m, ok = languages[base]
	}
	languagesMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("tag: no model for language %q; have %v", lang, Languages())
	}

	m.once.Do(func() {
		m.tg, m.err = FromBytes(m.data)
	})
	return m.tg, m.err
}

// LanguageModel returns the model registered for lang, in prose's flat format,
// as RegisterLanguage received it. It is for a program that ships models apart
// from its binary, as prose/v3/cmd/pack-models does; a program that only tags
// wants ForLanguage.
func LanguageModel(lang string) ([]byte, bool) {
	languagesMu.RLock()
	defer languagesMu.RUnlock()

	m, ok := languages[normalizeLang(lang)]
	if !ok {
		return nil, false
	}
	return m.data, true
}

// aliases maps language tags onto the registered model that serves them.
var aliases = map[string]string{
	"no": "nb", // Norwegian: the model was trained on Bokmål text
	"nn": "nb",
}

// Languages lists the languages with a model, sorted. English is always
// among them.
func Languages() []string {
	languagesMu.RLock()
	defer languagesMu.RUnlock()

	out := make([]string, 0, len(languages)+1)
	out = append(out, "en")
	for lang := range languages {
		out = append(out, lang)
	}
	sort.Strings(out)
	return out
}

// normalizeLang lower-cases a language tag and accepts an underscore in place
// of the hyphen, as locale names write it.
func normalizeLang(lang string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(lang), "_", "-"))
}
