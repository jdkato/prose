// Package conllu reads part-of-speech training data in the CoNLL-U format
// used by Universal Dependencies.
//
// Only what a tagger needs is kept: the surface form of each token and its
// universal (or, on request, language-specific) tag. Syntax is ignored.
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

// ReadFile reads every sentence in path.
func ReadFile(path string, col Column) ([]Sentence, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	sents, err := Read(f, col)
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
func Read(r io.Reader, col Column) ([]Sentence, error) {
	var (
		sents []Sentence
		cur   Sentence
		skip  int // syntactic words still covered by a range line
	)
	flush := func() {
		if len(cur.Words) > 0 {
			sents = append(sents, cur)
		}
		cur = Sentence{}
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
				cur.Words = append(cur.Words, rangeForm)
				cur.Tags = append(cur.Tags, tag)
				rangeForm = ""
			}
			continue
		}

		cur.Words = append(cur.Words, fields[1])
		cur.Tags = append(cur.Tags, tag)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	flush()

	return sents, nil
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
