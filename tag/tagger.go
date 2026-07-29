package tag

import (
	"fmt"
	"sort"
	"sync"
)

// Interface is what every tagger here satisfies.
//
// prose ships more than one, because they disagree in ways that matter to the
// caller rather than in ways one of them is simply better at. The perceptron
// reads context and is the more accurate on ordinary prose; a dictionary-driven
// tagger is steadier on text that is already ungrammatical, which is where a
// grammar checker does its work. A rule written against one of them expects
// that one's answers, so the choice belongs to the caller.
type Interface interface {
	// Tag assigns a tag to each word and returns the tagged tokens.
	//
	// The input is treated as one sentence.
	Tag(words []string) []Token

	// TagTokens fills in the Tag field of each token, using its Text. This is
	// the low-allocation entry point: it writes into the caller's slice.
	TagTokens(tokens []Token)

	// Name reports which tagger this is, as registered.
	Name() string
}

// A Factory builds a Tagger. Registered rather than constructed directly so a
// caller can name one in configuration.
type Factory func() (Interface, error)

var (
	registryMu sync.RWMutex
	registry   = map[string]Factory{}
)

// Register makes a tagger available by name.
//
// Registering the same name twice panics: it means two packages disagree about
// what that name refers to, and picking one silently would make the choice
// depend on import order.
func Register(name string, f Factory) {
	registryMu.Lock()
	defer registryMu.Unlock()

	if _, taken := registry[name]; taken {
		panic("tag: " + name + " is already registered")
	}
	registry[name] = f
}

// Open builds the tagger registered under name.
//
// An empty name gives the default, which is the perceptron.
func Open(name string) (Interface, error) {
	if name == "" {
		name = Default
	}

	registryMu.RLock()
	f, ok := registry[name]
	registryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("tag: no tagger named %q; have %v", name, Names())
	}

	return f()
}

// Names lists the registered taggers, sorted.
func Names() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	out := make([]string, 0, len(registry))
	for name := range registry {
		out = append(out, name)
	}
	sort.Strings(out)

	return out
}

// Default names the tagger used when a caller does not choose.
const Default = "perceptron"

func init() {
	Register(Default, func() (Interface, error) { return New() })
}
