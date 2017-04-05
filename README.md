# prose

[![Build Status](https://travis-ci.org/jdkato/prose.svg?branch=master)](https://travis-ci.org/jdkato/prose)

`prose` is an in-progress Go library for text processing that supports tokenization, part of speech tagging, text transforming, and text summarization.

```
go get github.com/jdkato/prose
```

## Tokenizing

`TreebankWordTokenizer` is a port if NLTK's [Treebank tokenizer](https://github.com/nltk/nltk/blob/develop/nltk/tokenize/treebank.py), which is based on a [sed script](https://github.com/andre-martins/TurboParser/blob/master/scripts/tokenizer.sed) written by Robert McIntyre.

```go
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

## Tagging

`PerceptronTagger` is a port of Textblob's "fast and accurate" [POS tagger](https://github.com/sloria/textblob-aptagger).

```go
import (
    "fmt"

    "github.com/jdkato/prose/tag"
)

func main() {
    text := "A fast and accurate part-of-speech tagger for Golang."
    words := tokenize.NewTreebankWordTokenizer().Tokenize(text)

    tagger := NewPerceptronTagger()
    for _, tok := range tagger.Tag(words) {
        fmt.Println(tok.Text, tok.Tag)
    }
}
```

It performs quite well on NLTK's `treebank` corpus:

| Library | Accuracy | Time (sec) |
|:--------|---------:|-----------:|
| NLTK    |    0.893 |       7.55 |
| `prose` |    0.961 |      3.056 |

(see [`scripts/test_model.py`](https://github.com/jdkato/aptag/blob/master/scripts/test_model.py).)

## Transforming

`Title` converts a string to title case, while attempting to adhere to common guidelines.

```go
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

## Summarizing
