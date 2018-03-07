package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestAddComment(t *testing.T) {
	removeAll()

	c := &Comment{
		TCID:    "tcid",
		UID:     "uid",
		Content: "cc",
	}
	err := AddComment(c)
	assert.NoError(t, err)

	c1, err := UpdateChildComment(c.ID, "add", goutil.Map{"content": "xx"})
	assert.NoError(t, err)
	assert.Len(t, c1.Children, 1)

	c2, err := UpdateChildComment(c.ID, "del", c1.Children[0])
	assert.NoError(t, err)
	assert.Len(t, c2.Children, 0)
}
