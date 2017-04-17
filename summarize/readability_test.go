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
		a := d.Assess()
		fmt.Printf("Case: %d\n", i)
		fmt.Printf("AutomatedReadability: %0.2f\n", a.AutomatedReadability)
		fmt.Printf("FleschKincaid: %0.2f\n", a.FleschKincaid)
		fmt.Printf("ReadingEase: %0.2f\n", a.ReadingEase)
		fmt.Printf("SMOG: %0.2f\n", a.SMOG)
		fmt.Printf("Gunningfog: %0.2f\n", a.GunningFog)
		fmt.Printf("ColemanLiau: %0.2f\n", a.ColemanLiau)
		fmt.Printf("DaleChall: %0.2f\n", a.DaleChall)
		fmt.Printf("MeanGrade: %0.2f\n", a.MeanGradeLevel)
		fmt.Printf("StdGrade: %0.2f\n", a.StdGradeLevel)
	}
}

func BenchmarkReadability(b *testing.B) {
	in := util.ReadDataFile(filepath.Join(testdata, "sherlock.txt"))

	d := NewDocument(string(in))
	for n := 0; n < b.N; n++ {
		d.Assess()
	}
}
