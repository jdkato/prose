package segment_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/jdkato/prose/v3/segment"
	_ "github.com/jdkato/prose/v3/segment/lang/all"
)

func joined(sents []segment.Sentence) string {
	out := make([]string, len(sents))
	for i, s := range sents {
		out[i] = s.Text
	}
	return strings.Join(out, "|")
}

func TestLanguageSegmentation(t *testing.T) {
	cases := []struct {
		lang string
		in   string
		want string
	}{
		// Punkt models know the language's abbreviations.
		{"de", "Das gilt z.B. für Berlin. Dr. Müller kam am 3. Mai. Er blieb.",
			"Das gilt z.B. für Berlin.|Dr. Müller kam am 3. Mai.|Er blieb."},
		{"fr", "« Bonjour ! » dit M. Dupont. Il est parti vers 10 h. Elle aussi.",
			"« Bonjour ! » dit M. Dupont.|Il est parti vers 10 h.|Elle aussi."},
		{"es", "¿Qué pasa? ¡Nada! El Sr. García no vino. Llegó la Sra. López.",
			"¿Qué pasa?|¡Nada!|El Sr. García no vino.|Llegó la Sra. López."},
		{"pt-BR", "O Dr. Silva chegou. Ele está bem.", "O Dr. Silva chegou.|Ele está bem."},
		{"nl", "Dhr. Jansen kwam. Hij bleef.", "Dhr. Jansen kwam.|Hij bleef."},

		// Rule-based: East Asian terminators need no space.
		{"ja", "東京に行きました。楽しかったです！また行きたい？「はい」と言った。",
			"東京に行きました。|楽しかったです！|また行きたい？|「はい」と言った。"},
		{"zh", "我爱北京。天安门很大！你去过吗？", "我爱北京。|天安门很大！|你去过吗？"},
		{"ko", "안녕하세요. 저는 학생입니다. 당신은요?", "안녕하세요.|저는 학생입니다.|당신은요?"},

		// Rule-based: initials, dotted abbreviations and list numbers do
		// not end a sentence; a lower-case continuation does not either.
		{"ru", "Он сказал, т.е. написал, что А. С. Пушкин был поэтом. Это правда.",
			"Он сказал, т.е. написал, что А. С. Пушкин был поэтом.|Это правда."},
		{"ru", "«Привет!» — сказал он. Она ушла... и не вернулась.",
			"«Привет!» — сказал он.|Она ушла... и не вернулась."},
		{"id", "Harganya 3.5 juta. Dia membeli 2. Lalu pergi.", "Harganya 3.5 juta.|Dia membeli 2. Lalu pergi."},
		{"hr", "Došao je u 10.30 sati. Otišao je u 11 sati.", "Došao je u 10.30 sati.|Otišao je u 11 sati."},
		{"ro", "A venit ieri. A plecat azi", "A venit ieri.|A plecat azi"},
	}

	for _, tc := range cases {
		seg, err := segment.ForLanguage(tc.lang)
		if err != nil {
			t.Fatalf("[%s] %v", tc.lang, err)
		}
		got := joined(seg.Segment(tc.in))
		if got != tc.want {
			t.Errorf("[%s] Segment(%q)\n  got  %s\n  want %s", tc.lang, tc.in, got, tc.want)
		}
		for _, s := range seg.Segment(tc.in) {
			if s.In(tc.in) != s.Text {
				t.Errorf("[%s] sentence %q at %d does not locate its text", tc.lang, s.Text, s.Start)
			}
		}
	}
}

func TestForLanguage(t *testing.T) {
	de, err := segment.ForLanguage("DE")
	if err != nil || !de.Trained() {
		t.Errorf("German should have a trained model (err=%v)", err)
	}
	ja, err := segment.ForLanguage("ja")
	if err != nil || ja.Trained() {
		t.Errorf("Japanese should be rule-based (err=%v)", err)
	}
	langs := segment.Languages()
	if len(langs) != 17 || langs[0] != "cs" {
		t.Errorf("Languages() = %v", langs)
	}
}

// The rule-based segmenter must cover every byte of its input: trimming is
// the caller's job, and a dropped span would lose a sentence.
func TestRulesCoverInput(t *testing.T) {
	seg, _ := segment.ForLanguage("xx")
	for _, src := range []string{"", "   ", "One. Two! Three?", "No terminator", "。。。", "a.b. c. D", "end.", "x!) y"} {
		var all []string
		for _, s := range seg.Segment(src) {
			all = append(all, s.Text)
		}
		if strings.Join(all, "") != strings.Join(strings.Fields(src), "")[:0]+strings.ReplaceAll(strings.Join(strings.Fields(src), ""), "", "") {
			// Compare without whitespace: trimming removes it.
			if strings.Join(strings.Fields(strings.Join(all, " ")), "") != strings.Join(strings.Fields(src), "") {
				t.Errorf("Segment(%q) = %q loses text", src, all)
			}
		}
	}
}

// A model taken out of the registry and read back through FromBytes segments
// as the registered one does, abbreviations included.
func TestFromBytesMatchesRegistered(t *testing.T) {
	gz, abbrevs, ok := segment.LanguageModel("de")
	if !ok {
		t.Fatal("German should be registered")
	}

	loose, err := segment.FromBytes(gz, abbrevs...)
	if err != nil {
		t.Fatal(err)
	}
	registered, err := segment.ForLanguage("de")
	if err != nil {
		t.Fatal(err)
	}

	src := "Dr. Müller kommt z.B. morgen. Er bleibt."
	if got, want := loose.SegmentText(src), registered.SegmentText(src); !reflect.DeepEqual(got, want) {
		t.Errorf("FromBytes = %q, ForLanguage = %q", got, want)
	}
	if n := len(loose.SegmentText(src)); n != 2 {
		t.Errorf("got %d sentences, want 2", n)
	}

	if _, _, ok := segment.LanguageModel("xx"); ok {
		t.Error("an unregistered language should not have a model")
	}
	if _, err := segment.FromBytes([]byte("not gzip")); err == nil {
		t.Error("garbage should not decode")
	}
}
