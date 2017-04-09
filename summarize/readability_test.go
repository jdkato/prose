package summarize

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jdkato/prose/internal/util"
)

var testdata = filepath.Join("..", "testdata")

func TestReadability(t *testing.T) {
	tests := make([]testCase, 0)
	cases := util.ReadDataFile(filepath.Join(testdata, "summarize.json"))

	util.CheckError(json.Unmarshal(cases, &tests))
	for i, test := range tests {
		d := NewDocument(test.Text)
		fmt.Printf("Case: %d\n", i)
		fmt.Printf("AutomatedReadability: %0.2f\n", d.AutomatedReadability())
		fmt.Printf("FleschKincaid: %0.2f\n", d.FleschKincaid())
		fmt.Printf("ReadingEase: %0.2f\n", d.FleschReadingEase())
		fmt.Printf("SMOG: %0.2f\n", d.SMOG())
		fmt.Printf("Gunningfog: %0.2f\n", d.GunningFog())
		fmt.Printf("ColemanLiau: %0.2f\n", d.ColemanLiau())
		fmt.Printf("DaleChall: %0.2f\n", d.DaleChall())
	}
}

func BenchmarkReadability(b *testing.B) {
	in := util.ReadDataFile(filepath.Join(testdata, "sherlock.txt"))

	d := NewDocument(string(in))
	for n := 0; n < b.N; n++ {
		d.Assess()
	}
}
