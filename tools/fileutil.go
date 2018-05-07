package tools

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"log"
)

func SaveFile(dir string, filename string, reader io.Reader) error {
	path := dir + string(filepath.Separator) + filename

	log.Printf("[SaveFile] save path(%v) ", path)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, os.ModePerm)
}

func UnmarshalJsonFile(filename string, dest interface{}) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	return json.Unmarshal(buf, &dest)
}

func MarshalJsonFile(filename string, data interface{})  error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(&data)
}