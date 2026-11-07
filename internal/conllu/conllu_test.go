package conllu

import (
	"strings"
	"testing"
)

const sample = `# sent_id = 1
# text = Er geht zum Arzt.
1	Er	er	PRON	PPER	_	2	nsubj	_	_
2	geht	gehen	VERB	VVFIN	_	0	root	_	_
3-4	zum	_	_	_	_	_	_	_	_
3	zu	zu	ADP	APPR	_	5	case	_	_
4	dem	der	DET	ART	_	5	det	_	_
5	Arzt	Arzt	NOUN	NN	_	2	obl	_	SpaceAfter=No
5.1	_	_	X	_	_	_	_	_	_
6	.	.	PUNCT	$.	_	2	punct	_	_

1	Ja	ja	INTJ	ITJ	_	0	root	_	_
`

func TestReadCollapsesRangesAndSkipsEmptyNodes(t *testing.T) {
	sents, err := Read(strings.NewReader(sample), UPOS, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(sents) != 2 {
		t.Fatalf("got %d sentences, want 2", len(sents))
	}

	words := strings.Join(sents[0].Words, " ")
	tags := strings.Join(sents[0].Tags, " ")
	if words != "Er geht zum Arzt ." {
		t.Errorf("words = %q", words)
	}
	if tags != "PRON VERB ADP NOUN PUNCT" {
		t.Errorf("tags = %q", tags)
	}
	if sents[1].Words[0] != "Ja" || sents[1].Tags[0] != "INTJ" {
		t.Errorf("second sentence = %v / %v", sents[1].Words, sents[1].Tags)
	}
}

func TestReadXPOS(t *testing.T) {
	sents, err := Read(strings.NewReader(sample), XPOS, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if got := sents[0].Tags[1]; got != "VVFIN" {
		t.Errorf("xpos = %q, want VVFIN", got)
	}
}

const hyphens = `# text = A well-known Hewlett-Packard device - not e-mail.
1	A	a	DET	DT	_	5	det	_	_
2	well	well	ADV	RB	_	4	advmod	_	SpaceAfter=No
3	-	-	PUNCT	HYPH	_	4	punct	_	SpaceAfter=No
4	known	know	VERB	VBN	_	5	amod	_	_
5	Hewlett	Hewlett	PROPN	NNP	_	7	compound	_	SpaceAfter=No
6	-	-	PUNCT	HYPH	_	5	punct	_	SpaceAfter=No
7	Packard	Packard	PROPN	NNP	_	8	compound	_	_
8	device	device	NOUN	NN	_	0	root	_	_
9	-	-	PUNCT	HYPH	_	12	punct	_	_
10	not	not	PART	RB	_	12	advmod	_	_
11	e	e	NOUN	NN	_	12	compound	_	SpaceAfter=No
12	-	-	PUNCT	HYPH	_	12	punct	_	SpaceAfter=No
13	mail	mail	NOUN	NN	_	8	conj	_	SpaceAfter=No
14	.	.	PUNCT	.	_	8	punct	_	_
`

func TestJoinHyphens(t *testing.T) {
	sents, err := Read(strings.NewReader(hyphens), XPOS, Options{JoinHyphens: true})
	if err != nil {
		t.Fatal(err)
	}
	words := strings.Join(sents[0].Words, " ")
	tags := strings.Join(sents[0].Tags, " ")
	if words != "A well-known Hewlett-Packard device - not e-mail ." {
		t.Errorf("words = %q", words)
	}
	// well-known modifies a noun, so it is an adjective; Hewlett-Packard
	// keeps its head's tag; the spaced hyphen is left alone.
	if tags != "DT JJ NNP NN HYPH RB NN ." {
		t.Errorf("tags = %q", tags)
	}

	plain, _ := Read(strings.NewReader(hyphens), XPOS, Options{})
	if len(plain[0].Words) != 14 {
		t.Errorf("without the option, got %d tokens, want 14", len(plain[0].Words))
	}
}
