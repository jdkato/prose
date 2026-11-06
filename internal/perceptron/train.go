package perceptron

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"

	"github.com/jdkato/prose/v3/internal/flat"
)

// Sentence is one tagged training sentence: Words[i] carries Tags[i].
type Sentence struct {
	Words []string
	Tags  []string
}

// Options tune training.
type Options struct {
	// Template selects the feature set. New models should use TemplateRunes.
	Template Template

	// Iterations is how many passes to make over the data. Five is the
	// usual number; more rarely helps an averaged perceptron.
	Iterations int

	// Seed orders the shuffling between passes, so a run is reproducible.
	Seed int64

	// MinFreq and MaxAmbiguity decide which words go in the tag map: a word
	// seen at least MinFreq times, with one tag at least MaxAmbiguity of
	// those times, is assigned that tag outright and never scored.
	MinFreq      int
	MaxAmbiguity float64

	// MinCount drops features seen fewer times than this in the training
	// data, counted with the gold tags in place. A feature seen once is a
	// memory of one token, not a generalization; leaving such features out
	// halves the model and costs nothing measurable.
	MinCount int

	// Prune drops averaged weights whose magnitude is below it. Most of a
	// trained model is weights near zero that decide nothing; dropping them
	// shrinks the file several-fold at a cost of a few hundredths of a
	// point.
	Prune float64

	// Fixed, if set, assigns tags the tagger decides without the model --
	// an emoticon is a symbol whatever the context -- so training sees the
	// same tag sequence the tagger will produce.
	Fixed func(word string) (string, bool)

	// Progress, if set, is called after each pass with the pass number and
	// the training accuracy of that pass.
	Progress func(pass int, accuracy float64)
}

// Stats describes a trained model.
type Stats struct {
	Classes  int
	Features int
	Pairs    int
	TagMap   int
}

// entry is one (feature, class) weight with the bookkeeping that averaging
// needs: the running total of the weight over time and when it last changed.
type entry struct {
	class  uint8
	weight float64
	total  float64
	stamp  int
}

type trainer struct {
	opts    Options
	classes []string
	classID map[string]uint8

	featID  map[string]int
	entries [][]entry       // by feature ID
	allowed map[string]bool // features frequent enough to learn; nil for all

	tagMap map[string]uint8
	feats  *Features
	ctx    []string
	scores []float64
	n      int // words scored so far: the clock the averaging runs on
}

// Train fits an averaged perceptron to sents and returns it in the flat
// format's source shape, ready for flat.Marshal.
func Train(sents []Sentence, opts Options) (*flat.Source, Stats, error) {
	if len(sents) == 0 {
		return nil, Stats{}, fmt.Errorf("perceptron: no training sentences")
	}
	if opts.Iterations <= 0 {
		opts.Iterations = 5
	}
	if opts.MinFreq <= 0 {
		opts.MinFreq = 20
	}
	if opts.MaxAmbiguity <= 0 {
		opts.MaxAmbiguity = 0.97
	}

	t := &trainer{
		opts:    opts,
		classID: map[string]uint8{},
		featID:  map[string]int{},
		tagMap:  map[string]uint8{},
		feats:   NewFeatures(opts.Template),
	}
	if err := t.collectClasses(sents); err != nil {
		return nil, Stats{}, err
	}
	t.scores = make([]float64, len(t.classes))
	t.makeTagMap(sents)
	if opts.MinCount > 1 {
		t.countFeatures(sents)
	}

	rng := rand.New(rand.NewSource(opts.Seed)) //nolint:gosec // reproducible shuffling, not security
	order := make([]int, len(sents))
	for i := range order {
		order[i] = i
	}

	for pass := 1; pass <= opts.Iterations; pass++ {
		rng.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })

		var right, total int
		for _, idx := range order {
			r, n := t.trainSentence(sents[idx])
			right += r
			total += n
		}
		if opts.Progress != nil && total > 0 {
			opts.Progress(pass, float64(right)/float64(total))
		}
	}

	return t.average()
}

// collectClasses fixes the class list: sorted, so class IDs -- and with them
// the tie-breaking in flat.Model.Best -- do not depend on data order.
func (t *trainer) collectClasses(sents []Sentence) error {
	seen := map[string]bool{}
	for _, s := range sents {
		if len(s.Words) != len(s.Tags) {
			return fmt.Errorf("perceptron: sentence has %d words and %d tags", len(s.Words), len(s.Tags))
		}
		for _, tag := range s.Tags {
			seen[tag] = true
		}
	}
	for c := range seen {
		t.classes = append(t.classes, c)
	}
	sort.Strings(t.classes)
	if len(t.classes) > flat.MaxClasses {
		return fmt.Errorf("perceptron: %d classes, format allows %d", len(t.classes), flat.MaxClasses)
	}
	for i, c := range t.classes {
		t.classID[c] = uint8(i)
	}
	return nil
}

// makeTagMap assigns outright the words that are frequent and unambiguous.
func (t *trainer) makeTagMap(sents []Sentence) {
	counts := map[string]map[string]int{}
	for _, s := range sents {
		for i, w := range s.Words {
			m := counts[w]
			if m == nil {
				m = map[string]int{}
				counts[w] = m
			}
			m[s.Tags[i]]++
		}
	}
	for w, tags := range counts {
		var n, best int
		var bestTag string
		for tag, c := range tags {
			n += c
			if c > best {
				best, bestTag = c, tag
			}
		}
		if n >= t.opts.MinFreq && float64(best)/float64(n) >= t.opts.MaxAmbiguity {
			t.tagMap[w] = t.classID[bestTag]
		}
	}
}

