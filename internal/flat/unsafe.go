package flat

import "unsafe"

// unsafeString views b as a string without copying.
//
// Safe here because the only caller hands it a slice of the model's backing
// bytes, which come from go:embed (read-only rodata) or a file read that the
// caller has contracted not to mutate. Copying instead would put a redundant
// 1.4 MB of feature strings on the heap for the GC to scan.
func unsafeString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}
