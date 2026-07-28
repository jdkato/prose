package re_test

import (
	"testing"

	"github.com/jdkato/prose/v3/internal/re"
)

// The one-letter Unicode class forms (\pL, \pN, \PL) failed to compile on the
// regexp2 fork prose used to depend on. They are valid in regexp/syntax, so
// callers reasonably expect them to work.
//
// See https://github.com/dlclark/regexp2/issues/65.
func TestUnicodeClassShorthand(t *testing.T) {
	cases := []struct {
		pattern string
		input   string
		want    bool
	}{
		{`\pL+`, "abc", true},
		{`\pN+`, "123", true},
		{`\PL+`, "123", true},
		{`\p{Greek}+`, "αβγ", true},
		{`\p{Han}+`, "漢字", true},
		{`\p{Cyrillic}+`, "привет", true},
		{`\pL+`, "123", false},
	}

	for _, c := range cases {
		r, err := re.Compile(c.pattern)
		if err != nil {
			t.Errorf("Compile(%q): %v", c.pattern, err)
			continue
		}
		if got := r.MatchString(c.input); got != c.want {
			t.Errorf("%q against %q: got %v, want %v", c.pattern, c.input, got, c.want)
		}
	}
}

// Lookarounds and backreferences are why prose uses regexp2 rather than the
// standard library for caller-supplied patterns.
func TestBeyondRE2Syntax(t *testing.T) {
	for _, pattern := range []string{`(?=foo)\w+`, `(?<=foo)bar`, `(\w)\1`} {
		if _, err := re.Compile(pattern); err != nil {
			t.Errorf("Compile(%q): %v", pattern, err)
		}
	}
}
