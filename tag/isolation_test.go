package tag_test

import (
	"os/exec"
	"strings"
	"testing"
)

// prose's package layout rests on one property: a consumer that only tags text
// must not link any model it does not use. Go's linker drops unimported
// packages entirely, including their go:embed data, so this is achievable with
// package boundaries alone — but only if nothing in the tagger's import graph
// reaches a foreign model.
//
// A linter that wants POS tagging and nothing else should not pay for the
// entity model in binary size or startup. If this test fails, it does.
func TestTaggerDoesNotLinkForeignModels(t *testing.T) {
	out, err := exec.Command("go", "list", "-deps", ".").Output()
	if err != nil {
		t.Fatalf("go list -deps: %v", err)
	}

	// Packages the tagger must never reach. The root prose package pulls in
	// ner, and with it the ~4 MB entity model.
	forbidden := []string{
		"github.com/jdkato/prose/v3",
		"github.com/jdkato/prose/v3/ner",
		"github.com/jdkato/prose/v3/ner/maxentmodel",
	}

	deps := make(map[string]bool)
	for _, line := range strings.Split(string(out), "\n") {
		deps[strings.TrimSpace(line)] = true
	}

	for _, f := range forbidden {
		if deps[f] {
			t.Errorf("package tag depends on %s; a tagging-only consumer would "+
				"link that package's embedded model", f)
		}
	}

	// Guard against the reverse mistake too: the model must be a leaf, so that
	// importing it never drags in unrelated code.
	out, err = exec.Command("go", "list", "-deps", "./aptagmodel").Output()
	if err != nil {
		t.Fatalf("go list -deps ./aptagmodel: %v", err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		pkg := strings.TrimSpace(line)
		if strings.HasPrefix(pkg, "github.com/jdkato/prose/") &&
			pkg != "github.com/jdkato/prose/v3/tag/aptagmodel" {
			t.Errorf("aptagmodel should be a dependency-free data leaf, but "+
				"it imports %s", pkg)
		}
	}
}
