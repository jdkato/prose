# prose

[![CI](https://github.com/jdkato/prose/actions/workflows/ci.yml/badge.svg)](https://github.com/jdkato/prose/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jdkato/prose/v3.svg)](https://pkg.go.dev/github.com/jdkato/prose/v3)
[![Coverage Status](https://coveralls.io/repos/github/jdkato/prose/badge.svg?branch=master)](https://coveralls.io/github/jdkato/prose?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/jdkato/prose)](https://goreportcard.com/report/github.com/jdkato/prose)

`prose` is a natural language processing library for Go — pure Go, no cgo, no
external services. English only.

- **Tokenization** — words and punctuation, with byte offsets
- **Segmentation** — sentences, via the Punkt algorithm
- **Tagging** — part of speech, using the Penn Treebank tag set
- **Entities** — people, places and organizations
- **Readability** — Flesch–Kincaid, Gunning fog, SMOG, Coleman–Liau and more
- **Case conversion** — title and sentence case, AP and Chicago styles

## Install

```console
go get github.com/jdkato/prose/v3
```

## Quick start

```go
package main

import (
    "fmt"
    "log"

    "github.com/jdkato/prose/v3"
)

func main() {
    doc, err := prose.NewDocument("Marie Curie studied in Paris. She won the Nobel Prize twice.")
    if err != nil {
        log.Fatal(err)
    }

    for _, sent := range doc.Sentences() {
        fmt.Println(sent.Text)
    }

    for _, tok := range doc.Tokens() {
        fmt.Println(tok.Text, tok.Tag)
    }

    for _, ent := range doc.Entities() {
        fmt.Println(ent.Text, ent.Label)
        // Marie Curie PERSON
        // Paris GPE
    }
}
```

## Everything knows where it came from

Every token, sentence and entity carries a byte offset into the source, and the
text at that offset is exactly the text you were given:

```go
src[tok.Start:tok.End()] == tok.Text
```

That holds for multi-byte input, repeated words and irregular whitespace:

```go
src := "He said “hi there” to me."
doc, _ := prose.NewDocument(src)

for _, tok := range doc.Tokens() {
    fmt.Printf("%-8q [%d:%d] %q\n", tok.Text, tok.Start, tok.End(), tok.In(src))
}
```

It is also why `prose` never normalizes text before tokenizing. Replacing curly
quotes with straight ones, or `&rsquo;` with `'`, changes byte lengths and
silently shifts every offset that follows. Normalization that would move bytes
is the caller's decision, made before `prose` sees the text.

## Pay only for what you use

`prose` ships two trained models: a part-of-speech tagger (~3 MB) and an entity
recognizer (~4 MB). Each lives in its own package, so the Go linker leaves out
whatever you do not import.

Importing the root package links both. If you only need tags, import the tagger
directly and the entity model never enters your binary — a 4.5 MB difference:

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
| [`ner`](https://pkg.go.dev/github.com/jdkato/prose/v3/ner) | named entities |
| [`token`](https://pkg.go.dev/github.com/jdkato/prose/v3/token) | the shared token type |
| [`summarize`](https://pkg.go.dev/github.com/jdkato/prose/v3/summarize) | readability metrics |
| [`strcase`](https://pkg.go.dev/github.com/jdkato/prose/v3/strcase) | title and sentence case |

## Choosing stages

Skipping entity extraction avoids most of the work:

```go
doc, err := prose.NewDocument(text, prose.WithExtraction(false))
```

`WithTokenization`, `WithSegmentation`, `WithTagging` and `WithExtraction` each
take a bool. Extraction requires tagging, and tagging requires tokenization, so
those are enabled automatically when needed.

## Cancellation

Work scales with the length of the input — roughly 1.4 s for a million words —
so `NewDocumentContext` accepts a `context.Context`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

doc, err := prose.NewDocumentContext(ctx, veryLargeText)
```

The individual tokenizers and taggers take no context: they work at microsecond
scale, where there is nothing worth cancelling.

## Concurrency

`Tokenizer`, `Segmenter`, `Tagger` and `Extracter` are all safe for concurrent
use. Their models are immutable once loaded, and per-call scratch space is
pooled.

## Readability and case conversion

```go
import "github.com/jdkato/prose/v3/summarize"

d := summarize.NewDocument("The quick brown fox jumps over the lazy dog.")
fmt.Println(d.FleschKincaid(), d.GunningFog(), d.SMOG())
```

```go
import "github.com/jdkato/prose/v3/strcase"

tc := strcase.NewTitleConverter(strcase.APStyle)
fmt.Println(tc.Convert("the quick brown fox"))
```

## License

MIT. See [LICENSE](LICENSE), which also records the upstream work `prose`
derives from.
