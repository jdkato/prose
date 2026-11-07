---
title: Getting started
summary: Install prose and run the whole pipeline over a string.
order: 1
---

# Getting started

`prose` is a natural language processing library for Go. It is pure Go, with
no cgo and no external services.

- **Tokenization**: words and punctuation, with byte offsets, in any language
- **Segmentation**: sentences, via the Punkt algorithm, in 12 languages
- **Tagging**: part of speech, in 12 languages
- **Entities**: people, places and organizations (English)
- **Readability**: Flesch–Kincaid, Gunning fog, SMOG, Coleman–Liau and more
- **Case conversion**: title and sentence case, AP and Chicago styles

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

`NewDocument` tokenizes, segments, tags and extracts entities. Each stage can
be switched off; see [Choosing stages](stages).

## Try it

The [playground](../) on the home page runs this same code in your browser,
compiled to WebAssembly, so you can see what each stage produces for your own
text.
