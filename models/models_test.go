package models

import "testing"

func TestLookup(t *testing.T) {
	for lang, want := range map[string]Info{
		"de":    {Language: "de", Tagger: true, Segmenter: true},
		"pt-BR": {Language: "pt", Tagger: true, Segmenter: true},
		"pt_br": {Language: "pt", Tagger: true, Segmenter: true},
		"da":    {Language: "da", Tagger: true, Segmenter: true},
		"nb":    {Language: "nb", Tagger: true, Segmenter: true},
	} {
		got, ok := Lookup(lang)
		if !ok || got != want {
			t.Errorf("Lookup(%q) = %+v, %t; want %+v", lang, got, ok, want)
		}
	}
	for _, lang := range []string{"en", "ja", "cs", "ca", "xx", ""} {
		if _, ok := Lookup(lang); ok {
			t.Errorf("Lookup(%q) should find nothing", lang)
		}
	}
	if Archive("de") != "de.zip" {
		t.Errorf("Archive = %q", Archive("de"))
	}
}
