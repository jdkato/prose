package aptag

import "encoding/json"

// AveragedPerceptron ...
type AveragedPerceptron struct {
	Classes   []string
	Instances int
	Stamps    map[string]int
	TagMap    map[string]string
	Totals    map[string]int
	Weights   map[string]map[string]float64
}

// NewAveragedPerceptron ...
func NewAveragedPerceptron() *AveragedPerceptron {
	var ap AveragedPerceptron
	var err error

	ap.Totals = make(map[string]int)
	ap.Stamps = make(map[string]int)
	err = json.Unmarshal(getAsset("classes.json"), &ap.Classes)
	checkError(err)
	err = json.Unmarshal(getAsset("tags.json"), &ap.TagMap)
	checkError(err)
	err = json.Unmarshal(getAsset("weights.json"), &ap.Weights)
	checkError(err)
	return &ap
}

// Predict ...
func (ap AveragedPerceptron) Predict(features map[string]float64) string {
	var weights map[string]float64
	var found bool

	scores := make(map[string]float64)
	for feat, value := range features {
		if weights, found = ap.Weights[feat]; !found || value == 0 {
			continue
		}
		for label, weight := range weights {
			if _, ok := scores[label]; ok {
				scores[label] += value * weight
			} else {
				scores[label] = value * weight
			}
		}
	}
	return max(scores)
}

func max(scores map[string]float64) string {
	var class string
	max := 0.0
	for label, value := range scores {
		if value > max {
			max = value
			class = label
		}
	}
	return class
}
