package main

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The same models make the same archive, entry for entry, so a release can
// be rebuilt and compared.
func TestWriteZipIsDeterministic(t *testing.T) {
	files := map[string][]byte{"tagger.bin": []byte("model"), "model.json": []byte("{}\n")}
	dir := t.TempDir()

	var archives [][]byte
	for i := range 2 {
		name := filepath.Join(dir, "de"+string(rune('0'+i))+".zip")
		if err := writeZip(name, "de/styles/config/models/de", files); err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		archives = append(archives, b)
	}
	if !bytes.Equal(archives[0], archives[1]) {
		t.Error("two runs wrote different archives")
	}

	r, err := zip.NewReader(bytes.NewReader(archives[0]), int64(len(archives[0])))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, f := range r.File {
		names = append(names, f.Name)
	}
	if want := "de/styles/config/models/de/model.json de/styles/config/models/de/tagger.bin"; strings.Join(names, " ") != want {
		t.Errorf("entries = %q, want %q", strings.Join(names, " "), want)
	}
}

func TestPublishedSource(t *testing.T) {
	src := string(publishedSource([]Entry{
		{Language: "cs", Segmenter: &Model{}},
		{Language: "de", Tagger: &Model{}, Segmenter: &Model{}},
	}))
	for _, want := range []string{
		"package models",
		`{Language: "cs", Tagger: false, Segmenter: true},`,
		`{Language: "de", Tagger: true, Segmenter: true},`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated source lacks %q:\n%s", want, src)
		}
	}
}
