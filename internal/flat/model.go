package flat

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sort"
)

// Model is a flat averaged-perceptron model, ready for scoring.
//
// Lookup is by feature string. Features are held in one sorted byte blob and
// found by binary search, which keeps the whole structure pointer-free: the GC
// never has to trace 75,447 individual strings.
type Model struct {
	classes []string // class ID -> name

	blob    string   // all feature strings, concatenated, sorted bytewise
	featOff []uint32 // nFeatures+1 offsets into blob
	pairOff []uint32 // nFeatures+1 offsets into classIDs/weights

	classIDs []uint8
	weights  []float32

	// tagMap holds unambiguous word -> class assignments, bypassing scoring.
	tagMap map[string]uint8
}

// Classes returns the class names, indexed by class ID.
func (m *Model) Classes() []string { return m.classes }

// NumFeatures reports how many distinct features the model knows.
func (m *Model) NumFeatures() int { return len(m.featOff) - 1 }

// NumPairs reports how many feature-class weights the model holds.
func (m *Model) NumPairs() int { return len(m.weights) }

// TagFor returns the pre-assigned tag for an unambiguous word.
func (m *Model) TagFor(word string) (string, bool) {
	if id, ok := m.tagMap[word]; ok {
		return m.classes[id], true
	}
	return "", false
}

// feature returns the index of feat, or -1. Binary search over the blob.
func (m *Model) feature(feat string) int {
	n := len(m.featOff) - 1
	i := sort.Search(n, func(i int) bool {
		return m.blob[m.featOff[i]:m.featOff[i+1]] >= feat
	})
	if i < n && m.blob[m.featOff[i]:m.featOff[i+1]] == feat {
		return i
	}
	return -1
}

// AccumulateInto adds the weights of feat, scaled by value, into scores.
// scores must be at least len(m.classes) long. Reports whether feat was known.
//
// This is the inner loop: no allocation, no map lookups, and the weights for
// a feature are contiguous.
func (m *Model) AccumulateInto(feat string, value float32, scores []float32) bool {
	i := m.feature(feat)
	if i < 0 {
		return false
	}
	lo, hi := m.pairOff[i], m.pairOff[i+1]
	ids := m.classIDs[lo:hi]
	wts := m.weights[lo:hi]
	for j := range ids {
		scores[ids[j]] += value * wts[j]
	}
	return true
}

// Best returns the highest-scoring class name.
//
// Ties break toward the lowest class ID, so the result is deterministic.
func (m *Model) Best(scores []float32) string {
	best, bestID := float32(math.Inf(-1)), -1
	for id := range m.classes {
		if scores[id] > best {
			best, bestID = scores[id], id
		}
	}
	if bestID < 0 {
		return ""
	}
	return m.classes[bestID]
}

// --- decoding -------------------------------------------------------------

type reader struct {
	b   []byte
	pos int
	err error
}

func (r *reader) fail(format string, args ...any) {
	if r.err == nil {
		r.err = fmt.Errorf(format, args...)
	}
}

func (r *reader) bytes(n int) []byte {
	if r.err != nil {
		return nil
	}
	if n < 0 || r.pos+n > len(r.b) {
		r.fail("truncated model: wanted %d bytes at offset %d, have %d",
			n, r.pos, len(r.b)-r.pos)
		return nil
	}
	out := r.b[r.pos : r.pos+n]
	r.pos += n
	return out
}

func (r *reader) u32() uint32 {
	b := r.bytes(4)
	if b == nil {
		return 0
	}
	return binary.LittleEndian.Uint32(b)
}

func (r *reader) u16() uint16 {
	b := r.bytes(2)
	if b == nil {
		return 0
	}
	return binary.LittleEndian.Uint16(b)
}

func (r *reader) u8() uint8 {
	b := r.bytes(1)
	if b == nil {
		return 0
	}
	return b[0]
}

func (r *reader) str16() string {
	n := int(r.u16())
	return string(r.bytes(n))
}

// u32s copies n little-endian uint32s. The copy is deliberate: the caller's
// backing array may be unaligned embedded data.
func (r *reader) u32s(n int) []uint32 {
	b := r.bytes(n * 4)
	if b == nil {
		return nil
	}
	out := make([]uint32, n)
	for i := range out {
		out[i] = binary.LittleEndian.Uint32(b[i*4:])
	}
	return out
}

