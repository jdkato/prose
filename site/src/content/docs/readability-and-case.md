---
title: Readability and case
summary: Readability scores and title or sentence case conversion.
order: 6
---

# Readability and case conversion

## Readability

`summarize` computes the usual readability scores over English text:

```go
import "github.com/jdkato/prose/v3/summarize"

d := summarize.NewDocument("The quick brown fox jumps over the lazy dog.")
fmt.Println(d.FleschKincaid(), d.GunningFog(), d.SMOG())
```

| Method | Score |
| --- | --- |
| `FleschKincaid` | Flesch–Kincaid grade level |
| `FleschReadingEase` | Flesch reading ease, 0 to 100 |
| `GunningFog` | Gunning fog index |
| `SMOG` | SMOG grade |
| `ColemanLiau` | Coleman–Liau index |
| `AutomatedReadability` | Automated readability index |
| `DaleChall` | Dale–Chall score |
| `LIX` | Läsbarhetsindex |

`Assess` returns all of them at once, with the mean and standard deviation of
the grade-level scores. The document also exposes its counts (`NumWords`,
`NumSentences`, `NumSyllables` and so on) and `Keywords`, `WordDensity` and
`Summary` for the most representative paragraphs.

## Case conversion

`strcase` writes a string in a case convention:

```go
import "github.com/jdkato/prose/v3/strcase"

tc := strcase.NewTitleConverter(strcase.APStyle)
fmt.Println(tc.Convert("the quick brown fox"))
// The Quick Brown Fox
```

`APStyle` and `ChicagoStyle` differ in which short words stay lowercase.
`NewSentenceConverter` lowercases everything but the first word, proper nouns
and acronyms. Both take `UsingVocab` to keep a list of terms as written, and
`UsingPrefix` to leave a leading marker alone.

`Simple`, `Dash`, `Snake`, `Dot`, `Constant`, `Pascal` and `Camel` handle
identifier casing.
