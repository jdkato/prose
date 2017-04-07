# prose

[![Build Status](https://travis-ci.org/jdkato/prose.svg?branch=master)](https://travis-ci.org/jdkato/prose)
[![Build status](https://ci.appveyor.com/api/projects/status/24bepq85nnnk4scr/branch/master?svg=true)](https://ci.appveyor.com/project/jdkato/prose/branch/master)  [![GoDoc](https://godoc.org/github.com/jdkato/prose?status.svg)](https://godoc.org/github.com/jdkato/prose) [![Go Report Card](https://goreportcard.com/badge/github.com/jdkato/prose)](https://goreportcard.com/report/github.com/jdkato/prose) [![Code Climate](https://codeclimate.com/github/jdkato/prose/badges/gpa.svg)](https://codeclimate.com/github/jdkato/prose) [![license](https://img.shields.io/github/license/mashape/apistatus.svg)]()

`prose` is Go library for text processing that supports tokenization, part of speech tagging, and various other prose-related functions.

## Tokenizing

```go
package main

import (
    "fmt"

    "github.com/jdkato/prose/tokenize"
)

func main() {
    text := "They'll save and invest more."
    tokenizer := tokenize.NewTreebankWordTokenizer()
    for _, word := range tokenizer.Tokenize(text) {
        // [They 'll save and invest more .]
        fmt.Println(word)
    }
}
```

`TreebankWordTokenizer` is a port if NLTK's [Treebank tokenizer](https://github.com/nltk/nltk/blob/develop/nltk/tokenize/treebank.py), which is based on a [sed script](https://github.com/andre-martins/TurboParser/blob/master/scripts/tokenizer.sed) written by Robert McIntyre.

## Tagging

```go
package main

import (
    "fmt"

    "github.com/jdkato/prose/tag"
    "github.com/jdkato/prose/tokenize"
)

func main() {
    text := "A fast and accurate part-of-speech tagger for Golang."
    words := tokenize.NewTreebankWordTokenizer().Tokenize(text)

    tagger := tag.NewPerceptronTagger()
    for _, tok := range tagger.Tag(words) {
        fmt.Println(tok.Text, tok.Tag)
    }
}
```

`PerceptronTagger` is a port of Textblob's "fast and accurate" [POS tagger](https://github.com/sloria/textblob-aptagger). It performs quite well on NLTK's `treebank` corpus:

| Library | Accuracy | Time (sec) |
|:--------|---------:|-----------:|
| NLTK    |    0.893 |       7.55 |
| `prose` |    0.961 |      3.056 |

(see [`scripts/test_model.py`](https://github.com/jdkato/aptag/blob/master/scripts/test_model.py).)

## Transforming

```go
package main

import (
    "fmt"
    "strings"

    "github.com/jdkato/prose/transform"
)

func main() {
    text := "the last of the mohicans"
    fmt.Println(strings.Title(text))   // The Last Of The Mohicans
    fmt.Println(transform.Title(text)) // The Last of the Mohicans
}
```

`Title` converts a string to title case, while attempting to adhere to common guidelines. Inspiration and test data taken from [python-titlecase](https://github.com/ppannuto/python-titlecase) and [to-title-case](https://github.com/gouch/to-title-case).

## Summarizing

```go
package main

import (
    "fmt"

    "github.com/jdkato/prose/summarize"
)

func main() {
    doc := summarize.NewDocument("This is some interesting text.")
    fmt.Println(doc.SMOG())
}
```
