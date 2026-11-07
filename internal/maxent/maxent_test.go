package maxent

import (
	"testing"

	"github.com/jdkato/prose/v3/ner/maxentmodel"
)

// Vector must agree with Weight for every key, including ones whose value
// contains a dash or extends another value.
func TestVectorMatchesWeight(t *testing.T) {
	m, err := Unmarshal(maxentmodel.English)
	if err != nil {
		t.Fatal(err)
	}
	n := m.NumKeys()
	if n == 0 {
		t.Fatal("empty model")
	}
	vec := make([]float32, len(m.labels))
	checked := 0
	for i := 0; i < n; i++ {
		key := m.blob[m.keyOff[i]:m.keyOff[i+1]]
		for j, label := range m.labels {
			if len(key) <= len(label)+1 || key[len(key)-len(label):] != label || key[len(key)-len(label)-1] != '-' {
				continue
			}
			prefix := key[:len(key)-len(label)]
			m.Vector(prefix, vec)
			w, _ := m.Weight(key)
			if vec[j] != w {
				t.Fatalf("%q label %q: vector %v, weight %v", prefix, label, vec[j], w)
			}
			checked++
		}
	}
	if checked < n {
		t.Fatalf("checked %d of %d keys", checked, n)
	}
	m.Vector("no-such-feature-", vec)
	for j, w := range vec {
		if w != 0 {
			t.Errorf("unknown key has weight %v for %s", w, m.labels[j])
		}
	}
}
