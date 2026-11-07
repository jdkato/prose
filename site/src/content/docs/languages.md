---
title: Languages
summary: English is built in; other languages are packages you import.
order: 3
---

# Languages

English is built in. Other languages are packages you import, so a binary
carries only the languages it uses:

```go
import (
    "github.com/jdkato/prose/v3"

    _ "github.com/jdkato/prose/v3/segment/lang/de"
    _ "github.com/jdkato/prose/v3/tag/lang/de"
)

doc, err := prose.NewDocument("Dr. Müller fährt z.B. nach Berlin. Er bleibt dort.",
    prose.WithLanguage("de"))

for _, tok := range doc.Tokens() {
    fmt.Println(tok.Text, tok.Tag)
    // Dr. NOUN
    // Müller PROPN
    // fährt VERB
    // ...
}
```

Each piece can be asked for a language on its own: `tokenize.ForLanguage`,
`segment.ForLanguage` and `tag.ForLanguage`. Import `segment/lang/all` and
`tag/lang/all` when the language is decided at run time.

| | Languages |
| --- | --- |
| Tagging | Danish, Dutch, English, Finnish, French, German, Italian, Norwegian, Portuguese, Spanish, Swedish, Turkish |
| Segmentation | The same twelve; punctuation rules for everything else, including Chinese, Japanese and Korean |
| Tokenization | Every language. French and Italian elisions (`l'homme`), French clitics (`donne-moi`), Spanish `¿¡`, every quotation style, and script boundaries for Chinese and Japanese |

## Tag sets

The English tagger uses the Penn Treebank tag set; every other model emits the
[Universal Dependencies](https://universaldependencies.org/u/pos/) tag set, so
one rule reads across languages, and `tag.ToUniversal` maps Penn tags onto it.

The models are averaged perceptrons trained on the Universal Dependencies
treebanks. The [models page](../models) lists each one's training data,
license and accuracy, and `cmd/train-tagger` trains one from any CoNLL-U
treebank.

Entity extraction and readability are English only.
