package db

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestAddUser(t *testing.T) {
	removeAll()

	var err error
	//1 add user1
	user1 := &User{
		Username: "user1",
	}
	err = AddUser(user1)
	assert.NoError(t, err)

	//2. add user2
	user2 := &User{
		Username: "user2",
		Role:     RoleTeacher,
	}
	err = AddUser(user2)
	assert.NoError(t, err)
	assert.Equal(t, RoleTeacher, user2.Role)

	//test err
	//3. add user2 again
	err = AddUser(&User{Username: "user2"})
	assert.Error(t, err)

	//4. nil user
	err = AddUser(nil)
	assert.Error(t, err)

	//5. empty username
	err = AddUser(&User{})
	assert.Error(t, err)
}

func TestUpdateUser(t *testing.T) {
	removeAll()

	var err error
	user := &User{
		Username: "user",
	}
	err = AddUser(user)
	assert.NoError(t, err)

	u := &User{
		Username: user.Username,
		Status:   UserStatsNormal,
		Role:     RoleTeacher,
		Icon:     "icon",
		Password: "654321",
	}

	err = UpdateUser(u)
	assert.NoError(t, err)

	verifyUser, err := VerifyUser(u.Username, u.Password)
	assert.NoError(t, err)
	assert.NotNil(t, verifyUser)

	assert.Equal(t, u.Status, verifyUser.Status)
	assert.Equal(t, u.Role, verifyUser.Role)
	assert.Equal(t, u.Icon, verifyUser.Icon)
}

func TestVerifyUser(t *testing.T) {
	removeAll()

	var err error
	user := &User{
		Username: "user",
	}
	err = AddUser(user)
	assert.NoError(t, err)

	u, err := VerifyUser(user.Username, DefaultPassword)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
