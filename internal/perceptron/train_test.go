package perceptron

import (
	"testing"

	"github.com/jdkato/prose/v3/internal/flat"
)

func corpus() []Sentence {
	var out []Sentence
	for i := 0; i < 30; i++ {
		out = append(out,
			Sentence{Words: []string{"der", "Hund", "bellt", "."}, Tags: []string{"DET", "NOUN", "VERB", "PUNCT"}},
			Sentence{Words: []string{"die", "Katze", "schläft", "."}, Tags: []string{"DET", "NOUN", "VERB", "PUNCT"}},
			Sentence{Words: []string{"Anna", "läuft", "schnell", "."}, Tags: []string{"PROPN", "VERB", "ADV", "PUNCT"}},
			Sentence{Words: []string{"er", "geht", "nach", "Hause", "."}, Tags: []string{"PRON", "VERB", "ADP", "NOUN", "PUNCT"}},
		)
	}
	return out
}

// A trained model must round-trip through the flat format and tag its own
// training data, including the unambiguous words it fixed outright.
func TestTrainRoundTrip(t *testing.T) {
	src, stats, err := Train(corpus(), Options{Template: TemplateRunes, Seed: 1, MinFreq: 100, Iterations: 20})
	if err != nil {
		t.Fatal(err)
	}
	if stats.Classes != 8 {
		t.Errorf("classes = %d, want 8", stats.Classes)
	}
	if src.TagMap["."] != "PUNCT" {
		t.Errorf("tag map: %q is %q, want PUNCT", ".", src.TagMap["."])
	}
	if src.Meta[flat.MetaTemplate] != "1" {
		t.Errorf("template meta = %q", src.Meta[flat.MetaTemplate])
	}

	data, err := flat.Marshal(*src)
	if err != nil {
		t.Fatal(err)
	}
	m, err := flat.Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}

	feats := NewFeatures(TemplateRunes)
	scores := make([]float32, len(m.Classes()))
	for _, s := range corpus()[:4] {
		ctx := Context(nil, s.Words)
		p1, p2 := Start, Start2
		for i, w := range s.Words {
			got, ok := m.TagFor(w)
			if !ok {
				for j := range scores {
					scores[j] = 0
				}
				feats.Build(ctx, i+2, w, p1, p2)
				for k := 0; k < feats.Len(); k++ {
					m.AccumulateInto(feats.Key(k), 1, scores)
				}
				got = m.Best(scores)
			}
			if got != s.Tags[i] {
				t.Errorf("%q tagged %s, want %s", w, got, s.Tags[i])
			}
			p2, p1 = p1, got
		}
	}
}

// The cutoff must drop rare features and keep frequent ones.
func TestMinCountDropsRareFeatures(t *testing.T) {
	sents := corpus()
	sents = append(sents, Sentence{Words: []string{"Zebra"}, Tags: []string{"NOUN"}})

	_, all, err := Train(sents, Options{Template: TemplateRunes, MinCount: 1, MinFreq: 1000})
	if err != nil {
		t.Fatal(err)
	}
	_, cut, err := Train(sents, Options{Template: TemplateRunes, MinCount: 2, MinFreq: 1000})
	if err != nil {
		t.Fatal(err)
	}
	if cut.Features >= all.Features {
		t.Errorf("min-count 2 kept %d features, min-count 1 kept %d", cut.Features, all.Features)
	}
}

func TestFeaturesKeysAreStableViews(t *testing.T) {
	f := NewFeatures(TemplateRunes)
	ctx := Context(nil, []string{"Große", "Straße"})
	f.Build(ctx, 3, "Straße", "ADJ", Start)
	keys := f.Strings()
	if len(keys) != f.Len() || keys[1] != "b"+"aße" || keys[2] != "cS" {
		t.Errorf("keys = %q", keys)
	}
	if got := Shape("iPhone"); got != "xX" {
		t.Errorf("Shape(iPhone) = %q", got)
	}
	if got := Shape("Москва"); got != "Xx" {
		t.Errorf("Shape(Москва) = %q", got)
	}
}
