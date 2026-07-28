package strcase

import "github.com/jdkato/prose/v3/internal/re"

// A RegexConverter converts a string to a regex-specified case.
type RegexConverter struct {
	CaseOpts
	Pattern *re.Regexp
}

// NewRegexConverter returns a new RegexConverter.
func NewRegexConverter(regex string, opts ...CaseOptFunc) (*RegexConverter, error) {
	pattern, err := re.Compile(regex)
	if err != nil {
		return nil, err
	}
	re := &RegexConverter{Pattern: pattern}

	base := defaultSentOpts
	for _, opt := range opts {
		opt(&base)
	}

	re.vocab = base.vocab
	re.indicator = base.indicator
	re.prefix = base.prefix

	return re, nil
}
