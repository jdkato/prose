package segment

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// ruleSegmenter finds sentence boundaries without a trained model.
//
// Punkt learns a language's abbreviations from a corpus, and prose ships its
// models for seventeen languages. For the rest, this does what can be done
// from punctuation alone: a sentence ends at a terminator that is followed by
// a space and then something that can begin a sentence, unless the word before
// the terminator looks like an abbreviation or an initial. East Asian
// terminators (。！？) end a sentence outright, since those scripts put no
// space between sentences.
//
// It is deliberately simple, and documented as such. It gets ordinary prose
// right and errs toward not breaking, which for a linter means a sentence
// that is occasionally two rather than a word that is occasionally torn in
// half.
type ruleSegmenter struct {
	// abbrevs holds the language's common abbreviations, lower-cased and
	// without their period, so "ул. Ленина" is not two sentences.
	abbrevs map[string]bool
}

// newRules returns the rule-based segmenter for lang, with that language's
// abbreviations if it is one of the few listed below.
func newRules(lang string) *ruleSegmenter {
	r := &ruleSegmenter{abbrevs: map[string]bool{}}
	for _, a := range strings.Fields(ruleAbbrevs[lang]) {
		r.abbrevs[a] = true
	}
	return r
}

// ruleAbbrevs lists a few common abbreviations, by language, that end in a
// period and do not end a sentence: titles, "etc." and the like. Dotted
// abbreviations (т.е., EE.UU.) and initials need no listing.
var ruleAbbrevs = map[string]string{
	"ru": "г гг ул д кв стр см др проф тел руб коп млн млрд тыс пр им ст св обл дер пос напр",
	"uk": "р рр вул буд кв див ін проф тел грн млн млрд тис ст св обл напр",
	"ca": "sr sra srta dr dra prof av pàg tel núm aprox etc ex ref vol cap fig p pp ed",
	"ro": "dl dna dra dr prof nr str bd tel pag vol cap fig ex etc aprox cca",
	"hr": "g gđa gđica dr prof ing mr br str tel sv tzv npr itd tj sl god",
	"sr": "g gđa gđica dr prof ing mr br str tel sv tzv npr itd tj sl god",
	"lv": "g k dr prof nr lpp tel u.c piem utt t.i sk",
	"lt": "p dr prof nr psl tel pvz t.y ir t.t žr g",
	"hu": "dr prof úr pl stb kb ún ill tel u sz évf old",
	"id": "bpk ibu sdr dr prof ir drs hal tel no jl dll dsb yth",
	"ko": "",
	"ja": "",
	"zh": "",
}

// split returns the sentence spans of text. Every byte of text falls in
// exactly one span, so the caller can trim them without losing anything.
func (rs *ruleSegmenter) split(text string) []span {
	var out []span
	start := 0

	for i := 0; i < len(text); {
		r, size := utf8.DecodeRuneInString(text[i:])
		next := i + size

		switch {
		case isEastAsianTerminator(r):
			// Any run of terminators and closers belongs to the sentence:
			// "本当に？！」" ends after the bracket.
			next = skipClosers(text, skipWhile(text, next, isEastAsianTerminator))
			out = append(out, span{start, next})
			start = next

		case r == '.' || r == '!' || r == '?' || r == '…':
			end := skipWhile(text, next, func(r rune) bool { return r == '.' || r == '!' || r == '?' || r == '…' })
			end = skipClosers(text, end)
			if rs.breaksAt(text, start, i, end) {
				out = append(out, span{start, end})
				start = end
			}
			next = end
		}
		i = next
	}

	if start < len(text) {
		out = append(out, span{start, len(text)})
	}
	return out
}

// breaksAt decides whether the terminator at i, whose closing punctuation
// runs to end, ends the sentence that began at start.
func (rs *ruleSegmenter) breaksAt(text string, start, i, end int) bool {
	// What follows must be a space and then something that can open a
	// sentence. No space means a decimal point, a domain name or a word
	// like "Yahoo!"; a lower-case letter means the writer went on -- also
	// when a dash stands between, as in dialogue: «Привет!» — сказал он.
	j := end
	if j >= len(text) {
		return true
	}
	r, size := utf8.DecodeRuneInString(text[j:])
	if !unicode.IsSpace(r) {
		return false
	}
	for j < len(text) && (unicode.IsSpace(r) || unicode.Is(unicode.Pd, r)) {
		j += size
		r, size = utf8.DecodeRuneInString(text[j:])
	}
	if j >= len(text) {
		return true
	}
	if unicode.IsLower(r) {
		return false
	}

	if text[i] != '.' {
		return true
	}

	// A period after an initial ("J. Smith"), after a dotted abbreviation
	// ("z.B.", "т.е.", "EE.UU.") or after a short run of digits ("3. Punkt")
	// is probably not the end of the sentence.
	word := text[start:i]
	if k := strings.LastIndexFunc(word, unicode.IsSpace); k >= 0 {
		word = word[k+1:]
	}
	word = strings.TrimLeftFunc(word, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })

	switch {
	case word == "":
		return true
	case utf8.RuneCountInString(word) == 1 && isLetter(word):
		return false
	case strings.Contains(word, "."):
		return false
	case isDigits(word) && len(word) <= 2:
		return false
	case rs.abbrevs[strings.ToLower(word)]:
		return false
	}
	return true
}

func isLetter(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsLetter(r)
}

func isDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return s != ""
}

func isEastAsianTerminator(r rune) bool {
	switch r {
	case '。', '！', '？', '｡', '‼', '⁇', '⁈', '⁉':
		return true
	}
	return false
}

// skipWhile advances i past runes satisfying f.
func skipWhile(text string, i int, f func(rune) bool) int {
	for i < len(text) {
		r, size := utf8.DecodeRuneInString(text[i:])
		if !f(r) {
			break
		}
		i += size
	}
	return i
}

// skipClosers advances i past the punctuation that may close a sentence after
// its terminator: quotation marks and brackets of every style.
func skipClosers(text string, i int) int {
	return skipWhile(text, i, func(r rune) bool {
		switch r {
		case '"', '\'', ')', ']', '}', '”', '’', '»', '“', '‘', '›', '」', '』', '）', '】', '》', '〉', '〕':
			return true
		}
		return false
	})
}
