package transform

import (
	"fmt"
	"path/filepath"
)

var testdata = filepath.Join("..", "testdata")

func ExampleNewTitleConverter() {
	tc := NewTitleConverter(APStyle)
	fmt.Println(tc.Title("the last of the mohicans"))
	// Output: The Last of the Mohicans
}
