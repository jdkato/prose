package prose

import (
	"os"
	"path/filepath"
	"testing"
)

func TestModelFromDisk(t *testing.T) {
	data := filepath.Join(testdata, "PRODUCT")

	model := ModelFromDisk(data)
	if model.Name != "PRODUCT" {
		t.Errorf("ModelFromDisk() expected = PRODUCT, got = %v", model.Name)
	}

	temp := filepath.Join(testdata, "temp")
	_ = os.RemoveAll(temp)

	err := model.Write(temp)
	if err != nil {
		panic(err)
	}
	model = ModelFromDisk(temp)
	if model.Name != "temp" {
		t.Errorf("ModelFromDisk() expected = temp, got = %v", model.Name)
	}
}
