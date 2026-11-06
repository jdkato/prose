package segment

import (
	"bytes"
	"compress/gzip"
	"io"
	"sort"
	"strings"
	"sync"

	"gopkg.in/neurosnap/sentences.v1"
)

// A registered model is decoded on first use, so a program importing every
// language pays only for the ones it segments.
type languageModel struct {
	gz      []byte
	abbrevs []string
	once    sync.Once
	storage *sentences.Storage
	err     error
}

var (
	languagesMu sync.RWMutex
	languages   = map[string]*languageModel{}
)

// RegisterLanguage makes a Punkt model available to ForLanguage under lang, a
// BCP 47 tag such as "de". gz is the model's JSON, gzip-compressed, in the
// format gopkg.in/neurosnap/sentences reads. abbrevs are abbreviations the
// model should know but did not learn from its corpus, lower-cased and
// without the trailing period.
//
// The packages under prose/v3/segment/lang call this from their init
// functions; import one to make its language available. Registering a
// language twice panics: two packages disagreeing about a name is a bug, and
// picking one silently would make the choice depend on import order.
func RegisterLanguage(lang string, gz []byte, abbrevs ...string) {
	lang = normalizeLang(lang)

	languagesMu.Lock()
	defer languagesMu.Unlock()

	if _, taken := languages[lang]; taken {
		panic("segment: language " + lang + " is already registered")
	}
	languages[lang] = &languageModel{gz: gz, abbrevs: abbrevs}
}

// aliases maps language tags onto the registered model that serves them.
var aliases = map[string]string{
	"no": "nb", // Norwegian: the Punkt model was trained on Bokmål text
	"nn": "nb",
}

func (m *languageModel) load() (*sentences.Storage, error) {
	m.once.Do(func() {
		zr, err := gzip.NewReader(bytes.NewReader(m.gz))
		if err != nil {
			m.err = err
			return
		}
		b, err := io.ReadAll(zr)
		if err != nil {
			m.err = err
			return
		}
		m.storage, m.err = sentences.LoadTraining(b)
		if m.err != nil {
			return
		}
		// Added before anyone segments with it, so the shared storage is
		// never written to afterwards.
		for _, a := range m.abbrevs {
			m.storage.AbbrevTypes.Add(a)
		}
	})
	return m.storage, m.err
}

func normalizeLang(lang string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(lang), "_", "-"))
}

// ForLanguage returns a Segmenter for lang, a BCP 47 tag such as "de" or
// "pt-BR".
//
// English is built in. Sixteen more languages have a trained Punkt model,
// which knows the language's abbreviations, each in its own package under
// prose/v3/segment/lang: import the ones you need, or prose/v3/segment/lang/all
// for every one. They are Czech, Danish, Dutch, Estonian, Finnish, French,
// German, Greek, Italian, Norwegian, Polish, Portuguese, Slovene, Spanish,
// Swedish and Turkish.
//
// Any other language, or one whose package was not imported, gets a
// rule-based segmenter that works from punctuation alone and handles the East
// Asian terminators (。！？) that need no space after them. See Languages, and
// Trained to tell the two apart.
//
// A model is loaded on first use and shared, so calling this repeatedly is
// cheap, and the returned Segmenter is safe for concurrent use.
func ForLanguage(lang string) (*Segmenter, error) {
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
		base, _, _ := strings.Cut(lang, "-")
		return &Segmenter{rules: newRules(base)}, nil
	}
	storage, err := m.load()
	if err != nil {
		return nil, err
	}
	return newWithStorage(storage), nil
}

// Languages lists the languages with a trained segmentation model, sorted.
// English is always among them; the others are whatever packages under
// prose/v3/segment/lang the program imported.
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

// Trained reports whether the segmenter uses a trained model rather than
// punctuation rules.
func (s *Segmenter) Trained() bool { return s.rules == nil }
