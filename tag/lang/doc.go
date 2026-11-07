// Package lang holds prose's part-of-speech models for languages other than
// English, one package each.
//
// A model is registered with prose/v3/tag when its package is imported, and
// nothing is linked that was not imported:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/de" // German only
//	import _ "github.com/jdkato/prose/v3/tag/lang/all" // every language
//
// Every model emits the Universal Dependencies tag set, so a rule written as
// "a NOUN followed by a VERB" reads the same in each language. English is the
// exception: its model is older, larger and Penn Treebank-tagged, and
// tag.ToUniversal maps its tags across.
//
// A program that would rather download a model than link it reads
// prose/v3/models for what is published and fetches the archive
// prose/v3/cmd/pack-models writes from the matching prose release.
//
// MODELS.md lists each model's training data, license and accuracy; build.sh
// retrains them all. The models are averaged perceptrons trained by
// prose/v3/cmd/train-tagger, and a caller with a treebank of their own can
// train one the same way and load it with tag.FromBytes or register it with
// tag.RegisterLanguage.
package lang
