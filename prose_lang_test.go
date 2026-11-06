package prose_test

import (
	"testing"

	"github.com/jdkato/prose/v3"
	_ "github.com/jdkato/prose/v3/segment/lang/de"
	_ "github.com/jdkato/prose/v3/tag/lang/de"
)

func TestDocumentWithLanguage(t *testing.T) {
	src := "Dr. Müller fährt z.B. nach Berlin. Er bleibt dort."
	doc, err := prose.NewDocument(src, prose.WithLanguage("de"))
	if err != nil {
		t.Fatal(err)
	}
	if n := len(doc.Sentences()); n != 2 {
		t.Errorf("got %d sentences, want 2: %+v", n, doc.Sentences())
	}
	if doc.Entities() != nil {
		t.Error("entity extraction is English only and should have been skipped")
	}
	var saw bool
	for _, tok := range doc.Tokens() {
		if tok.In(src) != tok.Text {
			t.Errorf("token %q at %d does not locate its text", tok.Text, tok.Start)
		}
		if tok.Text == "Berlin" && tok.Tag == "PROPN" {
			saw = true
		}
	}
	if !saw {
		t.Errorf("Berlin should be a PROPN: %+v", doc.Tokens())
	}

	if _, err := prose.NewDocument("こんにちは。", prose.WithLanguage("ja")); err == nil {
		t.Error("a language without a tagging model should be an error when tagging is on")
	}
	doc, err = prose.NewDocument("こんにちは。元気ですか？", prose.WithLanguage("ja"), prose.WithTagging(false))
	if err != nil {
		t.Fatal(err)
	}
	if n := len(doc.Sentences()); n != 2 {
		t.Errorf("got %d Japanese sentences, want 2", n)
	}
}
