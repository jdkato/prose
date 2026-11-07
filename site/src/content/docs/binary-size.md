---
title: Binary size
summary: Import only the models you use and the linker leaves the rest out.
order: 4
---

# Pay only for what you use

`prose` ships two trained English models: a part-of-speech tagger (~3 MB) and
an entity recognizer (~4 MB). Each lives in its own package, so the Go linker
leaves out whatever you do not import.

Importing the root package links both. If you only need tags, import the tagger
directly and the entity model never enters your binary, a 4.5 MB difference:

```go
import "github.com/jdkato/prose/v3/tag"

tagger, err := tag.New()
if err != nil {
    log.Fatal(err)
}

for _, tok := range tagger.Tag([]string{"Kubernetes", "orchestrates", "containers"}) {
    fmt.Println(tok.Text, tok.Tag)
}
```

The same applies to languages: each one's tagger and segmenter is its own
package, between 0.6 and 1.8 MB, linked only when imported.

Models load once per process and are shared, so building many documents does
not repeat the work.

## Packages

Each component stands alone if you want just that piece:

| Package | What it does |
| --- | --- |
| [`prose`](https://pkg.go.dev/github.com/jdkato/prose/v3) | `Document`, the whole pipeline |
| [`tokenize`](https://pkg.go.dev/github.com/jdkato/prose/v3/tokenize) | words and punctuation |
| [`segment`](https://pkg.go.dev/github.com/jdkato/prose/v3/segment) | sentences |
| [`tag`](https://pkg.go.dev/github.com/jdkato/prose/v3/tag) | part-of-speech tags |
| [`tag/lang`](https://pkg.go.dev/github.com/jdkato/prose/v3/tag/lang) | tagging models, one package per language |
| [`segment/lang`](https://pkg.go.dev/github.com/jdkato/prose/v3/segment/lang) | segmentation models, one package per language |
| [`ner`](https://pkg.go.dev/github.com/jdkato/prose/v3/ner) | named entities |
| [`token`](https://pkg.go.dev/github.com/jdkato/prose/v3/token) | the shared token type |
| [`summarize`](https://pkg.go.dev/github.com/jdkato/prose/v3/summarize) | readability metrics |
| [`strcase`](https://pkg.go.dev/github.com/jdkato/prose/v3/strcase) | title and sentence case |
