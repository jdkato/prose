package main

import (
	"encoding/json"
	"testing"
)

func TestAnalyze(t *testing.T) {
	src := "Marie Curie studied in Paris. She won the Nobel Prize twice."
	got, err := analyze(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(got.Sentences) != 2 {
		t.Fatalf("sentences = %d, want 2", len(got.Sentences))
	}
	for _, s := range got.Sentences {
		if src[s.Start:s.End] != s.Text {
			t.Errorf("sentence %q is not at [%d:%d]", s.Text, s.Start, s.End)
		}
	}
	for _, tok := range got.Tokens {
		if src[tok.Start:tok.End] != tok.Text {
			t.Errorf("token %q is not at [%d:%d]", tok.Text, tok.Start, tok.End)
		}
		if tok.Tag == "" {
			t.Errorf("token %q has no tag", tok.Text)
		}
	}

	var ents []string
	for _, e := range got.Entities {
		if src[e.Start:e.End] != e.Text {
			t.Errorf("entity %q is not at [%d:%d]", e.Text, e.Start, e.End)
		}
		ents = append(ents, e.Text+" "+e.Label)
	}
	want := []string{"Marie Curie PERSON", "Paris GPE"}
	if len(ents) < 2 || ents[0] != want[0] || ents[1] != want[1] {
		t.Errorf("entities = %v, want prefix %v", ents, want)
	}

	if got.Readability.Words == 0 || got.Readability.Sentences != 2 {
		t.Errorf("readability counts = %+v", got.Readability)
	}

	// The page reads the JSON, so its field names are part of the contract.
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"tokens", "sentences", "entities", "readability"} {
		if _, ok := m[key]; !ok {
			t.Errorf("JSON is missing %q", key)
		}
	}
}

func TestConvertCase(t *testing.T) {
	got := convertCase("the quick brown fox jumps over the lazy dog")
	if got.AP != "The Quick Brown Fox Jumps Over the Lazy Dog" {
		t.Errorf("AP = %q", got.AP)
	}
	if got.Chicago == "" || got.Sentence == "" {
		t.Errorf("cases = %+v", got)
	}
}
