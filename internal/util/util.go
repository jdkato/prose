package util

import (
	"io/ioutil"
	"path/filepath"

	"github.com/jdkato/prose/internal/model"
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

// GetAsset returns the named Asset.
func GetAsset(name string) []byte {
	b, err := model.Asset("model/" + name)
	CheckError(err)
	return b
}
