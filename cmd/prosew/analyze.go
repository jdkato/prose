// Command prosew is prose compiled to WebAssembly for the playground on the
// documentation site. It installs ProseAnalyze and ProseCase on the
// JavaScript global object; see main.go. The work itself is here, in plain
// Go, so the host can test it.
package main

import (
	"github.com/jdkato/prose/v3"
	"github.com/jdkato/prose/v3/strcase"
	"github.com/jdkato/prose/v3/summarize"
)

// Analysis is everything the playground shows for one text.
type Analysis struct {
	Tokens      []Token     `json:"tokens"`
	Sentences   []Span      `json:"sentences"`
	Entities    []Entity    `json:"entities"`
	Readability Readability `json:"readability"`
}

// Token is a prose.Token with its end offset spelled out, so the page never
// has to add lengths itself.
type Token struct {
	Text  string `json:"text"`
	Tag   string `json:"tag"`
	Label string `json:"label,omitempty"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

// Span is a sentence: its text and where it sits in the source.
type Span struct {
	Text  string `json:"text"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

// Entity is a named entity and where it sits in the source.
type Entity struct {
	Text  string `json:"text"`
	Label string `json:"label"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

// Readability carries the counts and the scores from summarize.
type Readability struct {
	Words      int `json:"words"`
	Sentences  int `json:"sentences"`
	Syllables  int `json:"syllables"`
	Characters int `json:"characters"`

	FleschKincaid        float64 `json:"fleschKincaid"`
	FleschReadingEase    float64 `json:"fleschReadingEase"`
	GunningFog           float64 `json:"gunningFog"`
	SMOG                 float64 `json:"smog"`
	ColemanLiau          float64 `json:"colemanLiau"`
	AutomatedReadability float64 `json:"automatedReadability"`
	DaleChall            float64 `json:"daleChall"`
	LIX                  float64 `json:"lix"`
}

// Cases is one string in each case convention strcase offers.
type Cases struct {
	AP       string `json:"ap"`
	Chicago  string `json:"chicago"`
	Sentence string `json:"sentence"`
}

// analyze runs the whole English pipeline over text.
func analyze(text string) (*Analysis, error) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		return nil, err
	}

	out := &Analysis{
		Tokens:    make([]Token, 0, len(doc.Tokens())),
		Sentences: make([]Span, 0, len(doc.Sentences())),
		Entities:  make([]Entity, 0, len(doc.Entities())),
	}
	for _, t := range doc.Tokens() {
		out.Tokens = append(out.Tokens, Token{
			Text: t.Text, Tag: t.Tag, Label: t.Label, Start: t.Start, End: t.End(),
		})
	}
	for _, s := range doc.Sentences() {
		out.Sentences = append(out.Sentences, Span{Text: s.Text, Start: s.Start, End: s.End()})
	}
	for _, e := range doc.Entities() {
		out.Entities = append(out.Entities, Entity{
			Text: e.Text, Label: e.Label, Start: e.Start, End: e.End(),
		})
	}

	sd := summarize.NewDocument(text)
	out.Readability = Readability{
		Words:                int(sd.NumWords),
		Sentences:            int(sd.NumSentences),
		Syllables:            int(sd.NumSyllables),
		Characters:           int(sd.NumCharacters),
		FleschKincaid:        sd.FleschKincaid(),
		FleschReadingEase:    sd.FleschReadingEase(),
		GunningFog:           sd.GunningFog(),
		SMOG:                 sd.SMOG(),
		ColemanLiau:          sd.ColemanLiau(),
		AutomatedReadability: sd.AutomatedReadability(),
		DaleChall:            sd.DaleChall(),
		LIX:                  sd.LIX(),
	}
	return out, nil
}

// convertCase writes text in each convention.
func convertCase(text string) Cases {
	return Cases{
		AP:       strcase.NewTitleConverter(strcase.APStyle).Convert(text),
		Chicago:  strcase.NewTitleConverter(strcase.ChicagoStyle).Convert(text),
		Sentence: strcase.NewSentenceConverter().Convert(text),
	}
}
