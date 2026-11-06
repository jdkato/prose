// Command train-tagger trains a part-of-speech model for prose from CoNLL-U
// treebanks, such as those of Universal Dependencies.
//
// Usage:
//
//	go run ./cmd/train-tagger -lang de -train de_gsd-ud-train.conllu \
//	    -dev de_gsd-ud-dev.conllu -test de_gsd-ud-test.conllu \
//	    -source "UD German-GSD r2.15" -license "CC BY-SA 4.0" \
//	    -out tag/lang/de/de.bin
//
// Several training files may be given, separated by commas. The model is
// evaluated on the dev and test files, through the same tagger that will read
// it in production, so the reported accuracy is what a user gets.
//
// With -eval and no -train, an existing model is evaluated instead.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jdkato/prose/v3/internal/conllu"
	"github.com/jdkato/prose/v3/internal/flat"
	"github.com/jdkato/prose/v3/internal/perceptron"
	"github.com/jdkato/prose/v3/tag"
)

func main() {
	var (
		lang     = flag.String("lang", "", "BCP 47 language tag the model is for (required)")
		train    = flag.String("train", "", "comma-separated CoNLL-U training files")
		dev      = flag.String("dev", "", "comma-separated CoNLL-U development files, for reporting")
		test     = flag.String("test", "", "comma-separated CoNLL-U test files, for reporting")
		out      = flag.String("out", "", "where to write the model")
		eval     = flag.String("eval", "", "an existing model to evaluate instead of training")
		source   = flag.String("source", "", "description of the training data, stored in the model")
		license  = flag.String("license", "", "license of the training data, stored in the model")
		iters    = flag.Int("iter", 5, "training passes")
		prune    = flag.Float64("prune", 0.5, "drop averaged weights below this magnitude")
		minCount = flag.Int("min-count", 2, "drop features seen fewer times than this")
		seed     = flag.Int64("seed", 1, "shuffle seed")
		template = flag.Int("template", int(perceptron.TemplateRunes), "feature template")
		xpos     = flag.Bool("xpos", false, "train on the language-specific XPOS column rather than UPOS")
	)
	flag.Parse()

	if *lang == "" {
		fatal("-lang is required")
	}

	col := conllu.UPOS
	if *xpos {
		col = conllu.XPOS
	}

	var tagger *tag.Tagger
	var size int

	switch {
	case *eval != "":
		data, err := os.ReadFile(*eval)
		if err != nil {
			fatal(err)
		}
		tagger, err = tag.FromBytes(data)
		if err != nil {
			fatal(err)
		}
		size = len(data)

	case *train != "":
		sents := read(*train, col)
		fmt.Fprintf(os.Stderr, "training on %d sentences, %d tokens\n", len(sents), tokens(sents))

		tagset := tag.Universal
		if *xpos {
			tagset = tag.Penn
		}

		started := time.Now()
		src, stats, err := perceptron.Train(toPerceptron(sents), perceptron.Options{
			Template:   perceptron.Template(*template),
			Iterations: *iters,
			Seed:       *seed,
			Prune:      *prune,
			MinCount:   *minCount,
			Fixed:      fixedFor(tagset),
			Progress: func(pass int, acc float64) {
				fmt.Fprintf(os.Stderr, "  pass %d: %.2f%% on training data\n", pass, acc*100)
			},
		})
		if err != nil {
			fatal(err)
		}
		src.Meta[flat.MetaLanguage] = *lang
		src.Meta[flat.MetaTagset] = tagset
		if *source != "" {
			src.Meta[flat.MetaSource] = *source
		}
		if *license != "" {
			src.Meta[flat.MetaLicense] = *license
		}

		data, err := flat.Marshal(*src)
		if err != nil {
			fatal(err)
		}
		size = len(data)
		fmt.Fprintf(os.Stderr, "trained in %s: %d classes, %d features, %d weights, %d fixed words, %.1f MB\n",
			time.Since(started).Round(time.Second), stats.Classes, stats.Features, stats.Pairs, stats.TagMap,
			float64(size)/1e6)

		if *out != "" {
			if err := os.WriteFile(*out, data, 0o644); err != nil { //nolint:gosec // a model file is public data
				fatal(err)
			}
		}

		// Evaluate through the production reader, not the trainer's own
		// weights: that proves the file round-trips.
		tagger, err = tag.FromBytes(data)
		if err != nil {
			fatal(err)
		}

	default:
		fatal("give -train to train a model or -eval to evaluate one")
	}

	fmt.Printf("lang=%s tagset=%s size=%d\n", tagger.Language(), tagger.Tagset(), size)
	for _, set := range []struct{ name, files string }{{"dev", *dev}, {"test", *test}} {
		if set.files == "" {
			continue
		}
		sents := read(set.files, col)
		right, total := accuracy(tagger, sents)
		fmt.Printf("%s: %.2f%% (%d/%d tokens)\n", set.name, 100*float64(right)/float64(total), right, total)
	}
}

func read(files string, col conllu.Column) []conllu.Sentence {
	var all []conllu.Sentence
	for _, path := range strings.Split(files, ",") {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		sents, err := conllu.ReadFile(path, col)
		if err != nil {
			fatal(err)
		}
		all = append(all, sents...)
	}
	if len(all) == 0 {
		fatal("no sentences in " + files)
	}
	return all
}

func tokens(sents []conllu.Sentence) int {
	n := 0
	for _, s := range sents {
		n += len(s.Words)
	}
	return n
}

func toPerceptron(sents []conllu.Sentence) []perceptron.Sentence {
	out := make([]perceptron.Sentence, len(sents))
	for i, s := range sents {
		out[i] = perceptron.Sentence{Words: s.Words, Tags: s.Tags}
	}
	return out
}

// fixedFor mirrors the tagger's own short-circuits, so training sees the tag
// sequence tagging will.
func fixedFor(tagset string) func(string) (string, bool) {
	data, err := flat.Marshal(flat.Source{
		Weights: map[string]map[string]float64{"bias": {"X": 0}},
		Classes: []string{"X"},
		Meta:    map[string]string{flat.MetaTagset: tagset},
	})
	if err != nil {
		fatal(err)
	}
	probe, err := tag.FromBytes(data)
	if err != nil {
		fatal(err)
	}
	return func(word string) (string, bool) {
		toks := probe.Tag([]string{word})
		if toks[0].Tag == "X" {
			return "", false
		}
		return toks[0].Tag, true
	}
}

func accuracy(t *tag.Tagger, sents []conllu.Sentence) (right, total int) {
	for _, s := range sents {
		toks := t.Tag(s.Words)
		for i, tok := range toks {
			if tok.Tag == s.Tags[i] {
				right++
			}
			total++
		}
	}
	return right, total
}

func fatal(v any) {
	fmt.Fprintln(os.Stderr, "train-tagger:", v)
	os.Exit(1)
}
