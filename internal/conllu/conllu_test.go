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
	sents, err := Read(strings.NewReader(sample), UPOS)
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
	sents, err := Read(strings.NewReader(sample), XPOS)
	if err != nil {
		t.Fatal(err)
	}
	if got := sents[0].Tags[1]; got != "VVFIN" {
		t.Errorf("xpos = %q, want VVFIN", got)
	}
}
