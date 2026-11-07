package models

import "testing"

func TestLookup(t *testing.T) {
	for lang, want := range map[string]Info{
		"de":    {Language: "de", Tagger: true, Segmenter: true},
		"pt-BR": {Language: "pt", Tagger: true, Segmenter: true},
		"pt_br": {Language: "pt", Tagger: true, Segmenter: true},
		"cs":    {Language: "cs", Tagger: false, Segmenter: true},
		"ca":    {Language: "ca", Tagger: true, Segmenter: false},
	} {
		got, ok := Lookup(lang)
		if !ok || got != want {
			t.Errorf("Lookup(%q) = %+v, %t; want %+v", lang, got, ok, want)
		}
	}
	for _, lang := range []string{"en", "ja", "xx", ""} {
		if _, ok := Lookup(lang); ok {
			t.Errorf("Lookup(%q) should find nothing", lang)
		}
	}
	if Archive("de") != "de.zip" {
		t.Errorf("Archive = %q", Archive("de"))
	}
}
