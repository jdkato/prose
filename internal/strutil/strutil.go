// Package strutil holds small string helpers shared across prose.
package strutil

import (
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// HasAnySuffix reports whether s ends with any of the given suffixes.
func HasAnySuffix(s string, suffixes []string) bool {
	for _, suf := range suffixes {
		if strings.HasSuffix(s, suf) {
			return true
		}
	}
	return false
}

// HasAnyPrefix reports whether s starts with any of the given prefixes.
func HasAnyPrefix(s string, prefixes []string) bool {
	for _, pre := range prefixes {
		if strings.HasPrefix(s, pre) {
			return true
		}
	}
	return false
}

// HasAnyIndex returns the index of the earliest of substrs found in s, or -1.
func HasAnyIndex(s string, substrs []string) int {
	best := -1
	for _, sub := range substrs {
		if i := strings.Index(s, sub); i >= 0 && (best == -1 || i < best) {
			best = i
		}
	}
	return best
}

// StringInSlice reports whether slice contains a.
func StringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if b == a {
			return true
		}
	}
	return false
}

// CharAt returns the byte at index i, or 0 if out of range.
func CharAt(s string, i int) byte {
	if i >= 0 && i < len(s) {
		return s[i]
	}
	return 0
}

// Min returns the smaller of a and b.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ToTitle uppercases the first rune of m. If strict, the remainder is
// lowercased.
func ToTitle(m string, strict bool) string {
	if m == "" {
		return m
	}
	r, size := utf8.DecodeRuneInString(m)

	other := m[size:]
	if strict {
		other = strings.ToLower(other)
	}

	return string(unicode.ToTitle(r)) + other
}

// ReadDataFile reads a test fixture, panicking on failure. Test-support only.
func ReadDataFile(path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}

// EqualFloat reports whether two floats agree when formatted to two decimal
// places. Test-support only.
//
// Readability metrics are reported to two decimals, so that is the precision
// the tests assert.
func EqualFloat(expected, observed float64) bool {
	return fmt.Sprintf("%0.2f", expected) == fmt.Sprintf("%0.2f", observed)
}
