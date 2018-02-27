package db

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func init() {
	err := InitDB("111.230.242.177:27017", "ocss_test")
	if err != nil {
		panic(err)
	}
}

func TestDB(t *testing.T) {
	C("test").RemoveAll(nil)

	err := C("test").Insert(map[string]string{
		"_id": "1",
	})
	assert.NoError(t, err)

	var res map[string]string
	err = C("test").FindId("1").One(&res)
	assert.NoError(t, err)
}

func removeAll() {
	C(CollectionUser).RemoveAll(nil)
}