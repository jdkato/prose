package util

import "github.com/jdkato/prose/internal/model"

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

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

func GetAsset(name string) []byte {
	b, err := model.Asset("model/" + name)
	CheckError(err)
	return b
}
