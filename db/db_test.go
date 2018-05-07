package db

import (
	"github.com/hzxiao/goutil/assert"
	"github.com/hzxiao/ocss-srv/config"
	"testing"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/tools"
	"fmt"
	"strconv"
	"strings"
)

func init() {
	var err error
	err = config.InitConfig("config-test", "../config")
	if err != nil {
		panic(err)
	}
	err = InitDB("111.230.242.177:27017", "ocss-test")
	if err != nil {
		panic(err)
	}
}

func TestDB(t *testing.T) {
	C("test").RemoveAll(nil)

	err := C("test").Insert(map[string]string{
		"_id": "1",
	})
	assert.NoError(t, err)

	var res map[string]string
	err = C("test").FindId("1").One(&res)
	assert.NoError(t, err)
}

func removeAll() {
	C(CollectionUser).RemoveAll(nil)
	C(CollectionComment).RemoveAll(nil)
	C(CollectionCourse).RemoveAll(nil)
	C(CollectionTeacher).RemoveAll(nil)
	C(CollectionTeachCourse).RemoveAll(nil)
}

func TestAddComment2(t *testing.T) {
	var data goutil.Map
	var tc []goutil.Map
	err := tools.UnmarshalJsonFile("../data/courses.json", &data)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(data)
	tc = data.GetMapArray("data")
	var id = 11
	var dept []goutil.Map
	var m = map[string]string{}
	natures := map[string]bool{}
	attrs := map[string]bool{}
	for _, c := range tc {
		m[c.GetString("开课学院")] = ""
		natures[c.GetString("课程性质")] = true
		attrs[c.GetString("课程归属")] = true
	}

	var as []string
	for k := range attrs {
		as = append(as, k)
	}
	fmt.Println(goutil.Struct2Json(as))

	var ns []string
	for k := range natures {
		ns = append(ns, k)
	}
	fmt.Println(goutil.Struct2Json(ns))

	for k := range m {
		dept = append(dept, goutil.Map{
			"id": strconv.Itoa(id),
			"name": k,
		})
		m[k] = strconv.Itoa(id)
		id++
	}

	var users []*User
	var teachers []*Teacher
	var tCount = map[string]int{}
	var cCount = map[string]int{}
	var courses []*Course
	for _, c := range tc {
		te := &Teacher{
			ID: m[c.GetString("开课学院")] + fmt.Sprintf("%03d", tCount[m[c.GetString("开课学院")]]+1),
		}

		tCount[m[c.GetString("开课学院")]]++
		te.Name = c.GetString("教师姓名")
		te.Dept = goutil.Map{
			"id": m[c.GetString("开课学院")],
			"name": c.GetString("开课学院"),
		}
		te.Title = "讲师"
		users = append(users, &User{
			ID: te.ID,
			Username: te.ID,
			Role: RoleTeacher,
			Status:UserStatsNormal,
		})
		teachers = append(teachers, te)

		crs := &Course{
			ID: "18"+ m[c.GetString("开课学院")] + fmt.Sprintf("%03d", cCount[m[c.GetString("开课学院")]]+1),
			Name: c.GetString("课程名称"),
			Dept: goutil.Map{
				"id": m[c.GetString("开课学院")],
				"name": c.GetString("开课学院"),
			},
			Credit: c.GetString("学分"),
			Attr: c.GetString("课程归属"),
			Nature: c.GetString("课程性质"),
			Campus: c.GetString("学区代码"),
		}
		cCount[m[c.GetString("开课学院")]]++

		courses = append(courses, crs)
	}

	var tcs []*TeachCourse
	for _, c := range tc {

		tec := &TeachCourse{}
		for i := range courses {
			if c.GetString("课程名称") == courses[i].Name {
				tec.CID = courses[i].ID
			}
		}

		for i := range teachers {
			if c.GetString("教师姓名") == teachers[i].Name {
				tec.TID = teachers[i].ID
			}
		}

		tec.Addr = c.GetString("上课地点")
		tec.Capacity, _ = strconv.Atoi(c.GetString("容量"))

		weeks := strings.Split(c.GetString("起始结束周"), "-")
		tec.TakeWeek = goutil.Map{
			"startWeek": weeks[0],
			"endWeek": weeks[1],
		}

		times := []rune(c.GetString("上课时间"))
		//tec.TakeTime = goutil.Map{
		//	"dayOfWeek": string(times[1]),
		//	"sections": []string{"9", "10", "11"},
		//}
		tec.TakeTime = &TakeTime{
			DayOfWeek:[]string{string(times[1])},
			Sections: [][]int64{[]int64{10, 11}, },
		}
		tcs = append(tcs, tec)
	}
	err = tools.MarshalJsonFile("../data/users.json", users)
	if err != nil {
		t.Error(err)
		return
	}

	err = tools.MarshalJsonFile("../data/teachers.json", teachers)
	if err != nil {
		t.Error(err)
		return
	}

	err = tools.MarshalJsonFile("../data/course.json", courses)
	if err != nil {
		t.Error(err)
		return
	}


	err = tools.MarshalJsonFile("../data/tc.json", tcs)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(goutil.Struct2Json(dept))
}

func TestPrepareData(t *testing.T) {
	removeAll()
	err := PrepareData()
	if err != nil {
		t.Error(err)
		return
	}
}


func TestPrepareData2(t *testing.T) {
	removeAll()
	fileDir := "../data/"

	//var users []*User
	//err := tools.UnmarshalJsonFile(fileDir+"users.json", &users)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//
	//for i := range users {
	//	err = AddUser(users[i])
	//	if err != nil {
	//		t.Log(users[i].ID)
	//		t.Error(err)
	//		return
	//	}
	//}

	var courses []*Course
	err := tools.UnmarshalJsonFile(fileDir+"course.json", &courses)
	if err != nil {
		t.Error(err)
		return
	}

	for i := range courses {
		_, err = AddCourse(courses[i])
		if err != nil {
			t.Error(err)
			return
		}
	}

	var teachers []*Teacher
	err = tools.UnmarshalJsonFile(fileDir+"teachers.json", &teachers)
	if err != nil {
		t.Error(err)
		return
	}

	for i := range teachers {
		err = AddTeacher(teachers[i])
		if err != nil {
			t.Error(err)
			return
		}
	}

	var tcs []*TeachCourse
	err = tools.UnmarshalJsonFile(fileDir+"tc.json", &tcs)
	if err != nil {
		t.Error(err)
		return
	}

	for i := range tcs {
		err = AddTeachCourse(tcs[i])
		if err != nil {
			t.Error(err)
			return
		}
	}

}