// Package maxent defines the on-disk format for prose's binary maximum-entropy
// entity model, plus the reader for it.
//
// # Why not gob
//
// A map[string]int of 149,576 entries plus a []float64 of 149,589 costs ~58 ms,
// ~38.7 MB and ~304,000 allocations to build.
//
// # Layout
//
// Inference only ever asks one question: "what is the weight for this feature
// string?" Answering it in two hops — string to dense index, index to weight —
// keeps an index that is only an artifact of how the model was trained. The
// flat format fuses the two: keys are sorted, and each key's weight sits at
// the same position. One binary search, one load.
//
//	magic     [8]byte  "PROSENER"
//	version   uint32
//	cardinality uint32
//	nLabels   uint32   then per label: uint16 len + bytes
//	nKeys     uint32
//	blobLen   uint32   then blobLen bytes (feature keys, sorted bytewise)
//	keyOff    (nKeys+1) x uint32
//	weights   nKeys x float32
//	nExtra    uint32   then nExtra x float32
//
// The "extra" weights are the tail of the weight vector beyond the key count:
// GIS correction terms, which classification does not use but the
// probability-estimating path does. They are preserved so the format is
// lossless.
//
// float32 is safe here: the largest relative error from narrowing the shipped
// model's weights is 5.9e-08, and label scores are sums of a dozen or so of
// them, far from any decision boundary.
package maxent

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sort"
	"unsafe"
)

// Magic identifies a flat entity model file.
const Magic = "PROSENER"

// Version is the current format version.
const Version = 1

// Model is a flat maximum-entropy model, ready for scoring.
type Model struct {
	cardinality int
	labels      []string

	blob    string   // feature keys, concatenated, sorted bytewise
	keyOff  []uint32 // nKeys+1 offsets into blob
	weights []float32

	extra []float32
}

// Labels returns the model's label set, in its original order.
func (m *Model) Labels() []string { return m.labels }

// Cardinality returns the model's cardinality, used by the GIS encoding.
func (m *Model) Cardinality() int { return m.cardinality }

// NumKeys reports how many feature keys the model knows.
func (m *Model) NumKeys() int { return len(m.keyOff) - 1 }

// Weight returns the weight for a feature key, and whether it was known.
func (m *Model) Weight(key string) (float32, bool) {
	n := len(m.keyOff) - 1
	i := sort.Search(n, func(i int) bool {
		return m.blob[m.keyOff[i]:m.keyOff[i+1]] >= key
	})
	if i < n && m.blob[m.keyOff[i]:m.keyOff[i+1]] == key {
		return m.weights[i], true
	}
	return 0, false
}

// ExtraWeight returns the i'th GIS correction weight.
func (m *Model) ExtraWeight(i int) (float32, bool) {
	if i < 0 || i >= len(m.extra) {
		return 0, false
	}
	return m.extra[i], true
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

func (r *reader) str16() string {
	b := r.bytes(2)
	if b == nil {
		return ""
	}
	return string(r.bytes(int(binary.LittleEndian.Uint16(b))))
}

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

// Unmarshal decodes a flat entity model.
//
// data is retained: feature keys are sliced from it rather than copied, so the
// caller must not modify it afterwards. Embedded data satisfies this.
func Unmarshal(data []byte) (*Model, error) {
	r := &reader{b: data}

	if string(r.bytes(len(Magic))) != Magic {
		return nil, errors.New("not a prose flat entity model")
	}
	if v := r.u32(); v != Version {
		return nil, fmt.Errorf("unsupported model version %d (want %d)", v, Version)
	}

	m := &Model{}
	m.cardinality = int(r.u32())

	nLabels := int(r.u32())
	if r.err != nil {
		return nil, r.err
	}
	m.labels = make([]string, nLabels)
	for i := range m.labels {
		m.labels[i] = r.str16()
	}

	nKeys := int(r.u32())
	blobLen := int(r.u32())
	if r.err != nil {
		return nil, r.err
	}

	blob := r.bytes(blobLen)
	if r.err != nil {
		return nil, r.err
	}
	m.blob = unsafeString(blob)

	m.keyOff = r.u32s(nKeys + 1)
	m.weights = r.f32s(nKeys)

	nExtra := int(r.u32())
	if r.err != nil {
		return nil, r.err
	}
	m.extra = r.f32s(nExtra)

	if r.err != nil {
		return nil, r.err
	}
	if err := m.validate(); err != nil {
		return nil, err
	}
	return m, nil
}

// validate checks the invariants Weight relies on, so a corrupt model fails at
// load rather than mid-classification.
func (m *Model) validate() error {
	n := len(m.keyOff) - 1
	if len(m.weights) != n {
		return fmt.Errorf("%d keys but %d weights", n, len(m.weights))
	}
	if n > 0 && int(m.keyOff[n]) != len(m.blob) {
		return fmt.Errorf("key offsets end at %d, blob is %d bytes",
			m.keyOff[n], len(m.blob))
	}
	// Binary search is only correct on sorted input.
	for i := 1; i < n; i++ {
		prev := m.blob[m.keyOff[i-1]:m.keyOff[i]]
		cur := m.blob[m.keyOff[i]:m.keyOff[i+1]]
		if prev >= cur {
			return fmt.Errorf("keys not sorted at index %d: %q >= %q", i, prev, cur)
		}
	}
	return nil
}

// unsafeString views b as a string without copying. Safe because the only
// caller hands it a slice of the model's backing bytes, which come from
// an embedded asset (read-only) or a file the caller contracted not to mutate.
func unsafeString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}
