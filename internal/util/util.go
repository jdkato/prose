/*
Package util contains internals used across the other prose packages.
*/
package util

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ReadDataFile reads data from a file, panicking on any errors.
func ReadDataFile(path string) []byte {
	p, err := filepath.Abs(path)
	CheckError(err)

	data, ferr := ioutil.ReadFile(p)
	CheckError(ferr)

	return data
}

// CheckError panics if `err` is not `nil`.
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// Min returns the minimum of `a` and `b`.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// IsPunct determines if a character is a punctuation symbol.
func IsPunct(c byte) bool {
	for _, r := range []byte("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") {
		if c == r {
			return true
		}
	}
	return false
}

// IsSpace determines if a character is a whitespace character.
func IsSpace(c byte) bool {
	for _, r := range []byte("\t\n\r\f\v") {
		if c == r {
			return true
		}
	}
	return false
}

// IsLetter determines if a character is letter.
func IsLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// IsAlnum determines if a character is a letter or a digit.
func IsAlnum(c byte) bool {
	return (c >= '0' && c <= '9') || IsLetter(c)
}

// StringInSlice determines if `slice` contains the string `a`.
func StringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if a == b {
			return true
		}
	}
	return false
}

// HasAnySuffix determines if the string a has any suffixes contained in the
// slice b.
func HasAnySuffix(a string, slice []string) bool {
	for _, b := range slice {
		if strings.HasSuffix(a, b) {
			return true
		}
	}
	return false
}

// ContainsAny determines if the string a contains any fo the strings contained
// in the slice b.
func ContainsAny(a string, b []string) bool {
	for _, s := range b {
		if strings.Contains(a, s) {
			return true
		}
	}
	return false
}

// SliceToMap takes in a slice of strings.  It then converts the slice into a map
// where the keys of the map is the strings in the slice and the value of the map
// is an empty struct.  This map is useful for doing repeated contains searches
// on the input slice.
func SliceToMap(slice []string) map[string]struct{} {
	// The value of this map is an empty struct because it will have a width of zero
	// and therefore takes up no memory.
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	return set
}