func (r *reader) f32s(n int) []float32 {
	b := r.bytes(n * 4)
	if b == nil {
		return nil
	}
	out := make([]float32, n)
	for i := range out {
		out[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[i*4:]))
	}
	return out
}

// Unmarshal decodes a flat model.
//
// data is retained: feature strings are sliced from it rather than copied, so
// the caller must not modify it afterwards. Embedded data satisfies this.
func Unmarshal(data []byte) (*Model, error) {
	r := &reader{b: data}

	if string(r.bytes(len(Magic))) != Magic {
		return nil, errors.New("not a prose flat tagger model")
	}
	if v := r.u32(); v != Version {
		return nil, fmt.Errorf("unsupported model version %d (want %d)", v, Version)
	}

	m := &Model{}

	nClasses := int(r.u32())
	if r.err != nil {
		return nil, r.err
	}
	if nClasses > MaxClasses {
		return nil, fmt.Errorf("model has %d classes, format allows %d",
			nClasses, MaxClasses)
	}
	m.classes = make([]string, nClasses)
	for i := range m.classes {
		m.classes[i] = r.str16()
	}

	nFeatures := int(r.u32())
	nPairs := int(r.u32())
	blobLen := int(r.u32())
	if r.err != nil {
		return nil, r.err
	}

	blob := r.bytes(blobLen)
	if r.err != nil {
		return nil, r.err
	}
	// Zero-copy: the blob's bytes live in the caller's slice (rodata, for an
	// embedded model), so no 1.4 MB duplicate ends up on the heap.
	m.blob = unsafeString(blob)

	m.featOff = r.u32s(nFeatures + 1)
	m.pairOff = r.u32s(nFeatures + 1)

	if r.err != nil {
		return nil, r.err
	}
	ids := r.bytes(nPairs)
	if r.err != nil {
		return nil, r.err
	}
	m.classIDs = make([]uint8, nPairs)
	copy(m.classIDs, ids)

	m.weights = r.f32s(nPairs)

	nTagMap := int(r.u32())
	if r.err != nil {
		return nil, r.err
	}
	m.tagMap = make(map[string]uint8, nTagMap)
	for i := 0; i < nTagMap; i++ {
		w := r.str16()
		id := r.u8()
		if r.err != nil {
			return nil, r.err
		}
		if int(id) >= nClasses {
			return nil, fmt.Errorf("tagMap entry %q has class ID %d, only %d classes",
				w, id, nClasses)
		}
		m.tagMap[w] = id
	}
	if r.err != nil {
		return nil, r.err
	}

	if err := m.validate(); err != nil {
		return nil, err
	}
	return m, nil
}

// validate checks the invariants the scoring path relies on, so a corrupt
// model fails at load with a clear message instead of panicking mid-tag.
func (m *Model) validate() error {
	n := len(m.featOff) - 1
	if len(m.pairOff) != n+1 {
		return fmt.Errorf("offset arrays disagree: %d features vs %d pair offsets",
			n, len(m.pairOff)-1)
	}
	if len(m.classIDs) != len(m.weights) {
		return fmt.Errorf("classIDs (%d) and weights (%d) differ in length",
			len(m.classIDs), len(m.weights))
	}
	if n > 0 {
		if int(m.featOff[n]) != len(m.blob) {
			return fmt.Errorf("feature offsets end at %d, blob is %d bytes",
				m.featOff[n], len(m.blob))
		}
		if int(m.pairOff[n]) != len(m.weights) {
			return fmt.Errorf("pair offsets end at %d, weights has %d entries",
				m.pairOff[n], len(m.weights))
		}
	}
	for _, id := range m.classIDs {
		if int(id) >= len(m.classes) {
			return fmt.Errorf("class ID %d out of range (%d classes)",
				id, len(m.classes))
		}
	}
	// The binary search in feature() is only correct on sorted input.
	for i := 1; i < n; i++ {
		prev := m.blob[m.featOff[i-1]:m.featOff[i]]
		cur := m.blob[m.featOff[i]:m.featOff[i+1]]
		if prev >= cur {
			return fmt.Errorf("features not sorted at index %d: %q >= %q",
				i, prev, cur)
		}
	}
	return nil
}
