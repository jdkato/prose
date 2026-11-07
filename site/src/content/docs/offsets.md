---
title: Offsets
summary: Every token, sentence and entity carries a byte offset into the source.
order: 2
---

# Everything knows where it came from

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

Offsets are bytes, not runes or UTF-16 code units, because that is what Go
string slicing uses. A caller that needs another unit converts at the edge.

## No normalization

It is also why `prose` never normalizes text before tokenizing. Replacing curly
quotes with straight ones, or `&rsquo;` with `'`, changes byte lengths and
silently shifts every offset that follows. Normalization that would move bytes
is the caller's decision, made before `prose` sees the text.
