package tag

import "testing"

// The closed-class lists hold words, so only a word of letters is checked
// against them; a token with punctuation in it is left to the model.
func TestAllLetters(t *testing.T) {
	for word, want := range map[string]bool{
		"write": true, "Straße": true, "'It": false, "It--it": false, "co-op": false, "": false,
	} {
		if got := allLetters(word); got != want {
			t.Errorf("allLetters(%q) = %v, want %v", word, got, want)
		}
	}
}

func TestParseClosed(t *testing.T) {
	got := parseClosed("PDT all such;XX nope;DT", []string{"NN", "PDT", "DT"})
	if len(got) != 1 || !got[1]["all"] || !got[1]["such"] || got[1]["write"] {
		t.Errorf("parseClosed = %v, want PDT (class 1) with all and such", got)
	}
	if parseClosed("", []string{"NN"}) != nil {
		t.Error("an empty list should parse to nil")
	}
}
