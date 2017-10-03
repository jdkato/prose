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
	Input    string
	None     string
	Snake    string
	Param    string
	Dot      string
	Constant string
	Pascal   string
	Camel    string
}

func TestTransform(t *testing.T) {
	tests := make([]testFormat, 0)
	cases := util.ReadDataFile(filepath.Join(testdata, "case.json"))

	util.CheckError(json.Unmarshal(cases, &tests))
	for _, test := range tests {
		assert.Equal(t, test.None, Simple(test.Input))
		assert.Equal(t, test.Snake, Snake(test.Input))
		assert.Equal(t, test.Param, Dash(test.Input))
		assert.Equal(t, test.Dot, Dot(test.Input))
		assert.Equal(t, test.Constant, Constant(test.Input))
		assert.Equal(t, test.Pascal, Pascal(test.Input))
		assert.Equal(t, test.Camel, Camel(test.Input))
	}
}

func ExampleNewTitleConverter() {
	tc := NewTitleConverter(APStyle)
	fmt.Println(tc.Title("the last of the mohicans"))
	// Output: The Last of the Mohicans
}
