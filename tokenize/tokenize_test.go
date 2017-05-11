package tokenize

func ExampleNewWordBoundaryTokenizer() {
	t := NewWordBoundaryTokenizer()
	t.Tokenize("They'll save and invest more.")
	// Output: [They'll save and invest more]
}

func ExampleNewWordPunctTokenizer() {
	t := NewWordPunctTokenizer()
	t.Tokenize("They'll save and invest more.")
	// Output: [They 'll save and invest more .]
}
