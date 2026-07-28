package ner

import "unsafe"

// unsafeString views b as a string without copying. Used only for model
// lookups that do not retain the result.
func unsafeString(b []byte) string {
	return unsafe.String(&b[0], len(b))
}
