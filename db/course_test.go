package db

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
	"github.com/hzxiao/goutil"
)

func TestAddCourse(t *testing.T) {
	removeAll()

	var err error
	//1 add user1
	course1 := &Course{
		Name: "yuwen",
	}
	course, err := AddCourse(course1)
	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, course1.Name, course)

	//test err

	//4. nil user
	_, err = AddCourse(nil)
	assert.Error(t, err)

	//5. empty username
	_, err = AddCourse(&Course{})
	assert.Error(t, err)
}

func TestUpdateCourse(t *testing.T) {
	removeAll()

	var err error
	c := &Course{
		Name: "语文",
	}
	course, err := AddCourse(c)
	assert.NoError(t, err)
	assert.NotNil(t,course)

	cu := &Course{
		ID:course.ID,
		Name: "语文1",
		Status:   CourseStatusChecked,
		Credit: "3",
		Period:"4",
		Attr:"基础",
		Nature:"必修",
		Campus:"大学城",
		Desc:"基础课程",
	}

	err = UpdateCourse(cu)
	assert.NoError(t, err)

	ca, err := LoadCourse(course.ID)
	assert.NoError(t, err)
	assert.NotNil(t, ca)

	assert.Equal(t, ca.Name, cu.Name)
	assert.Equal(t, ca.Status, cu.Status)
	assert.Equal(t, ca.Credit, cu.Credit)
	assert.Equal(t, ca.Period, cu.Period)
	assert.Equal(t, ca.Attr, cu.Attr)
	assert.Equal(t, ca.Nature, cu.Nature)
	assert.Equal(t, ca.Campus, cu.Campus)
	assert.Equal(t, ca.Desc, cu.Desc)

	//test user not exists
	u1 := &Course{
		ID:"12345",
	}

	err = UpdateCourse(u1)
	assert.Error(t, err)
}

func TestListCourse(t *testing.T) {
	var err error
	exactCondMap, fuzzyCondMap := goutil.Map{}, goutil.Map{}
	exactCondMap.Set("status", 1)
	fuzzyCondMap.Set("name", "语")
	var sort []string
	sort = append(sort, "name")
	courses, num, err := ListCourse(exactCondMap, fuzzyCondMap, sort, 0, 5)
	assert.NoError(t, err)
	assert.NotNil(t, courses)
	assert.Equal(t, exactCondMap.GetInt64("status"), int64(courses[0].Status))
	assert.Equal(t, 1, num)

	course, err := LoadCourse(courses[0].ID)
	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, exactCondMap.GetInt64("status"), int64(course.Status))
}

