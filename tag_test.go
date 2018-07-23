package prose

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var wsj = "Pierre|NNP Vinken|NNP ,|, 61|CD years|NNS old|JJ ,|, will|MD " +
	"join|VB the|DT board|NN as|IN a|DT nonexecutive|JJ director|NN " +
	"Nov.|NNP 29|CD .|.\nMr.|NNP Vinken|NNP is|VBZ chairman|NN of|IN " +
	"Elsevier|NNP N.V.|NNP ,|, the|DT Dutch|NNP publishing|VBG " +
	"group|NN .|. Rudolph|NNP Agnew|NNP ,|, 55|CD years|NNS old|JJ " +
	"and|CC former|JJ chairman|NN of|IN Consolidated|NNP Gold|NNP " +
	"Fields|NNP PLC|NNP ,|, was|VBD named|VBN a|DT nonexecutive|JJ " +
	"director|NN of|IN this|DT British|JJ industrial|JJ conglomerate|NN " +
	".|.\nA|DT form|NN of|IN asbestos|NN once|RB used|VBN to|TO make|VB " +
	"Kent|NNP cigarette|NN filters|NNS has|VBZ caused|VBN a|DT high|JJ " +
	"percentage|NN of|IN cancer|NN deaths|NNS among|IN a|DT group|NN " +
	"of|IN workers|NNS exposed|VBN to|TO it|PRP more|RBR than|IN " +
	"30|CD years|NNS ago|IN ,|, researchers|NNS reported|VBD .|."

func makeTagger(text string) (*Document, error) {
	return NewDocument(
		text,
		WithSegmentation(false),
		WithExtraction(false))
}

func ExampleReadTagged() {
	tagged := "Pierre|NNP Vinken|NNP ,|, 61|CD years|NNS"
	fmt.Println(ReadTagged(tagged, "|"))
	// Output: [[[Pierre Vinken , 61 years] [NNP NNP , CD NNS]]]
}

func TestTagSimple(t *testing.T) {
	doc, err := makeTagger("Pierre Vinken, 61 years old, will join the board as a nonexecutive director Nov. 29.")
	if err != nil {
		panic(err)
	}
	tags := []string{}
	for _, tok := range doc.Tokens() {
		tags = append(tags, tok.Tag)
	}
	assert.Equal(t, []string{
		"NNP", "NNP", ",", "CD", "NNS", "JJ", ",", "MD", "VB", "DT", "NN",
		"IN", "DT", "JJ", "NN", "NNP", "CD", "."}, tags)
}

func TestTagTreebank(t *testing.T) {
	tagger := newPerceptronTagger()
	tokens, expected := []*Token{}, []string{}

	tags := readDataFile(filepath.Join(testdata, "treebank_tags.json"))
	checkError(json.Unmarshal(tags, &expected))

	treebank := readDataFile(filepath.Join(testdata, "treebank_tokens.json"))
	checkError(json.Unmarshal(treebank, &tokens))

	correct := 0.0
	for i, tok := range tagger.tag(tokens) {
		if expected[i] == tok.Tag {
			correct++
		}
	}
	assert.True(t, correct/float64(len(expected)) >= 0.957477) // baseline
}

func BenchmarkTag(b *testing.B) {
	tagger := newPerceptronTagger()
	tokens := []*Token{}

	treebank := readDataFile(filepath.Join(testdata, "treebank_tokens.json"))
	checkError(json.Unmarshal(treebank, &tokens))
	for n := 0; n < b.N; n++ {
		_ = tagger.tag(tokens)
	}
}

/* TODO: POS training API

func TestTrain(t *testing.T) {
	sentences := ReadTagged(wsj, "|")
	iter := random(5, 20)
	tagger.Train(sentences, iter)

	tagSet := []string{}
	nrWords := 0
	for _, tuple := range sentences {
		nrWords += len(tuple[0])
		for _, tag := range tuple[1] {
			if !util.StringInSlice(tag, tagSet) {
				tagSet = append(tagSet, tag)
			}
		}
	}

	assert.Equal(t, nrWords*iter, int(tagger.model.instances))
	assert.Subset(t, tagger.Classes(), tagSet)
}*/

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
