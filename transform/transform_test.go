package transform

func ExampleNewTitleConverter() {
	tc := NewTitleConverter(APStyle)
	tc.Title("the last of the mohicans")
	// Output: The Last of the Mohicans
}
