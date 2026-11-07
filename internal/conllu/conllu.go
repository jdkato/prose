// Package conllu reads part-of-speech training data in the CoNLL-U format
// used by Universal Dependencies.
//
// Only what a tagger needs is kept: the surface form of each token and its
// universal (or, on request, language-specific) tag. Syntax is read only to
// decide what a joined hyphenated span is; see Options.
package conllu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Sentence is one tagged sentence: Words[i] carries Tags[i].
type Sentence struct {
	Words []string
	Tags  []string
}

// Column selects which tag column to read.
type Column int

const (
	// UPOS is the universal part-of-speech tag, column 4.
	UPOS Column = 3
	// XPOS is the language-specific tag, column 5.
	XPOS Column = 4
)

// Options adjusts what Read keeps.
type Options struct {
	// JoinHyphens reads a hyphenated span written without spaces -- "well",
	// "-", "known" -- as the one token a tokenizer that keeps hyphens sees.
	// The joined token takes the tag of the piece that heads the span, or
	// the adjective tag when the span modifies a noun ("well-known" is JJ,
	// not VBN).
	JoinHyphens bool
}

// A token is one line of a sentence, with what the join needs.
type token struct {
	id     int
	form   string
	tag    string
	head   int
	deprel string
	glued  bool // SpaceAfter=No
}

// ReadFile reads every sentence in path.
func ReadFile(path string, col Column, opt Options) ([]Sentence, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	sents, err := Read(f, col, opt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return sents, nil
}

// Read parses CoNLL-U from r.
//
// A multiword token -- a range line such as "6-7 zum" for the syntactic words
// "zu" and "dem" -- is read as one word, the surface form, carrying the tag
// of its first syntactic word. That is the token a tokenizer working on
// surface text produces, and the first word's tag (ADP for a contracted
// preposition, VERB for a verb with a clitic attached) is the one a rule is
// after. Empty nodes ("8.1") are skipped.
func Read(r io.Reader, col Column, opt Options) ([]Sentence, error) {
	var (
		sents []Sentence
		cur   []token
		skip  int // syntactic words still covered by a range line
	)
	flush := func() {
		if len(cur) > 0 {
			if opt.JoinHyphens {
				cur = joinHyphens(cur, col)
			}
			s := Sentence{Words: make([]string, len(cur)), Tags: make([]string, len(cur))}
			for i, t := range cur {
				s.Words[i], s.Tags[i] = t.form, t.tag
			}
			sents = append(sents, s)
		}
		cur = nil
		skip = 0
	}

	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

	var rangeForm string
	for n := 1; sc.Scan(); n++ {
		line := sc.Text()
		if line == "" {
			flush()
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) < 5 {
			return nil, fmt.Errorf("line %d: want at least 5 tab-separated fields, have %d", n, len(fields))
		}
		id := fields[0]

		if strings.Contains(id, ".") {
			continue
		}

		if lo, hi, ok := parseRange(id); ok {
			rangeForm = fields[1]
			skip = hi - lo + 1
			continue
		}

		tag := fields[col]
		if skip > 0 {
			skip--
			if rangeForm != "" {
				// First syntactic word of the range: it names the tag.
				cur = append(cur, token{form: rangeForm, tag: tag})
				rangeForm = ""
			}
			continue
		}

		t := token{form: fields[1], tag: tag}
		t.id, _ = strconv.Atoi(id)
		if len(fields) >= 10 {
			t.head, _ = strconv.Atoi(fields[6])
			t.deprel, _, _ = strings.Cut(fields[7], ":")
			t.glued = strings.Contains(fields[9], "SpaceAfter=No")
		}
		cur = append(cur, t)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	flush()

	return sents, nil
}

// joinHyphens merges each run of pieces glued together by hyphens into one
// token. The run's tag comes from the piece whose head lies outside it.
func joinHyphens(toks []token, col Column) []token {
	adj := "ADJ"
	if col == XPOS {
		adj = "JJ"
	}

	out := make([]token, 0, len(toks))
	for i := 0; i < len(toks); {
		j := i
		for isPiece(toks[j]) && j+2 < len(toks) && toks[j].glued &&
			toks[j+1].form == "-" && toks[j+1].glued && isPiece(toks[j+2]) {
			j += 2
		}
		if j == i {
			out = append(out, toks[i])
			i++
			continue
		}

		span := toks[i : j+1]
		inside := map[int]bool{}
		for _, t := range span {
			inside[t.id] = true
		}
		head := span[0]
		for _, t := range span {
			if t.form != "-" && !inside[t.head] {
				head = t
				break
			}
		}

		joined := token{id: head.id, tag: head.tag, head: head.head, deprel: head.deprel, glued: span[len(span)-1].glued}
		if head.deprel == "amod" {
			joined.tag = adj
		}
		var b strings.Builder
		for _, t := range span {
			b.WriteString(t.form)
		}
		joined.form = b.String()

		out = append(out, joined)
		i = j + 1
	}
	return out
}

// isPiece reports whether a token can be one side of a hyphen: a word, not
// punctuation and not itself a hyphen.
func isPiece(t token) bool {
	if t.form == "" || t.form == "-" || strings.ContainsAny(t.form, " \t") {
		return false
	}
	for _, r := range t.form {
		if r > ' ' && !strings.ContainsRune(`"'()[]{},.;:!?/\`, r) {
			return true
		}
	}
	return false
}

func parseRange(id string) (int, int, bool) {
	dash := strings.IndexByte(id, '-')
	if dash <= 0 {
		return 0, 0, false
	}
	lo, err1 := strconv.Atoi(id[:dash])
	hi, err2 := strconv.Atoi(id[dash+1:])
	if err1 != nil || err2 != nil || hi < lo {
		return 0, 0, false
	}
	return lo, hi, true
}
