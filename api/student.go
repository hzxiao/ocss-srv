package api

import (
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/kataras/iris/context"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func AddStudent(ctx context.Context) {
	var stu db.Student
	err := ctx.ReadJSON(&stu)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	err = db.AddStudent(&stu)
	if err != nil {
		log.Errorf("[AddStudent] add stu(%v) error(%v)", goutil.Struct2Json(stu), err)
		if strings.Contains(err.Error(), "already exists") {
			WriteResultErrByMsg(ctx, CodeAlreadyExists, "学号已存在", err)
			return
		}
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"student": stu,
	})
}

func UpdateStudent(ctx context.Context) {
	id := ctx.Params().Get("id")
	var stu db.Student
	err := ctx.ReadJSON(&stu)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	stu.ID = id
	err = db.UpdateStudent(&stu)
	if err != nil {
		log.Errorf("[UpdateStudent] update stu(%v) error(%v)", goutil.Struct2Json(stu), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"student": stu,
	})
}

func DeleteStudent(ctx context.Context) {
	var ids []string
	err := ctx.ReadJSON(&ids)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	err = db.UpdateStudentByIDs(ids, goutil.Map{"status": db.UserStatsDelete})
	if err != nil {
		log.Errorf("[DeleteStudent] delete ids(%v) error(%v)", ids, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"ids": ids,
	})
}

func GetStudent(ctx context.Context) {
	id := ctx.Params().Get("id")
	stu, err := db.LoadStudent(id)
	if err != nil {
		log.Errorf("[GetStudent] get stu(%v) error(%v)", id, err)
		if err == db.ErrNotFound {
			WriteResultErrByMsg(ctx, CodeUserNotFound, "学号不存在", err)
			return
		}
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"student": stu,
	})
}

func GetStudents(ctx context.Context) {
	argMap, err := CheckURLArg(ctx.FormValues(), []*Arg{
		{Key: "name", Type: "string"},
		{Key: "id", Type: "string"},
		{Key: "deptId", Type: "string"},
		{Key: "majorId", Type: "string"},
		{Key: "status", Type: "int", DefaultValue: strconv.Itoa(db.UserStatsNormal)},
		{Key: "page", Type: "int", DefaultValue: "1"},
		{Key: "pageSize", Type: "int", DefaultValue: "20"},
		{Key: "sort", Type: "string"},
	})
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}
	//handle args
	exactCondMap, fuzzyCondMap := goutil.Map{}, goutil.Map{}
	if argMap.Exist("name") {
		fuzzyCondMap.Set("name", argMap.Get("name"))
	}
	if argMap.Exist("id") {
		fuzzyCondMap.Set("id", argMap.Get("id"))
	}
	if argMap.Exist("deptId") {
		exactCondMap.Set("dept.id", argMap.Get("deptId"))
	}
	if argMap.Exist("majorId") {
		exactCondMap.Set("major.id", argMap.Get("majorId"))
	}
	exactCondMap.Set("status", int(argMap.GetInt64("status")))
	limit := int(argMap.GetInt64("pageSize"))
	skip := int(argMap.GetInt64("page")) * limit
	var sort []string
	if argMap.Exist("sort") {
		sort = append(sort, argMap.GetString("sort"))
	}
	studentList, total, err := db.ListStudent(exactCondMap, fuzzyCondMap, sort, skip, limit)
	if err != nil {
		log.Errorf("[GetStudents] get stu by(%v) error(%v)", argMap, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	WriteResultSuccess(ctx, goutil.Map{
		"studentList": studentList,
		"total":       total,
	})
}

func CallGetStudents(token string, argMap goutil.Map) (goutil.Map, error) {
	url := fmt.Sprintf("http://%v/students?", SrvAddr)
	result, err := tools.HttpGet(appendArgs(url, argMap), token)
	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}

func Count(ctx context.Context) {
	var info goutil.Map
	err := ctx.ReadJSON(&info)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	result := goutil.Map{}
	for k := range info {
		var count int
		switch k {
		case "student":
			cond := info.GetMap(k)
			tools.ReplaceKeys(cond, goutil.Map{
				"deptId":  "dept.id",
				"majorId": "major.id",
			})
			count, err = db.CountStudent(cond)
		case "course":
			cond := info.GetMap(k)
			tools.ReplaceKeys(cond, goutil.Map{
				"deptId": "dept.id",
			})
			count, err = db.CountCourse(cond)
		case "teacher":
			cond := info.GetMap(k)
			tools.ReplaceKeys(cond, goutil.Map{
				"deptId": "dept.id",
			})
			count, err = db.CountTeacher(cond)
		case "teachCourse":
			cond := info.GetMap(k)
			count, err = db.CountTeachCourse(cond)
		}
		if err != nil {
			log.Errorf("[Count] by key(%v), cond(%v) error(%v)", k, info.GetMap(k), err)
			WriteResultWithSrvErr(ctx, err)
			return
		}
		result.Set(k, count)
	}

	WriteResultSuccess(ctx, result)
}
