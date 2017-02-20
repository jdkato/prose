package tag

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/jdkato/prose/tokenize"
	"github.com/jdkato/prose/util"
)

var none = regexp.MustCompile(`^(?:0|\*[\w?]\*|\*\-\d{1,3}|\*[A-Z]+\*\-\d{1,3}|\*)$`)
var keep = regexp.MustCompile(`^\-[A-Z]{3}\-$`)

// AveragedPerceptron ...
type AveragedPerceptron struct {
	Classes   []string
	Instances int
	Stamps    map[string]int
	TagMap    map[string]string
	Totals    map[string]int
	Weights   map[string]map[string]float64
}

// NewAveragedPerceptron ...
func NewAveragedPerceptron() *AveragedPerceptron {
	var ap AveragedPerceptron
	var err error

	ap.Totals = make(map[string]int)
	ap.Stamps = make(map[string]int)
	err = json.Unmarshal(util.GetAsset("classes.json"), &ap.Classes)
	util.CheckError(err)
	err = json.Unmarshal(util.GetAsset("tags.json"), &ap.TagMap)
	util.CheckError(err)
	err = json.Unmarshal(util.GetAsset("weights.json"), &ap.Weights)
	util.CheckError(err)
	return &ap
}

// Predict ...
func (ap AveragedPerceptron) Predict(features map[string]float64) string {
	var weights map[string]float64
	var found bool

	scores := make(map[string]float64)
	for feat, value := range features {
		if weights, found = ap.Weights[feat]; !found || value == 0 {
			continue
		}
		for label, weight := range weights {
			if _, ok := scores[label]; ok {
				scores[label] += value * weight
			} else {
				scores[label] = value * weight
			}
		}
	}
	return max(scores)
}

// PerceptronTagger ...
type PerceptronTagger struct {
	Model *AveragedPerceptron
}

// NewPerceptronTagger ...
func NewPerceptronTagger() *PerceptronTagger {
	var pt PerceptronTagger
	pt.Model = NewAveragedPerceptron()
	return &pt
}

// Tag ...
func (pt PerceptronTagger) Tag(words []string) []tokenize.Token {
	var tokens []tokenize.Token
	var clean []string
	var tag string
	var found bool

	p1, p2 := "-START-", "-START2-"
	context := []string{p1, p2}
	for _, w := range words {
		if w == "" {
			continue
		}
		context = append(context, normalize(w))
		clean = append(clean, w)
	}
	context = append(context, []string{"-END-", "-END2-"}...)
	for i, word := range clean {
		if none.MatchString(word) {
			tag = "-NONE-"
		} else if keep.MatchString(word) {
			tag = word
		} else if tag, found = pt.Model.TagMap[word]; !found {
			tag = pt.Model.Predict(featurize(i, word, context, p1, p2))
		}
		tokens = append(tokens, tokenize.Token{Tag: tag, Text: word})
		p2 = p1
		p1 = tag
	}

	return tokens
}

// TokenizeAndTag ...
func (pt PerceptronTagger) TokenizeAndTag(corpus string) []tokenize.Token {
	var tokens []tokenize.Token
	for _, s := range tokenize.SentenceTokenizer(corpus) {
		tokens = append(tokens, pt.Tag(tokenize.WordTokenizer(s.Text))...)
	}
	return tokens
}

func max(scores map[string]float64) string {
	var class string
	max := 0.0
	for label, value := range scores {
		if value > max {
			max = value
			class = label
		}
	}
	return class
}

func featurize(i int, w string, ctx []string, p1 string, p2 string) map[string]float64 {
	feats := make(map[string]float64)
	suf := util.Min(len(w), 3)
	i = util.Min(len(ctx)-2, i+2)
	iminus := util.Min(len(ctx[i-1]), 3)
	iplus := util.Min(len(ctx[i+1]), 3)
	feats = add([]string{"bias"}, feats)
	feats = add([]string{"i suffix", w[len(w)-suf:]}, feats)
	feats = add([]string{"i pref1", string(w[0])}, feats)
	feats = add([]string{"i-1 tag", p1}, feats)
	feats = add([]string{"i-2 tag", p2}, feats)
	feats = add([]string{"i tag+i-2 tag", p1, p2}, feats)
	feats = add([]string{"i word", ctx[i]}, feats)
	feats = add([]string{"i-1 tag+i word", p1, ctx[i]}, feats)
	feats = add([]string{"i-1 word", ctx[i-1]}, feats)
	feats = add([]string{"i-1 suffix", ctx[i-1][len(ctx[i-1])-iminus:]}, feats)
	feats = add([]string{"i-2 word", ctx[i-2]}, feats)
	feats = add([]string{"i+1 word", ctx[i+1]}, feats)
	feats = add([]string{"i+1 suffix", ctx[i+1][len(ctx[i+1])-iplus:]}, feats)
	feats = add([]string{"i+2 word", ctx[i+2]}, feats)
	return feats
}

func add(args []string, features map[string]float64) map[string]float64 {
	key := strings.Join(args, " ")
	if _, ok := features[key]; ok {
		features[key]++
	} else {
		features[key] = 1
	}
	return features
}

func normalize(word string) string {
	if word == "" {
		return word
	}
	first := string(word[0])
	if strings.Contains(word, "-") && first != "-" {
		return "!HYPHEN"
	} else if _, err := strconv.Atoi(word); err == nil && len(word) == 4 {
		return "!YEAR"
	} else if _, err := strconv.Atoi(first); err == nil {
		return "!DIGITS"
	}
	return strings.ToLower(word)
}
