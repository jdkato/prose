package summarize

// WordDensity returns a map of each word and its density in the Document d.
func (d *Document) WordDensity() map[string]float64 {
	density := make(map[string]float64)
	for word, count := range d.WordFrequency {
		density[word] = float64(count) / d.NumWords
	}
	return density
}
