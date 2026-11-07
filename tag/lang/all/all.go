// Package all registers every part-of-speech model prose ships.
//
// Import it for its side effect when the language is decided at run time:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/all"
//
// A program that tags one language should import that language's package
// instead and leave the other models out of its binary; together they come
// to about 24 MB. See MODELS.md for what each one is.
package all

import (
	// Each import registers its language; nothing else is used.
	_ "github.com/jdkato/prose/v3/tag/lang/da"
	_ "github.com/jdkato/prose/v3/tag/lang/de"
	_ "github.com/jdkato/prose/v3/tag/lang/es"
	_ "github.com/jdkato/prose/v3/tag/lang/fi"
	_ "github.com/jdkato/prose/v3/tag/lang/fr"
	_ "github.com/jdkato/prose/v3/tag/lang/it"
	_ "github.com/jdkato/prose/v3/tag/lang/nb"
	_ "github.com/jdkato/prose/v3/tag/lang/nl"
	_ "github.com/jdkato/prose/v3/tag/lang/pt"
	_ "github.com/jdkato/prose/v3/tag/lang/sv"
	_ "github.com/jdkato/prose/v3/tag/lang/tr"
)
