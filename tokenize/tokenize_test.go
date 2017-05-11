package tokenize

import "fmt"

func ExampleNewWordBoundaryTokenizer() {
	t := NewWordBoundaryTokenizer()
	fmt.Println(t.Tokenize("They'll save and invest more."))
	// Output: [They'll save and invest more]
}

func ExampleNewWordPunctTokenizer() {
	t := NewWordPunctTokenizer()
	fmt.Println(t.Tokenize("They'll save and invest more."))
	// Output: [They ' ll save and invest more .]
}

func ExampleNewTreebankWordTokenizer() {
	t := NewTreebankWordTokenizer()
	fmt.Println(t.Tokenize("They'll save and invest more."))
	// Output: [They 'll save and invest more .]
}
