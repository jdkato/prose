package flat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
)

// Source is a model in its trained shape: feature -> class -> weight, plus the
// unambiguous word -> tag map and the class list.
type Source struct {
	Weights map[string]map[string]float64
	TagMap  map[string]string
	Classes []string
}

// Marshal encodes src into the flat format.
//
// Classes are renumbered to a dense 0..n-1 space in the order given by
// src.Classes, extended with any class that appears in the weights or tag map
// but is missing from that list.
func Marshal(src Source) ([]byte, error) {
	classID := make(map[string]uint8, len(src.Classes))
	var classes []string
	addClass := func(name string) (uint8, error) {
		if id, ok := classID[name]; ok {
			return id, nil
		}
		if len(classes) >= MaxClasses {
			return 0, fmt.Errorf("more than %d classes; format stores IDs as uint8",
				MaxClasses)
		}
		id := uint8(len(classes))
		classID[name] = id
		classes = append(classes, name)
		return id, nil
	}
	for _, c := range src.Classes {
		if _, err := addClass(c); err != nil {
			return nil, err
		}
	}

	feats := make([]string, 0, len(src.Weights))
	for f := range src.Weights {
		feats = append(feats, f)
	}
	sort.Strings(feats)

	var (
		blob     bytes.Buffer
		featOff  = make([]uint32, 0, len(feats)+1)
		pairOff  = make([]uint32, 0, len(feats)+1)
		classIDs []uint8
		weights  []float32
	)
	featOff = append(featOff, 0)
	pairOff = append(pairOff, 0)

	for _, f := range feats {
		blob.WriteString(f)
		featOff = append(featOff, uint32(blob.Len()))

		// Sort each feature's classes by ID so the layout is deterministic
		// and a given model always encodes to identical bytes.
		inner := src.Weights[f]
		ids := make([]uint8, 0, len(inner))
		byID := make(map[uint8]float64, len(inner))
		for cls, w := range inner {
			id, err := addClass(cls)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
			byID[id] = w
		}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

		for _, id := range ids {
			classIDs = append(classIDs, id)
			weights = append(weights, float32(byID[id]))
		}
		pairOff = append(pairOff, uint32(len(weights)))
	}

	// Sort tagMap for reproducible output.
	tagWords := make([]string, 0, len(src.TagMap))
	for w := range src.TagMap {
		tagWords = append(tagWords, w)
	}
	sort.Strings(tagWords)

	var out bytes.Buffer
	out.WriteString(Magic)
	putU32(&out, Version)

	putU32(&out, uint32(len(classes)))
	for _, c := range classes {
		if err := putStr16(&out, c); err != nil {
			return nil, err
		}
	}

	putU32(&out, uint32(len(feats)))
	putU32(&out, uint32(len(weights)))
	putU32(&out, uint32(blob.Len()))
	out.Write(blob.Bytes())

	for _, v := range featOff {
		putU32(&out, v)
	}
	for _, v := range pairOff {
		putU32(&out, v)
	}
	out.Write(classIDs)
	for _, w := range weights {
		putU32(&out, math.Float32bits(w))
	}

	putU32(&out, uint32(len(tagWords)))
	for _, w := range tagWords {
		id, err := addClass(src.TagMap[w])
		if err != nil {
			return nil, err
		}
		if err := putStr16(&out, w); err != nil {
			return nil, err
		}
		out.WriteByte(id)
	}

	// A class first seen in the tag map would not be in the header written
	// above, which would corrupt the file. Fail loudly rather than emit
	// something subtly wrong.
	if len(classes) != int(readU32(out.Bytes()[len(Magic)+4:])) {
		return nil, fmt.Errorf(
			"class set grew from %d to %d while writing the tag map; "+
				"tag map references a class absent from the weights",
			readU32(out.Bytes()[len(Magic)+4:]), len(classes))
	}

	return out.Bytes(), nil
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

func readU32(b []byte) uint32 { return binary.LittleEndian.Uint32(b) }
