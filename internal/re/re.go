// Package re adapts github.com/dlclark/regexp2 to the parts of the standard
// library's regexp API that prose needs.
//
// prose uses regexp2 where patterns come from callers — vocabulary entries and
// case-conversion regexes — because regexp2 accepts lookarounds and
// backreferences that regexp/syntax rejects outright. Everything with a fixed,
// in-package pattern uses the standard library directly.
//
// regexp2's matching methods return an error, most often a timeout on a
// pathological pattern. The wrappers here treat that as "no match": a single
// bad vocabulary entry should not abort a conversion, and prose has no way to
// surface the error from the middle of one.
package re

import "github.com/dlclark/regexp2/v2"

// Regexp is a compiled pattern.
type Regexp struct {
	inner *regexp2.Regexp
}

// Compile parses a pattern in RE2 compatibility mode.
func Compile(pattern string) (*Regexp, error) {
	inner, err := regexp2.Compile(pattern, regexp2.RE2)
	if err != nil {
		return nil, err
	}
	return &Regexp{inner: inner}, nil
}

// MustCompile is Compile but panics on an invalid pattern. It is meant for
// patterns fixed at compile time.
func MustCompile(pattern string) *Regexp {
	r, err := Compile(pattern)
	if err != nil {
		panic("prose/internal/re: " + err.Error())
	}
	return r
}

// Unwrap returns the underlying regexp2 pattern.
func (r *Regexp) Unwrap() *regexp2.Regexp { return r.inner }

// MatchString reports whether s contains a match.
func (r *Regexp) MatchString(s string) bool {
	ok, err := r.inner.MatchString(s)
	return err == nil && ok
}

// FindAllString returns up to n successive matches, or all of them if n < 0.
func (r *Regexp) FindAllString(s string, n int) []string {
	if n == 0 {
		return nil
	}

	var out []string
	m, err := r.inner.FindStringMatch(s)
	for err == nil && m != nil {
		out = append(out, m.String())
		if n > 0 && len(out) >= n {
			break
		}
		m, err = r.inner.FindNextMatch(m)
	}
	return out
}

// MatchString reports whether s matches pattern. An invalid pattern does not
// match rather than panicking.
func MatchString(pattern, s string) bool {
	r, err := Compile(pattern)
	if err != nil {
		return false
	}
	return r.MatchString(s)
}
