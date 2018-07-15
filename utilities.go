package prose

import (
	"bytes"
	"encoding/gob"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
)

var maxLogDiff = math.Log2(1e-30)

type mappedProbDist struct {
	dict map[string]float64
	log  bool
}

func (m *mappedProbDist) prob(label string) float64 {
	if p, found := m.dict[label]; found {
		return math.Pow(2, p)
	}
	return 0.0
}

func (m *mappedProbDist) max() string {
	var class string
	max := math.Inf(-1)
	for label, value := range m.dict {
		if value > max {
			max = value
			class = label
		}
	}
	return class
}

func newMappedProbDist(dict map[string]float64, normalize bool) *mappedProbDist {
	if normalize {
		values := []float64{}
		for _, v := range dict {
			values = append(values, v)
		}
		sum := sumLogs(values)
		if sum <= math.Inf(-1) {
			p := math.Log2(1.0 / float64(len(dict)))
			for k := range dict {
				dict[k] = p
			}
		} else {
			for k := range dict {
				dict[k] -= sum
			}
		}
	}
	return &mappedProbDist{dict: dict, log: true}
}

func addLogs(x, y float64) float64 {
	if x < y+maxLogDiff {
		return y
	} else if y < x+maxLogDiff {
		return x
	}
	base := math.Min(x, y)
	return base + math.Log2(math.Pow(2, x-base)+math.Pow(2, y-base))
}

func sumLogs(logs []float64) float64 {
	if len(logs) == 0 {
		return math.Inf(-1)
	}
	sum := logs[0]
	for _, log := range logs[1:] {
		sum = addLogs(sum, log)
	}
	return sum
}

// checkError panics if `err` is not `nil`.
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// min returns the minimum of `a` and `b`.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// isPunct determines if the string represents a number.
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// stringInSlice determines if `slice` contains the string `a`.
func stringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if a == b {
			return true
		}
	}
	return false
}

func getAsset(folder, name string) *gob.Decoder {
	b, err := Asset(path.Join("model", folder, name))
	checkError(err)
	return gob.NewDecoder(bytes.NewReader(b))
}

func getDiskAsset(path string) *gob.Decoder {
	f, err := os.Open(path)
	checkError(err)
	return gob.NewDecoder(f)
}

func nSuffix(word string, length int) string {
	return strings.ToLower(word[len(word)-min(len(word), length):])
}

func nPrefix(word string, length int) string {
	return strings.ToLower(word[:min(len(word), length)])
}

func isBasic(word string) string {
	if stringInSlice(word, enWordList) {
		return "True"
	}
	return "False"
}
