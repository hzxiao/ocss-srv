package api

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"github.com/hzxiao/ocss-srv/db"
	"testing"
)

func TestCallGetStudents(t *testing.T) {
	_, token, err := createUserAndLogin(db.RoleAdmin)
	assert.NoError(t, err)

	res, err := CallGetStudents(token, goutil.Map{
		"name": "xx",
		"id":   "",
	})
	assert.NoError(t, err)

	t.Log(res)
}
