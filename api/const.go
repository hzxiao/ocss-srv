package api

//result code
const (
	CodeSuccess = iota
	CodeArgErr
	CodeSrvErr
	CodeUserNotFound
	CodeForbid
	CodeDeleted
	CodeAlreadyExists
	CodeCourseNotFound
)

const (
	AdminCN = "管理员"
	TeacherCN = "教师"
	StudentCN = "学生"
)