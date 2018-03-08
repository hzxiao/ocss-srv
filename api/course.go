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
	"bytes"
)

func AddCourse(ctx context.Context) {
	var crs db.Course
	err := ctx.ReadJSON(&crs)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	course,err := db.AddCourse(&crs)
	if err != nil {
		log.Printf("[AddCourse] add crs(%v) error(%v)", goutil.Struct2Json(crs), err)
		if strings.Contains(err.Error(), "already exists") {
			WriteResultErrByMsg(ctx, CodeAlreadyExists, "课程已存在", err)
			return
		}
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"course": course,
	})
}

func CallAddCourse(token string, course *db.Course) (goutil.Map, error) {
	result, err := tools.HttpPost(fmt.Sprintf("http://%v/courses", SrvAddr), token,
		"application/json", bytes.NewBufferString(goutil.Struct2Json(course)))
	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}

func UpdateCourse(ctx context.Context) {
	//log.Printf("[UpdateCourse] update crs(%v) ", goutil.Struct2Json(ctx.Params().Get("id")))
	id := ctx.Params().Get("id")
	var crs db.Course
	err := ctx.ReadJSON(&crs)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	crs.ID = id
	err = db.UpdateCourse(&crs)
	if err != nil {
		log.Printf("[UpdateCourse] update crs(%v) error(%v)", goutil.Struct2Json(crs), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"course": crs,
	})
}


func CallUpdateCourse(token,id string, course *db.Course) (goutil.Map, error) {
	result, err := tools.HttpPut(fmt.Sprintf("http://%v/courses/%v", SrvAddr, id), token,
		"application/json", bytes.NewBufferString(goutil.Struct2Json(course)))
	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}

func DeleteCourse(ctx context.Context) {
	var ids []string
	err := ctx.ReadJSON(&ids)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	err = db.UpdateCourseByIDs(ids, goutil.Map{"status": db.CourseStatusDelete})
	if err != nil {
		log.Printf("[DeleteCourse] delete ids(%v) error(%v)", ids, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"ids": ids,
	})
}

func CallDeleteCourse(token string, id string) (goutil.Map, error) {

	result, err := tools.HttpGet(fmt.Sprintf("http://%v/courses/%v", SrvAddr, id), token)

	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}

func GetCourse(ctx context.Context) {
	id := ctx.Params().Get("id")
	crs, err := db.LoadCourse(id)
	if err != nil {
		log.Printf("[GetCourse] get crs(%v) error(%v)", id, err)
		if err == db.ErrNotFound {
			WriteResultErrByMsg(ctx, CodeCourseNotFound, "课程不存在", err)
			return
		}
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"course": crs,
	})
}

func CallGetCourse(token string, id string) (goutil.Map, error) {
	result, err := tools.HttpGet(fmt.Sprintf("http://%v/courses/%v", SrvAddr, id), token)
	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}

func GetCourses(ctx context.Context) {
	argMap, err := CheckURLArg(ctx.FormValues(), []*Arg{
		{Key: "name", Type: "string"},
		{Key: "id", Type: "string"},
		{Key: "deptId", Type: "string"},
		{Key: "majorId", Type: "string"},
		{Key: "period", Type: "string"},
		{Key: "credit", Type: "string"},
		{Key: "attr", Type: "string"},
		{Key: "nature", Type: "string"},
		{Key: "campus", Type: "string"},
		{Key: "status", Type: "int", DefaultValue: strconv.Itoa(db.CourseStatusChecked)},
		{Key: "page", Type: "int", DefaultValue: "0"},
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
	if argMap.Exist("attr") {
		fuzzyCondMap.Set("attr", argMap.Get("attr"))
	}
	if argMap.Exist("nature") {
		fuzzyCondMap.Set("name", argMap.Get("name"))
	}
	if argMap.Exist("campus") {
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
	if argMap.Exist("period") {
		exactCondMap.Set("period", argMap.Get("period"))
	}
	if argMap.Exist("credit") {
		exactCondMap.Set("credit", argMap.Get("credit"))
	}
	exactCondMap.Set("status", int(argMap.GetInt64("status")))
	limit := int(argMap.GetInt64("pageSize"))
	skip := int(argMap.GetInt64("page")) * limit
	var sort []string
	if argMap.Exist("sort") {
		sort = append(sort, argMap.GetString("sort"))
	}
	//log.Printf("[GetCourses] get exactCondMap(%v) fuzzyCondMap(%v)", exactCondMap, fuzzyCondMap)
	courseList, total, err := db.ListCourse(exactCondMap, fuzzyCondMap, sort, skip, limit)
	if err != nil {
		log.Printf("[GetCourses] get crs by(%v) error(%v)", argMap, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	log.Printf("[GetCourses] get courseList(%v) fuzzyCondMap(%v)", courseList, fuzzyCondMap)
	WriteResultSuccess(ctx, goutil.Map{
		"courseList": courseList,
		"total":       total,
	})
}

func CallGetCourses(token string, argMap goutil.Map) (goutil.Map, error) {
	url := fmt.Sprintf("http://%v/courses/list?", SrvAddr)
	result, err := tools.HttpGet(appendArgs(url, argMap), token)
	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}
