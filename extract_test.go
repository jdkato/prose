package prose

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeNER(text string, model *Model) (*Document, error) {
	return NewDocument(text,
		WithSegmentation(false),
		UsingModel(model))
}

type prodigyOuput struct {
	Text   string
	Spans  []LabeledEntity
	Answer string
}

func readProdigy(jsonLines []byte) []prodigyOuput {
	dec := json.NewDecoder(bytes.NewReader(jsonLines))
	entries := []prodigyOuput{}
	for {
		ent := prodigyOuput{}
		err := dec.Decode(&ent)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		entries = append(entries, ent)
	}
	return entries
}

func split(data []prodigyOuput) ([]EntityContext, []prodigyOuput) {
	cutoff := int(float64(len(data)) * 0.8)

	train, test := []EntityContext{}, []prodigyOuput{}
	for i := range data {
		if i < cutoff {
			train = append(train, EntityContext{
				Text:   data[i].Text,
				Spans:  data[i].Spans,
				Accept: data[i].Answer == "accept"})
		} else {
			test = append(test, data[i])
		}
	}

	return train, test
}

func TestSumLogs(t *testing.T) {
	assert.Equal(t, 3.0, sumLogs([]float64{math.Log2(3), math.Log2(5)}))
}

func TestNERProdigy(t *testing.T) {
	data := filepath.Join(testdata, "reddit_product.jsonl")

	file, e := ioutil.ReadFile(data)
	if e != nil {
		panic(e)
	}

	train, test := split(readProdigy(file))
	correct := 0.0

	model := ModelFromData("PRODUCT", UsingEntities(train))
	for _, entry := range test {
		doc, _ := makeNER(entry.Text, model)
		ents := doc.Entities()
		if entry.Answer != "accept" && len(ents) == 0 {
			correct++
		} else {
			expected := []string{}
			for _, span := range entry.Spans {
				expected = append(expected, entry.Text[span.Start:span.End])
			}
			if reflect.DeepEqual(expected, ents) {
				correct++
			}
		}
	}
	assert.True(t, correct/float64(len(test)) >= 0.819444) // baseline
}
