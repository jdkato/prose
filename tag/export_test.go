package tag_test

import "github.com/jdkato/prose/v3/tag/aptagmodel"

// embeddedModel gives benchmarks access to the raw model bytes so they can
// measure FromBytes without reaching into package internals.
var embeddedModel = aptagmodel.English
