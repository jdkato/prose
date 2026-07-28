// Package segment splits text into sentences that know where they came from.
//
// Segmentation uses the Punkt algorithm via gopkg.in/neurosnap/sentences, with
// prose's customizations for English (see punkt.go).
//
// The trained model is immutable once built, so it is loaded once per process
// and shared. Sentences carry their byte offsets, so a caller never has to go
// looking for where one started — which would be both slow and wrong for
// repeated sentences.
package segment

import (
	"sync"
	"unicode"
	"unicode/utf8"

	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/data"
)

// Sentence is a segmented span of text and its position in the source.
//
// The same invariant the token package documents holds here:
//
//	src[s.Start:s.End()] == s.Text
//
// Leading and trailing whitespace is trimmed from Text, and Start is advanced
// to match, so the offsets stay accurate rather than merely close. Enforced by
// TestSentenceOffsetsLocateTheirText.
type Sentence struct {
	// Text is the sentence, whitespace-trimmed.
	Text string

	// Start is the byte offset of Text within the source.
	Start int
}

// End returns the byte offset just past this sentence in the source.
func (s Sentence) End() int { return s.Start + len(s.Text) }

// In returns the sentence's text as located in src, or "" if the offsets do
// not fit. For any sentence prose produces, s.In(src) == s.Text.
func (s Sentence) In(src string) string {
	if s.Start < 0 || s.End() > len(src) {
		return ""
	}
	return src[s.Start:s.End()]
}

// Segmenter splits text into sentences.
type Segmenter struct {
	tokenizer *sentences.DefaultSentenceTokenizer
}

// trainedStorage loads and prepares the English punkt model exactly once.
//
// The supervised abbreviations are added here rather than per-Segmenter so the
// shared storage is already complete and is never mutated afterwards — which
// is what makes sharing a Segmenter across goroutines safe.
var trainedStorage = sync.OnceValues(func() (*sentences.Storage, error) {
	b, err := data.Asset("data/english.json")
	if err != nil {
		return nil, err
	}
	training, err := sentences.LoadTraining(b)
	if err != nil {
		return nil, err
	}
	for _, abbr := range supervisedAbbrevs {
		training.AbbrevTypes.Add(abbr)
	}
	return training, nil
})

var supervisedAbbrevs = []string{"sgt", "gov", "no", "mt"}

// New returns a Segmenter using the built-in English model.
//
// The model is loaded on the first call and shared thereafter, so calling New
// repeatedly is cheap. The returned Segmenter is safe for concurrent use.
func New() (*Segmenter, error) {
	training, err := trainedStorage()
	if err != nil {
		return nil, err
	}
	return newWithStorage(training), nil
}

// WithTraining returns a Segmenter using a caller-supplied punkt model.
//
// The storage is not copied. If the caller mutates it afterwards — by adding
// abbreviations, say — the Segmenter is no longer safe to use concurrently.
func WithTraining(training *sentences.Storage) *Segmenter {
	return newWithStorage(training)
}

func newWithStorage(training *sentences.Storage) *Segmenter {
	lang := sentences.NewPunctStrings()
	word := newWordTokenizer(lang)
	annotations := sentences.NewAnnotations(training, lang, word)

	ortho := &sentences.OrthoContext{
		Storage:      training,
		PunctStrings: lang,
		TokenType:    word,
		TokenFirst:   word,
	}

	annotations = append(annotations, &multiPunctWordAnnotation{
		Storage:      training,
		TokenParser:  word,
		TokenGrouper: &sentences.DefaultTokenGrouper{},
		Ortho:        ortho,
	})

	return &Segmenter{
		tokenizer: &sentences.DefaultSentenceTokenizer{
			Storage:       training,
			PunctStrings:  lang,
			WordTokenizer: word,
			Annotations:   annotations,
		},
	}
}

// Segment splits text into sentences.
//
// Sentences that are entirely whitespace are dropped rather than returned as
// empty spans.
func (s *Segmenter) Segment(text string) []Sentence {
	raw := s.tokenizer.Tokenize(text)

	out := make([]Sentence, 0, len(raw))
	for _, r := range raw {
		start, end := trimRange(text, r.Start, r.End)
		if start >= end {
			continue
		}
		out = append(out, Sentence{Text: text[start:end], Start: start})
	}
	return out
}

// SegmentText returns just the sentence strings, for callers that do not care
// where they came from.
func (s *Segmenter) SegmentText(text string) []string {
	sents := s.Segment(text)
	out := make([]string, len(sents))
	for i, sent := range sents {
		out[i] = sent.Text
	}
	return out
}

// trimRange narrows [start,end) past surrounding whitespace.
//
// strings.TrimSpace would produce the same text, but it throws away how much
// was trimmed, leaving the reported offset pointing at whitespace the text no
// longer contains. Trimming the range keeps text[start:end] equal to the
// sentence.
func trimRange(text string, start, end int) (int, int) {
	if start < 0 {
		start = 0
	}
	if end > len(text) {
		end = len(text)
	}
	for start < end {
		r, size := utf8.DecodeRuneInString(text[start:end])
		if size == 0 || !unicode.IsSpace(r) {
			break
		}
		start += size
	}
	for end > start {
		r, size := utf8.DecodeLastRuneInString(text[start:end])
		if size == 0 || !unicode.IsSpace(r) {
			break
		}
		end -= size
	}
	return start, end
}
