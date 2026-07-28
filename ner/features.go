// This file holds the named-entity feature set: the shape and part-of-speech
// simplifications, the IOB chunk coalescing, and the offset adjustment.

package ner

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jdkato/prose/v3/internal/wordlist"
	"github.com/jdkato/prose/v3/token"
)

var featureOrder = []string{
	"bias", "en-wordlist", "nextpos", "nextword", "pos", "pos+prevtag",
	"prefix3", "prevpos", "prevtag", "prevword", "shape", "shape+prevtag",
	"suffix3", "word", "word+nextpos", "word.lower", "wordlen"}

// binaryMaxentClassifier is a feature encoding that generates vectors
// containing binary joint-features of the form:
//
//    |  joint_feat(fs, l) = { 1 if (fs[fname] == fval) and (l == label)
//    |                      {
//    |                      { 0 otherwise
//
// where `fname` is the name of an input-feature, `fval` is a value for that
// input-feature, and `label` is a label.
//
// See https://www.nltk.org/_modules/nltk/classify/maxent.html for more
// information.

func parseEntities(ents []string) string {
	if stringInSlice("B-PERSON", ents) && len(ents) == 2 {
		// PERSON takes precedence because it's hard to identify.
		return "PERSON"
	}
	return strings.Split(ents[0], "-")[1]
}

// coalesce turns a run of tokens into one entity.
//
// The entity's text is the exact source span rather than the tokens rejoined
// with single spaces: Text must be what is actually at Start, so "New\nYork"
// in the source cannot come back as "New York".
func coalesce(src string, parts []token.Token) Entity {
	labels := make([]string, len(parts))
	for i, tok := range parts {
		labels[i] = tok.Label
	}

	start := parts[0].Start
	end := parts[len(parts)-1].End()
	if start < 0 || end > len(src) || start > end {
		// Should not happen, but never hand back an entity whose offsets do
		// not locate its text.
		return Entity{Label: parseEntities(labels), Text: "", Start: 0}
	}

	return Entity{
		Label: parseEntities(labels),
		Text:  src[start:end],
		Start: start,
	}
}

// extractInto fills feats with the model's features for token i.
//
// Writing into a caller-owned map lets the Extracter pool it rather than
// allocating one per token.
func extractInto(feats map[string]string, i int, ctx []token.Token, history []string) {

	word := ctx[i].Text
	prevShape := "None"

	feats["bias"] = "True"
	feats["word"] = word
	feats["pos"] = ctx[i].Tag
	feats["en-wordlist"] = boolFeature(wordlist.IsBasic(word))
	feats["word.lower"] = strings.ToLower(word)
	feats["suffix3"] = nSuffix(word, 3)
	feats["prefix3"] = nPrefix(word, 3)
	feats["shape"] = shape(word)
	feats["wordlen"] = strconv.Itoa(len(word))

	switch i {
	case 0:
		feats["prevtag"] = "None"
		feats["prevword"], feats["prevpos"] = "None", "None"
	case 1:
		feats["prevword"] = strings.ToLower(ctx[i-1].Text)
		feats["prevpos"] = ctx[i-1].Tag
		feats["prevtag"] = history[i-1]
	default:
		feats["prevword"] = strings.ToLower(ctx[i-1].Text)
		feats["prevpos"] = ctx[i-1].Tag
		feats["prevtag"] = history[i-1]
		prevShape = shape(ctx[i-1].Text)
	}

	if i == len(ctx)-1 {
		feats["nextword"], feats["nextpos"] = "None", "None"
	} else {
		feats["nextword"] = strings.ToLower(ctx[i+1].Text)
		feats["nextpos"] = strings.ToLower(ctx[i+1].Tag)
	}

	feats["word+nextpos"] = strings.Join(
		[]string{feats["word.lower"], feats["nextpos"]}, "+")
	feats["pos+prevtag"] = strings.Join(
		[]string{feats["pos"], feats["prevtag"]}, "+")
	feats["shape+prevtag"] = strings.Join(
		[]string{prevShape, feats["prevtag"]}, "+")
}

func shape(word string) string {
	if isNumeric(word) {
		return "number"
	} else if match, _ := regexp.MatchString(`\W+$`, word); match {
		return "punct"
	} else if match, _ := regexp.MatchString(`\w+$`, word); match {
		if strings.ToLower(word) == word {
			return "downcase"
		}
		// strings.Title is deprecated for applying title-casing rules that
		// mishandle Unicode punctuation. That behaviour is load-bearing here:
		// the model was trained against these exact shape features, and the
		// obvious replacements disagree on acronyms and hyphenated or
		// apostrophised names — "NASA", "U.S.A.", "McDonald" and "Jean-Luc"
		// all read as upcase under Title and as mixedcase without it.
		//nolint:staticcheck // deliberate: preserves trained feature values
		if strings.Title(word) == word {
			return "upcase"
		}
		return "mixedcase"
	}
	return "other"
}

func simplePOS(pos string) string {
	if strings.HasPrefix(pos, "V") {
		return "v"
	}
	return strings.Split(pos, "-")[0]
}

// --- small helpers ---------------------------------------------------------

func boolFeature(b bool) string {
	if b {
		return "True"
	}
	return "False"
}

func nSuffix(word string, length int) string {
	return strings.ToLower(word[len(word)-minInt(len(word), length):])
}

func nPrefix(word string, length int) string {
	return strings.ToLower(word[:minInt(len(word), length)])
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func stringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if b == a {
			return true
		}
	}
	return false
}
