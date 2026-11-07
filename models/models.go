// Package models says which languages prose publishes model archives for,
// without linking any of them.
//
// The archives are written by prose/v3/cmd/pack-models and attached to a
// prose release, one per language, for a program that downloads a model when
// a configuration first names its language rather than carrying every model
// in its binary. A program that would rather link them imports the packages
// under prose/v3/tag/lang and prose/v3/segment/lang instead.
package models

import "strings"

// Info is what a language's archive holds.
type Info struct {
	// Language is the BCP 47 tag the archive is named by.
	Language string
	// Tagger reports whether the archive holds a part-of-speech model.
	Tagger bool
	// Segmenter reports whether it holds a sentence-segmentation model. A
	// language without one is segmented by punctuation rules.
	Segmenter bool
}

// Lookup returns what is published for lang, a BCP 47 tag; a tag with a
// region ("pt-BR") finds its base language's archive.
func Lookup(lang string) (Info, bool) {
	lang = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(lang), "_", "-"))
	base, _, _ := strings.Cut(lang, "-")
	for _, info := range Published {
		if info.Language == lang || info.Language == base {
			return info, true
		}
	}
	return Info{}, false
}

// Archive is the file name of a language's archive in a release.
func Archive(lang string) string { return lang + ".zip" }
