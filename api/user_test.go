package api

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
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

func TestGetUser(t *testing.T) {
	username, token, err := createUserAndLogin(db.RoleAdmin)
	assert.NoError(t, err)

	res, err := CallGetUser(token, username)
	assert.NoError(t, err)
	assert.Equal(t, username, res.GetStringP("user/username"))
}

func TestUpdateUser(t *testing.T) {
	username, token, err := createUserAndLogin(db.RoleAdmin)
	assert.NoError(t, err)

	u := &db.User{
		Username:username,
		Icon: "cc",
	}
	res, err := CallUpdateUser(token, username, u)
	assert.NoError(t, err)
	assert.Equal(t, u.Icon, res.GetStringP("user/icon"))
}
