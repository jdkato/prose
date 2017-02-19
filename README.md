# aptag

An English-language Part-of-Speech Tagger:

```go
import (
    "fmt"

    "github.com/jdkato/aptag"
)

func main() {
    text := "Dive into NLTK: Part-of-speech tagging and POS Tagger."
    tagger := aptag.NewPerceptronTagger()

    tokens = tagger.TokenizeAndTag(string(text))
    for _, tok := range tokens {
        fmt.Println(tok.Text, tok.Tag)
    }
}
```

## Install

```
go get github.com/jdkato/aptag
```

## Performance

| Library | Accuracy | Time (sec) |
|:--------|---------:|-----------:|
| NLTK    |    0.893 |      6.755 |
| aptag   |    0.961 |      2.879 |

(see `scripts/test_model.py`.)

## Notice

This is a port of [`textblob-aptagger`](https://github.com/sloria/textblob-aptagger):

```
Copyright 2013 Matthew Honnibal

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
