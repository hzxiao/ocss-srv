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
	var data goutil.Map
	err := ctx.ReadJSON(&data)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	ids := data.GetStringArray("ids")
	err = db.UpdateStudentByIDs(ids, goutil.Map{"status": db.UserStatsDelete})
	if err != nil {
		log.Errorf("[DeleteStudent] delete ids(%v) error(%v)", ids, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	err = db.UpdateUserByIDs(ids, goutil.Map{"status": db.UserStatsDelete})
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
		{Key: "schoolYear", Type: "string"},
		{Key: "class", Type: "string"},
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
		fuzzyCondMap.Set("_id", argMap.Get("id"))
	}
	if argMap.Exist("deptId") {
		exactCondMap.Set("dept.id", argMap.Get("deptId"))
	}
	if argMap.Exist("majorId") {
		exactCondMap.Set("major.id", argMap.Get("majorId"))
	}
	if argMap.Exist("schoolYear") {
		exactCondMap.Set("schoolYear", argMap.Get("schoolYear"))
	}
	if argMap.Exist("class") {
		exactCondMap.Set("class", argMap.Get("class"))
	}
	exactCondMap.Set("status", int(argMap.GetInt64("status")))
	page := argMap.GetInt64("page")
	if page > 0 {
		page--
	}
	limit := int(argMap.GetInt64("pageSize"))
	skip := int(page) * limit
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

func CheckImportStudent(ctx context.Context) {
	var stus []*db.Student
	err := ctx.ReadJSON(&stus)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	//get all dept info
	depts, err := db.FindAllDept()
	if err != nil {
		log.Errorf("[CheckImportStudent] find all dept error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	//get all major
	majors, err := db.FindAllMajor("")
	if err != nil {
		log.Errorf("[CheckImportStudent] find all major error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	var stuOfDeptCountMap = make(map[string]int)
	var errInfo = make([]string, len(stus))
	for i, stu := range stus {
		var dept string
		if stu.Dept != nil {
			dept = stu.Dept.GetString("name")
		}
		if dept == "" {
			errInfo[i] = "学院为空"
			continue
		}
		var validDept bool
		for _, d := range depts {
			if d.Name == dept {
				stu.Dept.Set("id", d.ID)
				validDept = true
				break
			}
		}
		if !validDept {
			errInfo[i] = "无效的学院"
			continue
		}
		//major
		var major string
		if stu.Major != nil {
			major = stu.Major.GetString("name")
		}
		for _, m := range majors {
			if m.Name == major {
				stu.Major.Set("id", m.ID)
			}
		}
		if !tools.ContainElem([]string{"2014", "2015", "2016", "2017"}, stu.SchoolYear) {
			errInfo[i] = "无效的年级"
			continue
		}
		if stu.ID != "" {
			cnt, err := db.CountStudent(goutil.Map{"_id": stu.ID})
			if err != nil {
				log.Errorf("[CheckImportStudent] count student by id(%v) error(%v)", stu.Dept, err)
				WriteResultWithSrvErr(ctx, err)
				return
			}
			if cnt > 0 {
				errInfo[i] = "学号已存在"
				continue
			}
		} else {
			key := stu.Dept.GetString("id") + stu.SchoolYear
			num, ok := stuOfDeptCountMap[key]
			if !ok {
				cnt, err := db.CountStudent(goutil.Map{
					"dept.id":    stu.Dept.GetString("id"),
					"schoolYear": stu.SchoolYear,
				})
				if err != nil {
					log.Errorf("[CheckImportStudent] count student by error(%v)", err)
					WriteResultWithSrvErr(ctx, err)
					return
				}
				stuOfDeptCountMap[key] = cnt
				num = cnt
			}

			stu.ID = stu.SchoolYear[2:] + stu.Dept.GetString("id") + fmt.Sprintf("%03d", num+1)
			stuOfDeptCountMap[key]++
		}
	}

	WriteResultSuccess(ctx, goutil.Map{
		"studentList": stus,
		"errInfo":     errInfo,
	})
}

func ImportStudent(ctx context.Context) {
	var stus []*db.Student
	err := ctx.ReadJSON(&stus)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	for _, stu := range stus {
		err = db.AddStudent(stu)
		if err != nil {
			log.Errorf("[ImportStudent] add stu(%v) error(%v)", goutil.Struct2Json(stu), err)
			if strings.Contains(err.Error(), "already exists") {
				WriteResultErrByMsg(ctx, CodeAlreadyExists, "学号已存在", err)
				return
			}
			WriteResultWithSrvErr(ctx, err)
			return
		}
	}

	WriteResultSuccess(ctx, goutil.Map{
		"studentList": stus,
	})
}
