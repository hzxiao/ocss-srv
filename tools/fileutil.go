package tools

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SaveFile(dir string, filename string, reader io.Reader) error {
	path := dir + string(filepath.Separator) + filename

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, os.ModePerm)
}
