// Command pack-models writes prose's language models as archives a program
// downloads at run time instead of linking: one per language, in the layout
// Vale installs as a package.
//
//	go run ./cmd/pack-models -version v3.3.0 -out dist/models
//
// Each <lang>.zip holds <lang>/styles/config/models/<lang>/ with:
//
//   - tagger.bin, the part-of-speech model, where the language has one;
//   - segmenter.json.gz and abbreviations.txt, the Punkt model and the
//     abbreviations added to it, where the language has one;
//   - model.json, describing both: the prose version they were packed from,
//     and each model's training data and license.
//
// manifest.json beside the archives lists every language and what it has, and
// models/published.go in the module is rewritten to say the same, so that a
// program can tell a language with no model from one it has not downloaded
// without asking the network.
package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jdkato/prose/v3/segment"
	_ "github.com/jdkato/prose/v3/segment/lang/all"
	"github.com/jdkato/prose/v3/tag"
	_ "github.com/jdkato/prose/v3/tag/lang/all"
)

// Entry describes one language's models, in model.json and the manifest.
type Entry struct {
	Language  string `json:"language"`
	Prose     string `json:"prose"`
	Tagger    *Model `json:"tagger,omitempty"`
	Segmenter *Model `json:"segmenter,omitempty"`
}

// Model describes one model's provenance.
type Model struct {
	Source  string `json:"source"`
	License string `json:"license"`
	Tagset  string `json:"tagset,omitempty"`
}

// punkt is where every shipped Punkt model comes from.
var punkt = Model{Source: "gopkg.in/neurosnap/sentences", License: "MIT"}

func main() {
	var (
		version   = flag.String("version", "", "prose version the models are packed from (required)")
		out       = flag.String("out", "dist/models", "directory to write the archives to")
		published = flag.String("published", "models/published.go", "Go file listing the languages, or empty to skip")
	)
	flag.Parse()
	if *version == "" {
		fatal("-version is required")
	}

	if err := os.MkdirAll(*out, 0o755); err != nil { //nolint:gosec // release files are public
		fatal(err)
	}

	var entries []Entry
	for _, lang := range languages() {
		entry := Entry{Language: lang, Prose: *version}
		files := map[string][]byte{}

		if data, ok := tag.LanguageModel(lang); ok {
			t, err := tag.FromBytes(data)
			if err != nil {
				fatal(fmt.Errorf("%s: %w", lang, err))
			}
			entry.Tagger = &Model{
				Source:  t.Meta(tag.MetaSource),
				License: t.Meta(tag.MetaLicense),
				Tagset:  t.Tagset(),
			}
			files["tagger.bin"] = data
		}
		if gz, abbrevs, ok := segment.LanguageModel(lang); ok {
			m := punkt
			entry.Segmenter = &m
			files["segmenter.json.gz"] = gz
			files["abbreviations.txt"] = []byte(strings.Join(abbrevs, "\n") + "\n")
		}

		meta, err := json.MarshalIndent(entry, "", "  ")
		if err != nil {
			fatal(err)
		}
		files["model.json"] = append(meta, '\n')

		name := filepath.Join(*out, lang+".zip")
		if err := writeZip(name, path.Join(lang, "styles", "config", "models", lang), files); err != nil {
			fatal(err)
		}
		entries = append(entries, entry)
		fmt.Printf("%s: %s\n", name, describe(entry))
	}

	manifest, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fatal(err)
	}
	if err := os.WriteFile(filepath.Join(*out, "manifest.json"), append(manifest, '\n'), 0o644); err != nil { //nolint:gosec // public data
		fatal(err)
	}

	if *published != "" {
		if err := os.WriteFile(*published, publishedSource(entries), 0o644); err != nil { //nolint:gosec // source code
			fatal(err)
		}
	}
}

// languages is every language either registry knows, sorted, English aside:
// its models are built into prose.
func languages() []string {
	seen := map[string]bool{}
	for _, l := range tag.Languages() {
		seen[l] = true
	}
	for _, l := range segment.Languages() {
		seen[l] = true
	}
	delete(seen, "en")

	out := make([]string, 0, len(seen))
	for l := range seen {
		out = append(out, l)
	}
	sort.Strings(out)
	return out
}

func describe(e Entry) string {
	var parts []string
	if e.Tagger != nil {
		parts = append(parts, "tagger")
	}
	if e.Segmenter != nil {
		parts = append(parts, "segmenter")
	}
	return strings.Join(parts, ", ")
}

// writeZip writes files under dir in a new archive at name. Entries are
// sorted and carry a fixed time, so the same models make the same bytes.
func writeZip(name, dir string, files map[string][]byte) error {
	names := make([]string, 0, len(files))
	for n := range files {
		names = append(names, n)
	}
	sort.Strings(names)

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, n := range names {
		f, err := w.CreateHeader(&zip.FileHeader{
			Name:     path.Join(dir, n),
			Method:   zip.Deflate,
			Modified: time.Unix(0, 0).UTC(),
		})
		if err != nil {
			return err
		}
		if _, err := f.Write(files[n]); err != nil {
			return err
		}
	}
	if err := w.Close(); err != nil {
		return err
	}
	return os.WriteFile(name, buf.Bytes(), 0o644) //nolint:gosec // release files are public
}

// publishedSource is the Go source of the models package: the list of
// languages the archives exist for.
func publishedSource(entries []Entry) []byte {
	var b strings.Builder
	b.WriteString(`// Code generated by cmd/pack-models; DO NOT EDIT.

package models

// Published lists every language prose publishes a model archive for, and
// what the archive holds. English is not among them: its models are built
// into prose.
var Published = []Info{
`)
	for _, e := range entries {
		fmt.Fprintf(&b, "\t{Language: %q, Tagger: %t, Segmenter: %t},\n", e.Language, e.Tagger != nil, e.Segmenter != nil)
	}
	b.WriteString("}\n")
	return []byte(b.String())
}

func fatal(v any) {
	fmt.Fprintln(os.Stderr, "pack-models:", v)
	os.Exit(1)
}
