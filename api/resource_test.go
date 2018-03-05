package api

import (
	"github.com/hzxiao/goutil/assert"
	"github.com/hzxiao/ocss-srv/config"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadFile(t *testing.T) {
	res, err := CallUploadFile("./resource.go")
	assert.NoError(t, err)

	t.Log("file id: ", res.GetStringP("file/id"))

	err = os.Remove(config.GetString("file.location") + string(filepath.Separator) + res.GetStringP("file/id"))
	assert.NoError(t, err)
}
