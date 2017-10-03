package transform

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jdkato/prose/internal/util"
	"github.com/stretchr/testify/assert"
)

var testdata = filepath.Join("..", "testdata")

type testFormat struct {
	Input string
	None  string
	Snake string
	Param string
}

func TestTransform(t *testing.T) {
	tests := make([]testFormat, 0)
	cases := util.ReadDataFile(filepath.Join(testdata, "case.json"))

	util.CheckError(json.Unmarshal(cases, &tests))
	for _, test := range tests {
		assert.Equal(t, test.None, SimpleCase(test.Input))
		assert.Equal(t, test.Snake, SnakeCase(test.Input))
		assert.Equal(t, test.Param, DashCase(test.Input))
	}
}

func ExampleNewTitleConverter() {
	tc := NewTitleConverter(APStyle)
	fmt.Println(tc.Title("the last of the mohicans"))
	// Output: The Last of the Mohicans
}
