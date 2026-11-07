// Package all registers every sentence-segmentation model prose ships.
//
// Import it for its side effect when the language is decided at run time:
//
//	import _ "github.com/jdkato/prose/v3/segment/lang/all"
//
// A program that segments one language should import that language's package
// instead and leave the others out of its binary.
package all

import (
	// Each import registers its language; nothing else is used.
	_ "github.com/jdkato/prose/v3/segment/lang/da"
	_ "github.com/jdkato/prose/v3/segment/lang/de"
	_ "github.com/jdkato/prose/v3/segment/lang/es"
	_ "github.com/jdkato/prose/v3/segment/lang/fi"
	_ "github.com/jdkato/prose/v3/segment/lang/fr"
	_ "github.com/jdkato/prose/v3/segment/lang/it"
	_ "github.com/jdkato/prose/v3/segment/lang/nb"
	_ "github.com/jdkato/prose/v3/segment/lang/nl"
	_ "github.com/jdkato/prose/v3/segment/lang/pt"
	_ "github.com/jdkato/prose/v3/segment/lang/sv"
	_ "github.com/jdkato/prose/v3/segment/lang/tr"
)
