/*
Package model contains internals used by prose/tag.
*/
package model

import "github.com/jdkato/prose/internal/util"

// GetAsset returns the named Asset.
func GetAsset(name string) []byte {
	b, err := Asset("internal/model/" + name)
	util.CheckError(err)
	return b
}
