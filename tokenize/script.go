package tokenize

import (
	"unicode"
	"unicode/utf8"
)

// script classes a tokenizer without a dictionary can tell apart.
type script uint8

const (
	scriptOther script = iota
	scriptHan
	scriptHiragana
	scriptKatakana
	scriptCJKPunct
)

func scriptOf(r rune) script {
	switch {
	case unicode.Is(unicode.Han, r):
		return scriptHan
	case unicode.Is(unicode.Hiragana, r):
		return scriptHiragana
	case unicode.Is(unicode.Katakana, r), r == 'ー': // prolonged sound mark
		return scriptKatakana
	case r >= 0x3000 && r <= 0x303F, r >= 0xFF00 && r <= 0xFF0F,
		r >= 0xFF1A && r <= 0xFF20, r >= 0xFF3B && r <= 0xFF40,
		r >= 0xFF5B && r <= 0xFF65:
		// CJK symbols and punctuation, and the full-width ASCII punctuation
		// blocks, each its own token.
		return scriptCJKPunct
	}
	return scriptOther
}

// splitScripts appends span to tokens, divided wherever its script changes.
//
// Only the East Asian scripts count as a change: Latin letters, digits and
// Cyrillic are all "other", so a span of ordinary text comes back whole and
// the cost is one pass over its runes.
func splitScripts(span string, base int, tokens []Token) []Token {
	start := 0
	prev := scriptOther
	for i := 0; i < len(span); {
		r, size := utf8.DecodeRuneInString(span[i:])
		cur := scriptOf(r)
		// Every punctuation mark is its own token; runs of the same script
		// stay together.
		if i > 0 && (cur != prev || cur == scriptCJKPunct) {
			tokens = append(tokens, Token{Text: span[start:i], Start: base + start})
			start = i
		}
		prev = cur
		i += size
	}
	if start < len(span) {
		tokens = append(tokens, Token{Text: span[start:], Start: base + start})
	}
	return tokens
}
