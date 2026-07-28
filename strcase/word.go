package strcase

import (
	"regexp"
	"strings"
	"unicode"
)

var spaces = regexp.MustCompile(" +")

func removeCase(s string, sep string, t func(rune) rune) string {
	out := ""
	old := ' '
	for i, c := range s {
		alpha := unicode.IsLetter(c) || unicode.IsNumber(c)
		mat := (i > 1 && unicode.IsLower(old) && unicode.IsUpper(c))
		if mat || !alpha || (unicode.IsSpace(c) && c != ' ') {
			out += " "
		}
		if alpha || c == ' ' {
			out += string(t(c))
		}
		old = c
	}
	return spaces.ReplaceAllString(strings.TrimSpace(out), sep)
}

// Simple converts s to space-separated lower case.
func Simple(s string) string {
	return removeCase(s, " ", unicode.ToLower)
}

// Dash converts s to dash-separated lower case, as in "kebab-case".
func Dash(s string) string {
	return removeCase(s, "-", unicode.ToLower)
}

// Snake converts s to underscore-separated lower case.
func Snake(s string) string {
	return removeCase(s, "_", unicode.ToLower)
}

// Dot converts s to dot-separated lower case.
func Dot(s string) string {
	return removeCase(s, ".", unicode.ToLower)
}

// Constant converts s to underscore-separated upper case.
func Constant(s string) string {
	return removeCase(s, "_", unicode.ToUpper)
}

// Pascal converts s to PascalCase.
func Pascal(s string) string {
	out := ""
	wasSpace := false
	for i, c := range removeCase(s, " ", unicode.ToLower) {
		if i == 0 || wasSpace {
			c = unicode.ToUpper(c)
		}
		wasSpace = c == ' '
		if !wasSpace {
			out += string(c)
		}
	}
	return out
}

// Camel converts s to camelCase.
func Camel(s string) string {
	first := ' '
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			first = c
			break
		}
	}
	body := Pascal(s)
	if len(body) > 1 {
		return strings.TrimSpace(string(unicode.ToLower(first)) + body[1:])
	}
	return s
}
