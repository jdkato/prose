// Package ner finds named entities in tagged text.
//
// The classifier is a binary maximum-entropy model; see features.go for the
// feature set.
//
// The model loads once per process behind sync.OnceValues. Feature weights are
// looked up directly by key — one binary search over a sorted blob — rather
// than through an index indirection. Label selection is deterministic: ties
// break toward the label's position in the model's list.
//
// # Concurrency
//
// An Extracter is safe for concurrent use.
package ner

import (
	"math"
	"strings"
	"sync"

	"github.com/jdkato/prose/v3/internal/maxent"
	"github.com/jdkato/prose/v3/ner/maxentmodel"
	"github.com/jdkato/prose/v3/token"
)

// Token is prose's shared token type, re-exported for convenience.
type Token = token.Token

// Entity is a named entity found in text.
//
// Like Token, it carries its position, and the same invariant holds:
//
//	src[e.Start:e.End()] == e.Text
type Entity struct {
	// Text is the entity's content, verbatim from the source.
	Text string

	// Label is the entity type, such as "PERSON" or "GPE".
	Label string

	// Start is the byte offset of Text within the source.
	Start int
}

// End returns the byte offset just past this entity in the source.
func (e Entity) End() int { return e.Start + len(e.Text) }

// In returns the entity's text as located in src, or "" if the offsets do not
// fit. For any entity prose produces, e.In(src) == e.Text.
func (e Entity) In(src string) string {
	if e.Start < 0 || e.End() > len(src) {
		return ""
	}
	return src[e.Start:e.End()]
}

// Extracter labels tokens and groups them into entities.
type Extracter struct {
	model *maxent.Model
	pool  sync.Pool
}

// defaultModel loads the built-in English entity model on first use.
var defaultModel = sync.OnceValues(func() (*maxent.Model, error) {
	return maxent.Unmarshal(maxentmodel.English)
})

// New returns an Extracter backed by the built-in English model.
//
// The model is loaded on the first call and shared thereafter, so calling New
// repeatedly is cheap.
func New() (*Extracter, error) {
	m, err := defaultModel()
	if err != nil {
		return nil, err
	}
	return newWithModel(m), nil
}

// FromBytes returns an Extracter backed by a model in prose's flat format.
//
// data is retained and must not be modified afterwards.
func FromBytes(data []byte) (*Extracter, error) {
	m, err := maxent.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return newWithModel(m), nil
}

func newWithModel(m *maxent.Model) *Extracter {
	e := &Extracter{model: m}
	e.pool.New = func() any {
		return &scratch{
			scores: make([]float32, len(m.Labels())),
			key:    make([]byte, 0, 96),
			feats:  make(map[string]string, len(featureOrder)),
		}
	}
	return e
}

// scratch is the per-call working set, pooled so classification does not
// allocate a fresh score array and feature map for every token.
type scratch struct {
	scores  []float32
	key     []byte
	feats   map[string]string
	history []string
}

// Label assigns an IOB entity label to each token, in place.
//
// Tokens must already carry part-of-speech tags; the classifier uses them as
// features.
func (e *Extracter) Label(tokens []Token) {
	if len(tokens) == 0 {
		return
	}

	s := e.pool.Get().(*scratch)
	defer e.pool.Put(s)

	s.history = s.history[:0]
	labels := e.model.Labels()

	for i := range tokens {
		clear(s.feats)
		extractInto(s.feats, i, tokens, s.history)

		for j := range s.scores {
			s.scores[j] = 0
		}
		for j, label := range labels {
			s.scores[j] = e.score(s, label)
		}

		best := 0
		bestScore := float32(math.Inf(-1))
		for j := range labels {
			if s.scores[j] > bestScore {
				bestScore, best = s.scores[j], j
			}
		}

		tokens[i].Label = labels[best]
		s.history = append(s.history, simplePOS(labels[best]))
	}
}

// score sums the weights of the features present for one label.
//
// Features are walked and accumulated directly, with the key assembled in
// scratch space rather than joined into a new string per feature.
func (e *Extracter) score(s *scratch, label string) float32 {
	var total float32
	for _, name := range featureOrder {
		val, ok := s.feats[name]
		if !ok {
			val = ""
		}
		// The model's keys are "name-value-label"; build that in scratch
		// space rather than allocating a joined string per feature.
		s.key = s.key[:0]
		s.key = append(s.key, name...)
		s.key = append(s.key, '-')
		s.key = append(s.key, val...)
		s.key = append(s.key, '-')
		s.key = append(s.key, label...)

		if w, found := e.model.Weight(bytesToString(s.key)); found {
			total += w
		}
	}
	return total
}

// Chunk groups labeled tokens into entities.
//
// src is the text the tokens came from; it is used to recover each entity's
// exact span. Tokens must already be labeled by Label.
//
// Grouping extends an entity across matching part-of-speech tags and CD
// (numeral) runs, which is not what a textbook IOB chunker would do but is
// what prose's model was trained to produce.
func (e *Extracter) Chunk(src string, tokens []Token) []Entity {
	var (
		entities []Entity
		parts    []Token
		end      string
		idx      int
	)

	for _, tok := range tokens {
		label := tok.Label
		if (label != "O" && label != end) ||
			(idx > 0 && tok.Tag == parts[idx-1].Tag) ||
			(idx > 0 && tok.Tag == "CD" && parts[idx-1].Label != "O") {
			end = strings.Replace(label, "B", "I", 1)
			parts = append(parts, tok)
			idx++
		} else if (label == "O" && end != "") || label == end {
			// We have found the end of an entity.
			if label != "O" {
				parts = append(parts, tok)
			}
			entities = append(entities, coalesce(src, parts))

			end = ""
			parts = parts[:0]
			idx = 0
		}
	}

	return entities
}

// bytesToString views b as a string without copying. The result is used only
// for a lookup that does not retain it.
func bytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafeString(b)
}