// countFeatures finds the features that occur at least MinCount times, using
// the gold tags for the previous-tag features.
func (t *trainer) countFeatures(sents []Sentence) {
	// Keys are interned once and counted by ID. Assigning through a map
	// with the key's transient view would have the map keep that view --
	// Go replaces a string key on update -- and the buffer behind it is
	// overwritten by the next Build.
	ids := map[string]int{}
	counts := []int{}
	for _, s := range sents {
		t.ctx = Context(t.ctx, s.Words)
		p1, p2 := Start, Start2
		for i, word := range s.Words {
			if _, ok := t.tagMap[word]; !ok {
				t.feats.Build(t.ctx, i+2, word, p1, p2)
				for k, n := 0, t.feats.Len(); k < n; k++ {
					key := t.feats.Key(k)
					id, seen := ids[key]
					if !seen {
						id = len(counts)
						ids[strings.Clone(key)] = id
						counts = append(counts, 0)
					}
					counts[id]++
				}
			}
			p2, p1 = p1, s.Tags[i]
		}
	}
	t.allowed = make(map[string]bool, len(ids)/2)
	for key, id := range ids {
		if counts[id] >= t.opts.MinCount {
			t.allowed[key] = true
		}
	}
}

// trainSentence runs one pass over a sentence, updating the weights wherever
// the model guesses wrong. It reports how many words it got right.
//
// The previous-tag features use the model's own guesses rather than the gold
// tags, as the tagger will at run time.
func (t *trainer) trainSentence(s Sentence) (right, total int) {
	t.ctx = Context(t.ctx, s.Words)
	p1, p2 := Start, Start2

	for i, word := range s.Words {
		truth := s.Tags[i]
		var guess string

		switch {
		case t.opts.Fixed != nil && fixedTag(t.opts.Fixed, word) != "":
			guess = fixedTag(t.opts.Fixed, word)
		default:
			if id, ok := t.tagMap[word]; ok {
				guess = t.classes[id]
				break
			}
			t.feats.Build(t.ctx, i+2, word, p1, p2)
			ids := t.featureIDs()
			guess = t.predict(ids)
			// The clock advances on every scored word, not only on the
			// wrong ones: averaging weighs each weight by how long it
			// stood, and a weight that stood through a correct answer
			// earned that time.
			t.n++
			if guess != truth {
				t.update(ids, t.classID[truth], t.classID[guess])
			}
		}

		if guess == truth {
			right++
		}
		total++
		p2, p1 = p1, guess
	}
	return right, total
}

func fixedTag(fixed func(string) (string, bool), word string) string {
	if tag, ok := fixed(word); ok {
		return tag
	}
	return ""
}

// featureIDs interns the built features, creating any that are new and
// leaving out any that are too rare to learn.
func (t *trainer) featureIDs() []int {
	ids := make([]int, 0, t.feats.Len())
	for k, n := 0, t.feats.Len(); k < n; k++ {
		key := t.feats.Key(k)
		id, ok := t.featID[key]
		if !ok {
			if t.allowed != nil && !t.allowed[key] {
				continue
			}
			id = len(t.entries)
			t.featID[strings.Clone(key)] = id
			t.entries = append(t.entries, nil)
		}
		ids = append(ids, id)
	}
	return ids
}

// predict scores the classes and returns the best, breaking ties toward the
// lowest class ID exactly as flat.Model.Best does.
func (t *trainer) predict(ids []int) string {
	for c := range t.scores {
		t.scores[c] = 0
	}
	for _, id := range ids {
		for _, e := range t.entries[id] {
			t.scores[e.class] += e.weight
		}
	}
	best, bestID := math.Inf(-1), 0
	for c, s := range t.scores {
		if s > best {
			best, bestID = s, c
		}
	}
	return t.classes[bestID]
}

// update moves each feature's weight toward truth and away from guess.
func (t *trainer) update(ids []int, truth, guess uint8) {
	for _, id := range ids {
		t.bump(id, truth, 1)
		t.bump(id, guess, -1)
	}
}

func (t *trainer) bump(id int, class uint8, delta float64) {
	entries := t.entries[id]
	for k := range entries {
		e := &entries[k]
		if e.class == class {
			e.total += float64(t.n-e.stamp) * e.weight
			e.weight += delta
			e.stamp = t.n
			return
		}
	}
	t.entries[id] = append(entries, entry{class: class, weight: delta, stamp: t.n})
}

// average replaces each weight with its mean over every update, which is
// what makes the perceptron stable, then drops the weights that round to
// nothing.
func (t *trainer) average() (*flat.Source, Stats, error) {
	src := &flat.Source{
		Weights: make(map[string]map[string]float64, len(t.featID)),
		TagMap:  make(map[string]string, len(t.tagMap)),
		Classes: t.classes,
		Meta:    map[string]string{flat.MetaTemplate: fmt.Sprint(int(t.opts.Template))},
	}
	for w, id := range t.tagMap {
		src.TagMap[w] = t.classes[id]
	}

	var pairs int
	for feat, id := range t.featID {
		var inner map[string]float64
		for _, e := range t.entries[id] {
			total := e.total + float64(t.n-e.stamp)*e.weight
			avg := total / float64(max(t.n, 1))
			avg = math.Round(avg*1000) / 1000
			if avg == 0 || math.Abs(avg) < t.opts.Prune {
				continue
			}
			if inner == nil {
				inner = map[string]float64{}
			}
			inner[t.classes[e.class]] = avg
		}
		if len(inner) > 0 {
			src.Weights[feat] = inner
			pairs += len(inner)
		}
	}

	return src, Stats{
		Classes:  len(t.classes),
		Features: len(src.Weights),
		Pairs:    pairs,
		TagMap:   len(src.TagMap),
	}, nil
}
