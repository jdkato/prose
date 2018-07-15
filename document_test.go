package prose

import (
	"path/filepath"
	"testing"
)

func BenchmarkDoc(b *testing.B) {
	content := readDataFile(filepath.Join(testdata, "sherlock.txt"))
	text := string(content)
	for n := 0; n < b.N; n++ {
		NewDocument(text)
	}
}
