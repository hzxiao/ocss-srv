package api

import (
	//"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"github.com/hzxiao/ocss-srv/db"
	"testing"
	"github.com/hzxiao/goutil"
)


func TestAddCourse(t *testing.T) {
	removeAll()

	_, token, err := createUserAndLogin(db.RoleAdmin)
	assert.NoError(t, err)

	res, err := CallAddCourse(token, &db.Course{Name: "语文"})
	assert.NoError(t, err)
	t.Log(res)
}

func TestGetCourse(t *testing.T) {
	_, token, err := createUserAndLogin(db.RoleAdmin)
	assert.NoError(t, err)

	res, err := CallAddCourse(token, &db.Course{Name: "语文1"})
	assert.NoError(t, err)
	t.Log(res)

	res1, err := CallGetCourse(token, res.GetStringP("course/id"))
	assert.NoError(t, err)
	assert.Equal(t, res.GetStringP("course/name"), res1.GetStringP("course/name"))
}

func TestGetCourses(t *testing.T) {
	_, token, err := createUserAndLogin(db.RoleAdmin)
	assert.NoError(t, err)
	res, err := CallGetCourses(token, goutil.Map{"status": 1,"name":"语","pageSize":10,"page":0})
	assert.NoError(t, err)
	t.Log(res)

	assert.NotNil(t,res.Get("courseList"))
	t.Log(res)

}


func TestUpdateCourse(t *testing.T) {
	//removeAll()
	_, token, err := createUserAndLogin(db.RoleAdmin)
	assert.NoError(t, err)

	//c := &db.Course{
	//	Name: "语文",
	//}
	//course, err := CallAddCourse(token,c)
	//assert.NoError(t, err)
	//assert.NotNil(t,course)



	cu := &db.Course{
		ID:"5a9fe490a8c72e62387361de",
		Name: "语文1",
		Status:   db.CourseStatusChecked,
		Credit: "3",
		Period:"4",
		Attr:"基础",
		Nature:"必修",
		Campus:"大学城",
		Desc:"基础课程",
	}
	res, err := CallUpdateCourse(token,"5a9fe490a8c72e62387361de", cu)
	assert.NoError(t, err)

	ca,err:=CallGetCourse(token,res.GetStringP("course/id"))
	assert.NoError(t,err)
	assert.NotNil(t,ca.Get("course"))
	assert.Equal(t, cu.Name, ca.GetStringP("course/name"))
	assert.Equal(t, int64(cu.Status), ca.GetInt64P("course/status"))
	assert.Equal(t, cu.Credit, ca.GetStringP("course/credit"))
	assert.Equal(t, cu.Period, ca.GetStringP("course/period"))
	assert.Equal(t, cu.Attr, ca.GetStringP("course/attr"))
	assert.Equal(t, cu.Nature, ca.GetStringP("course/nature"))
	assert.Equal(t, cu.Campus, ca.GetStringP("course/campus"))
	assert.Equal(t, cu.Desc, ca.GetStringP("course/desc"))
}

func TestDeleteCourse(t *testing.T) {

}