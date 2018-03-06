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

func AddTeacher(ctx context.Context) {
	var teacher db.Teacher
	err := ctx.ReadJSON(&teacher)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	err = db.AddTeacher(&teacher)
	if err != nil {
		log.Errorf("[AddTeacher] add teacher(%v) error(%v)", goutil.Struct2Json(teacher), err)
		if strings.Contains(err.Error(), "already exists") {
			WriteResultErrByMsg(ctx, CodeAlreadyExists, "学号已存在", err)
			return
		}
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"teacher": teacher,
	})
}

func UpdateTeacher(ctx context.Context) {
	id := ctx.Params().Get("id")
	var teacher db.Teacher
	err := ctx.ReadJSON(&teacher)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	teacher.ID = id
	err = db.UpdateTeacher(&teacher)
	if err != nil {
		log.Errorf("[UpdateTeacher] update teacher(%v) error(%v)", goutil.Struct2Json(teacher), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"teacher": teacher,
	})
}

func DeleteTeacher(ctx context.Context) {
	var ids []string
	err := ctx.ReadJSON(&ids)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	err = db.UpdateTeacherByIDs(ids, goutil.Map{"status": db.UserStatsDelete})
	if err != nil {
		log.Errorf("[DeleteTeacher] delete ids(%v) error(%v)", ids, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"ids": ids,
	})
}

func GetTeacher(ctx context.Context) {
	id := ctx.Params().Get("id")
	teacher, err := db.LoadTeacher(id)
	if err != nil {
		log.Errorf("[GetTeacher] get teacher(%v) error(%v)", id, err)
		if err == db.ErrNotFound {
			WriteResultErrByMsg(ctx, CodeUserNotFound, "学号不存在", err)
			return
		}
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"teacher": teacher,
	})
}

func GetTeachers(ctx context.Context) {
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
	teacherList, total, err := db.ListTeacher(exactCondMap, fuzzyCondMap, sort, skip, limit)
	if err != nil {
		log.Errorf("[GetTeachers] get teacher by(%v) error(%v)", argMap, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	WriteResultSuccess(ctx, goutil.Map{
		"teacherList": teacherList,
		"total":       total,
	})
}

func CallGetTeachers(token string, argMap goutil.Map) (goutil.Map, error) {
	url := fmt.Sprintf("http://%v/teachers?", SrvAddr)
	result, err := tools.HttpGet(appendArgs(url, argMap), token)
	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}
