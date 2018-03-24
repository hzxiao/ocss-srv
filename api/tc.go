package api

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/kataras/iris/context"
	log "github.com/sirupsen/logrus"
	"strings"
)

func AddTeachCourse(ctx context.Context) {
	var tc db.TeachCourse
	err := ctx.ReadJSON(&tc)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	err = db.AddTeachCourse(&tc)
	if err != nil {
		log.Errorf("[AddTeachCourse] add tc(%v) error(%v)", goutil.Struct2Json(tc), err)
		if strings.Contains(err.Error(), "already exists") {
			WriteResultErrByMsg(ctx, CodeAlreadyExists, "课程号已存在", err)
			return
		}
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"tc": tc,
	})
}

func UpdateTeachCourse(ctx context.Context) {
	var tc db.TeachCourse
	err := ctx.ReadJSON(&tc)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	err = db.UpdateTeachCourseByIDs([]string{tc.ID}, &tc)
	if err != nil {
		log.Errorf("[UpdateTeachCourse] add tc(%v) error(%v)", goutil.Struct2Json(tc), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"tc": tc,
	})
}

func UpdateTeachCourses(ctx context.Context) {
	var info struct {
		Ids []string        `json:"ids"`
		Tc  *db.TeachCourse `json:"tc"`
	}
	err := ctx.ReadJSON(&info)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	err = db.UpdateTeachCourseByIDs(info.Ids, info.Tc)
	if err != nil {
		log.Errorf("[UpdateTeachCourse] add tc(%v) error(%v)", goutil.Struct2Json(info), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"ids": info.Ids,
	})
}

func ListTeachCourse(ctx context.Context) {
	argMap, err := CheckURLArg(ctx.FormValues(), []*Arg{
		{Key: "name", Type: "string"},
		{Key: "deptId", Type: "string"},
		{Key: "nature", Type: "string"},
		{Key: "attr", Type: "string"},
		{Key: "status", Type: "int"},
		{Key: "selectState", Type: "int"},
		{Key: "page", Type: "int", DefaultValue: "0"},
		{Key: "pageSize", Type: "int", DefaultValue: "20"},
		{Key: "sort", Type: "string"},
	})
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}
	var cids, tids []string
	if HasOneOfKeys(argMap, "name", "deptId", "nature") {
		crsExactMap := TakeByReplaceKeys(argMap, goutil.Map{"deptId": "dept.id", "nature": "nature", "attr": "attr"})
		crsExactMap.Set("status", db.CourseStatusChecking)
		crsFuzzyMap := TakeByKeys(argMap, "name")
		crs, _, err := db.ListCourse(crsExactMap, crsFuzzyMap, nil, 0, 0)
		if err != nil {
			log.Errorf("[ListTeachCourse] error(%v)", err)
			WriteResultWithSrvErr(ctx, err)
			return
		}
		for i := range crs {
			cids = append(cids, crs[i].ID)
		}
		if !argMap.Exist("name") && len(cids) == 0 {
			WriteResultSuccess(ctx, goutil.Map{
				"tcList": nil,
				"total":  0,
			})
			return
		}
		if argMap.Exist("name") {
			teExactMap := goutil.Map{"status": db.UserStatsNormal}
			teFuzzyMap := goutil.Map{"name": argMap.GetString("name")}
			tes, _, err := db.ListTeacher(teExactMap, teFuzzyMap, nil, 0, 0)
			if err != nil {
				log.Errorf("[ListTeachCourse] error(%v)", err)
				WriteResultWithSrvErr(ctx, err)
				return
			}
			for i := range tes {
				tids = append(tids, tes[i].ID)
			}
		}
		if len(cids) == 0 && len(tids) == 0 {
			WriteResultSuccess(ctx, goutil.Map{
				"tcList": nil,
				"total":  0,
			})
			return
		}
	}

	limit := int(argMap.GetInt64("pageSize"))
	skip := int(argMap.GetInt64("page")) * limit
	var sort []string
	if argMap.Exist("sort") {
		sort = append(sort, argMap.GetString("sort"))
	}
	status := int(argMap.GetInt64("status"))
	selectState := int(argMap.GetInt64("selectState"))
	tcs, total, err := db.ListTeachCourses(status, selectState, cids, tids, sort, skip, limit)
	if err != nil {
		log.Errorf("[ListTeachCourse] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	cids, tids = nil, nil
	for i := range tcs {
		cids = append(cids, tcs[i].CID)
		tids = append(tids, tcs[i].TID)
	}
	crsList, err := db.ListCourseByIds(cids)
	if err != nil {
		log.Errorf("[ListTeachCourse] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	teList, err := db.ListTeacherByIds(tids)
	if err != nil {
		log.Errorf("[ListTeachCourse] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	var tcList []goutil.Map
	for i := range tcs {
		tc := goutil.Struct2Map(tcs[i])
		for j := range crsList {
			if crsList[j].ID == tcs[i].CID {
				tc.Set("courseName", crsList[j].Name)
				if crsList[j].Dept != nil {
					tc.Set("deptName", crsList[j].Dept.GetString("name"))
					tc.Set("nature", crsList[j].Nature)
					tc.Set("attr", crsList[j].Attr)
				}
				break
			}
		}
		for k := range teList {
			if teList[k].ID == tcs[i].TID {
				tc.Set("teacherName", teList[k].Name)
			}
		}
		tcList = append(tcList, tc)
	}

	WriteResultSuccess(ctx, goutil.Map{
		"tcList": tcList,
		"total":  total,
	})
}

func DeleteTeachCourse(ctx context.Context) {
	var data goutil.Map
	err := ctx.ReadJSON(&data)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	ids := data.GetStringArray("ids")
	for _, id := range ids {
		err = db.DelTeachCourse(id)
		if err != nil {
			log.Errorf("[DeleteTeachCourse] delete ids(%v) error(%v)", ids, err)
			WriteResultWithSrvErr(ctx, err)
			return
		}
	}

	WriteResultSuccess(ctx, goutil.Map{
		"ids": ids,
	})
}
