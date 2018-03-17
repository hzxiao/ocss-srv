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

//课程状态
const (
	CourseStatusChecking = 1
	CourseStatusChecked  = 2
	CourseStatusDelete   = 3
)

//学生选课状态
const (
	SCourseStatusSelecting = 1
	SCourseStatusSelected = 2
	SCourseStatusGrade = 3
	SCourseStatusOver = 4
	SCourseStatusForbid = 5
)

//error
var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
)

//
const (
	StatusNormal  = 1
	StatusDeleted = 3
)

const (
	NoticeStateUnRead = 1
	NoticeStatusRead  = 2
)

//选课的状态
const (
	TeachCourseStatusSelectable    = 1
	TeachCourseStatusNonselectable = 2
	TeachCourseStatusLearning = 3
	TeachCourseStatusOver = 5

)
