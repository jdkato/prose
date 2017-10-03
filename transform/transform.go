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

// Simple returns a space-separated, lower-cased copy of the string s.
func Simple(s string) string {
	return removeCase(s, " ")
}

// Dash returns a dash-separated, lower-cased copy of the string s.
func Dash(s string) string {
	return removeCase(s, "-")
}

// Snake returns a underscore-separated, lower-cased copy of the string s.
func Snake(s string) string {
	return removeCase(s, "_")
}

// Dot returns a underscore-separated, lower-cased copy of the string s.
func Dot(s string) string {
	return removeCase(s, ".")
}
