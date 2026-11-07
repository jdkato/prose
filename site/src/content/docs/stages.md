---
title: Choosing stages
summary: Switch off the pipeline stages you do not need.
order: 5
---

# Choosing stages

Skipping entity extraction avoids most of the work:

```go
doc, err := prose.NewDocument(text, prose.WithExtraction(false))
```

`WithTokenization`, `WithSegmentation`, `WithTagging` and `WithExtraction` each
take a bool. Extraction requires tagging, and tagging requires tokenization, so
those are enabled automatically when needed.

## Cancellation

Work scales with the length of the input, roughly 1.4 s for a million words,
so `NewDocumentContext` accepts a `context.Context`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

doc, err := prose.NewDocumentContext(ctx, veryLargeText)
```

The context is polled between pipeline stages. The individual tokenizers and
taggers take no context: they work at microsecond scale, where there is nothing
worth cancelling.

## Concurrency

`Tokenizer`, `Segmenter`, `Tagger` and `Extracter` are all safe for concurrent
use. Their models are immutable once loaded, and per-call scratch space is
pooled.
