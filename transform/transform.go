/*
Package transform implements functions to manipulate UTF-8 encoded strings.
*/
package transform

import (
	"regexp"
	"strings"
	"unicode"
)

var spaces = regexp.MustCompile(" +")

func removeCase(s string, sep string) string {
	out := ""
	old := ' '
	for i, c := range s {
		alpha := unicode.IsLetter(c) || unicode.IsNumber(c)
		mat := (i > 1 && unicode.IsLower(old) && unicode.IsUpper(c))
		if mat || !alpha || (unicode.IsSpace(c) && c != ' ') {
			out += " "
		}
		if alpha || c == ' ' {
			out += string(unicode.ToLower(c))
		}
		old = c
	}
	return spaces.ReplaceAllString(strings.TrimSpace(out), sep)
}

// SimpleCase returns a space-separated, lower-cased copy of the string s.
func SimpleCase(s string) string {
	return removeCase(s, " ")
}

// DashCase returns a dash-separated, lower-cased copy of the string s.
func DashCase(s string) string {
	return removeCase(s, "-")
}

// SnakeCase returns a underscore-separated, lower-cased copy of the string s.
func SnakeCase(s string) string {
	return removeCase(s, "_")
}
