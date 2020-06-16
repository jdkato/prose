package prose

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelFromDisk(t *testing.T) {
	data := filepath.Join(testdata, "PRODUCT")
	model := ModelFromDisk(data)
	assert.Equal(t, model.Name, "PRODUCT")

	temp := filepath.Join(testdata, "temp")
	_ = os.RemoveAll(temp)

	err := model.Write(temp)
	if err != nil {
		panic(err)
	}
	model = ModelFromDisk(temp)

	assert.Equal(t, model.Name, "temp")
}
