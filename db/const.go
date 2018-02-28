package db

import (
	"github.com/juju/errors"
)

const DefaultPassword = "123456"

//角色
const (
	RoleAdmin   = 1
	RoleTeacher = 2
	RoleStudent = 3
)

//用户状态
const (
	UserStatsNormal = 1
	UserStatsForbid = 2
	UserStatsDelete = 3
)

//error
var (
	ErrNotFound = errors.New("not found")
)