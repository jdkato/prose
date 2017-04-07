package summarize

// WordDensity returns a map of each word and its density.
func (d *Document) WordDensity() map[string]float64 {
	density := make(map[string]float64)
	for word, count := range d.WordFrequency {
		density[word] = float64(count) / d.NumWords
	}
	return density
}

// AverageWordLength returns the average number of characters per word.
func (d *Document) AverageWordLength() float64 {
	return d.NumCharacters / d.NumWords
}
