package prose

import (
	"path/filepath"
	"strings"
	"testing"
)

func BenchmarkDoc(b *testing.B) {
	content := readDataFile(filepath.Join(testdata, "sherlock.txt"))
	text := string(content)
	for n := 0; n < b.N; n++ {
		_, err := NewDocument(text)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkCustomTokenizer(b *testing.B) {
	content := readDataFile(filepath.Join(testdata, "sherlock.txt"))
	tok := NewIterTokenizer(
		UsingSanitizer(strings.NewReplacer()), // Disable sanitizer
		UsingPrefixes([]string{"(", `"`, "[", "'"}),
		UsingSuffixes([]string{",", ")", `"`, "]", "!", ";", ".", "?", ":", "'"}),
	)
	text := string(content)
	for n := 0; n < b.N; n++ {
		_, err := NewDocument(text, UsingTokenizer(tok))
		if err != nil {
			panic(err)
		}
	}
}
