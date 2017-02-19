package aptag

import "github.com/jdkato/aptag/model"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getAsset(name string) []byte {
	b, err := model.Asset("model/" + name)
	checkError(err)
	return b
}
