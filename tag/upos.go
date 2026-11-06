package tag

// ToUniversal maps a Penn Treebank tag onto the Universal Dependencies tag
// set, so that a caller can read the English model's output in the same terms
// as every other language's.
//
// The mapping is the usual one and is lossy in the direction that matters
// least: Penn's inflectional distinctions (NN vs NNS, VB vs VBD) collapse into
// one universal class. Two of Penn's tags are genuinely ambiguous -- IN is
// both ADP and SCONJ, and TO is both PART and ADP -- and are mapped to the
// more common reading. Penn has no AUX, so an auxiliary verb comes back as
// VERB. An unknown tag maps to X.
func ToUniversal(penn string) string {
	if u, ok := pennToUniversal[penn]; ok {
		return u
	}
	return "X"
}

var pennToUniversal = map[string]string{
	"#": "SYM", "$": "SYM", "''": "PUNCT", "``": "PUNCT", ",": "PUNCT",
	"-LRB-": "PUNCT", "-RRB-": "PUNCT", "(": "PUNCT", ")": "PUNCT",
	".": "PUNCT", ":": "PUNCT", "-": "PUNCT", "HYPH": "PUNCT", "NFP": "PUNCT",
	"CC": "CCONJ", "CD": "NUM", "DT": "DET", "EX": "PRON", "FW": "X",
	"IN": "ADP", "JJ": "ADJ", "JJR": "ADJ", "JJS": "ADJ", "LS": "X",
	"MD": "AUX", "NN": "NOUN", "NNS": "NOUN", "NNP": "PROPN", "NNPS": "PROPN",
	"PDT": "DET", "POS": "PART", "PRP": "PRON", "PRP$": "PRON",
	"RB": "ADV", "RBR": "ADV", "RBS": "ADV", "RP": "ADP", "SYM": "SYM",
	"TO": "PART", "UH": "INTJ", "VB": "VERB", "VBD": "VERB", "VBG": "VERB",
	"VBN": "VERB", "VBP": "VERB", "VBZ": "VERB", "WDT": "DET", "WP": "PRON",
	"WP$": "PRON", "WRB": "ADV", "ADD": "X", "AFX": "ADJ", "GW": "X", "XX": "X",
	"-NONE-": "X",
}

// UniversalTags lists the Universal Dependencies tag set.
var UniversalTags = []string{
	"ADJ", "ADP", "ADV", "AUX", "CCONJ", "DET", "INTJ", "NOUN", "NUM",
	"PART", "PRON", "PROPN", "PUNCT", "SCONJ", "SYM", "VERB", "X",
}
