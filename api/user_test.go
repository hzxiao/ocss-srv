package api

import (
	"github.com/hzxiao/goutil/assert"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/db"
	"testing"
)

func TestLogin(t *testing.T) {
	removeAll()

	err := db.AddUser(&db.User{Username: "aaa"})
	assert.NoError(t, err)

	res, err := CallLogin(goutil.Map{
		"username": "aaa",
		"password": "123456",
	})
	assert.NoError(t, err)

	_ = res
}

func TestAddUser(t *testing.T) {
	removeAll()

	_, token, err := createUserAndLogin(db.RoleAdmin)
	assert.NoError(t, err)

	res, err := CallAddUser(token, &db.User{Username: "u1"})
	assert.NoError(t, err)
	t.Log(res)
}