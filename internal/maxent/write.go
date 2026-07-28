package maxent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
)

// Source is a model in its trained shape: a feature-string to dense-index
// mapping, a weight per index, and the label set.
type Source struct {
	Mapping map[string]int
	Weights []float64
	Labels  []string
}

// Marshal encodes src into the flat format.
//
// The mapping's dense indices are dissolved: keys are sorted and each key's
// weight is stored alongside it, so lookup does not need the index. Weights
// past len(Mapping) are GIS correction terms, preserved as "extra".
func Marshal(src Source) ([]byte, error) {
	keys := make([]string, 0, len(src.Mapping))
	for k := range src.Mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// cardinality mirrors newMaxentClassifier: the number of distinct feature
	// names in the mapping, plus one.
	names := make(map[string]struct{})
	for k := range src.Mapping {
		if i := indexByte(k, '-'); i >= 0 {
			names[k[:i]] = struct{}{}
		} else {
			names[k] = struct{}{}
		}
	}
	cardinality := len(names) + 1

	var (
		blob    bytes.Buffer
		keyOff  = make([]uint32, 0, len(keys)+1)
		weights = make([]float32, 0, len(keys))
	)
	keyOff = append(keyOff, 0)

	for _, k := range keys {
		idx := src.Mapping[k]
		if idx < 0 || idx >= len(src.Weights) {
			return nil, fmt.Errorf("key %q maps to index %d, outside the %d weights",
				k, idx, len(src.Weights))
		}
		blob.WriteString(k)
		keyOff = append(keyOff, uint32(blob.Len()))
		weights = append(weights, float32(src.Weights[idx]))
	}

	// Anything past the mapping's range is a GIS correction term.
	var extra []float32
	if len(src.Weights) > len(src.Mapping) {
		for _, w := range src.Weights[len(src.Mapping):] {
			extra = append(extra, float32(w))
		}
	}

	var out bytes.Buffer
	out.WriteString(Magic)
	putU32(&out, Version)
	putU32(&out, uint32(cardinality))

	putU32(&out, uint32(len(src.Labels)))
	for _, l := range src.Labels {
		if err := putStr16(&out, l); err != nil {
			return nil, err
		}
	}

	putU32(&out, uint32(len(keys)))
	putU32(&out, uint32(blob.Len()))
	out.Write(blob.Bytes())

	for _, v := range keyOff {
		putU32(&out, v)
	}
	for _, w := range weights {
		putU32(&out, math.Float32bits(w))
	}

	putU32(&out, uint32(len(extra)))
	for _, w := range extra {
		putU32(&out, math.Float32bits(w))
	}

	return out.Bytes(), nil
}

func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func putU32(b *bytes.Buffer, v uint32) {
	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], v)
	b.Write(tmp[:])
}

func putStr16(b *bytes.Buffer, s string) error {
	if len(s) > math.MaxUint16 {
		return fmt.Errorf("string of %d bytes exceeds the uint16 length field", len(s))
	}
	var tmp [2]byte
	binary.LittleEndian.PutUint16(tmp[:], uint16(len(s)))
	b.Write(tmp[:])
	b.WriteString(s)
	return nil
}
