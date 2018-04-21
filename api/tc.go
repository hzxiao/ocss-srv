package api

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/juju/errors"
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
		} else if strings.Contains(err.Error(), "time conflict") {
			WriteResultErrByMsg(ctx, CodeAlreadyExists, "存在课程与教师与该选课时间冲突的选课", err)
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
		{Key: "tid", Type: "string"},
	})
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}
	var cids, tids,sids []string
	if HasOneOfKeys(argMap, "name", "deptId", "nature", "attr") {
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
	if argMap.Exist("tid") {
		tids = nil
		tids = append(tids, argMap.GetString("tid"))
	}
	status := int(argMap.GetInt64("status"))
	selectState := int(argMap.GetInt64("selectState"))
	tcs, total, err := db.ListTeachCourses(status, selectState,nil, cids, tids, sort, skip, limit)
	if err != nil {
		log.Errorf("[ListTeachCourse] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	cids, tids = nil, nil
	for i := range tcs {
		cids = append(cids, tcs[i].CID)
		tids = append(tids, tcs[i].TID)

		if tcs[i].StuInfo != nil&&len(tcs[i].StuInfo)>0 {
			for j := range tcs[i].StuInfo{
				sids = append(sids,tcs[i].StuInfo[j].GetString("sid"))
			}	
		}
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
	stuList, err := db.ListStudentByIds(sids)
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
					tc.Set("credit", crsList[j].Credit)
				}
				break
			}
		}
		if tcs[i].StuInfo != nil&&len(tcs[i].StuInfo)>0 && stuList != nil && len(stuList) > 0{
			stuInfo := tcs[i].StuInfo
			// 成绩
			var grade,ordinaryGrade,examGrade float64
			var  gs,ogs,egs int
			for m := range stuInfo{
				for n := range stuList{
					if stuInfo[m].GetString("sid") == stuList[n].ID {
						stuInfo[m].Set("name",stuList[n].Name)
						stuInfo[m].Set("class",stuList[n].Class)
						stuInfo[m].Set("schoolYear",stuList[n].SchoolYear)
						stuInfo[m].Set("deptName",stuList[n].Dept.GetString("name"))
						sex := "男"
						if stuList[n].Dept.GetString("sex") != "male"{
							sex = "女"
						}
						stuInfo[m].Set("majorName",sex)
					}
				}
				if stuInfo[m].Get("grade") != nil{
					gs ++
					grade += stuInfo[m].GetFloat64("grade")
				}
				if stuInfo[m].Get("ordinaryGrade") != nil{
					ogs ++
					ordinaryGrade += stuInfo[m].GetFloat64("ordinaryGrade")
				}
				if stuInfo[m].Get("examGrade") != nil{
					egs ++
					examGrade += stuInfo[m].GetFloat64("examGrade")
				}
			}
			tc.Set("stuInfo", stuInfo)
			log.Printf("[ListTeachCourse] get  stuInfo(%v) ,students(%v)", stuInfo, stuList)

			if gs > 0{
				grade = grade/(float64(gs))
				tc.Set("grade", grade)
			}
			tc.Set("gradeNum", gs)
			if ogs > 0{
				ordinaryGrade = ordinaryGrade/float64(ogs)
				tc.Set("ordinaryGrade", ordinaryGrade)
			}
			if egs > 0{
				examGrade = examGrade/float64(egs)
				tc.Set("examGrade", examGrade)
			}
		}


		for k := range teList {
			if teList[k].ID == tcs[i].TID {
				tc.Set("teacherName", teList[k].Name)
			}
		}
		tcList = append(tcList, tc)
	}
	log.Printf("[ListTeachCourse] get courseList(%v) ", tcList)

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

func ListStudentCourse(ctx context.Context) {
	argMap, err := CheckURLArg(ctx.FormValues(), []*Arg{
		{Key: "name", Type: "string"},
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
	if HasOneOfKeys(argMap, "name") {
		crsExactMap := goutil.Map{}
		crsExactMap.Set("status", db.CourseStatusChecking)
		crsFuzzyMap := TakeByKeys(argMap, "name")
		//crsFuzzyMap := goutil.Map{}
		//if argMap.Exist("name") {
		//	crsFuzzyMap.Set("name", argMap.Get("name"))
		//}
		crs, _, err := db.ListCourse(crsExactMap, crsFuzzyMap, nil, 0, 0)
		if err != nil {
			log.Errorf("[ListStudentCourse] error(%v)", err)
			WriteResultWithSrvErr(ctx, err)
			return
		}
		for i := range crs {
			cids = append(cids, crs[i].ID)
		}
		//log.Printf("[ListStudentCourse] get courseList(%v) cids[%v] fuzzyCondMap(%v)", crs,cids, crsFuzzyMap)

		if  len(cids) == 0 {
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
	selectState := int(argMap.GetInt64("selectState"))
	sid := ctx.Values().GetString("uid")
	tcs, total, err := db.ListStudentCourse(selectState, sid, cids, sort, skip, limit)
	if err != nil {
		log.Errorf("[ListStudentCourse] error(%v)", err)
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
		log.Errorf("[ListStudentCourse] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	teList, err := db.ListTeacherByIds(tids)
	if err != nil {
		log.Errorf("[ListStudentCourse] error(%v)", err)
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
	log.Printf("[ListStudentCourse] get tcList(%v) ", tcList)

	WriteResultSuccess(ctx, goutil.Map{
		"tcList": tcList,
		"total":  total,
	})
}

func StuSelectCourse(ctx context.Context) {
	var data goutil.Map
	err := ctx.ReadJSON(&data)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	sid := ctx.Values().GetString("uid")
	switch data.GetString("method") {
	case "select":
		err = db.StuSelectCourse(data.GetStringArray("ids"), sid)
	case "cancel":
		err = db.StuCancelCourse(data.GetStringArray("ids"), sid)
	default:
		WriteResultWithArgErr(ctx, errors.New("unknown method"))
		return
	}
	if err != nil {
		log.Errorf("[StuSelectCourse] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	go db.NotifyTcFull2Adm(data.GetStringArray("ids"))
	WriteResultSuccess(ctx, "OK")
}

func GetTeachCourse(ctx context.Context) {
	id := ctx.Params().Get("id")
	tc, err := db.LoadTeachCourse(id)
	if err != nil {
		log.Errorf("[GetTeachCourse] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	course, err := db.LoadCourse(tc.CID)
	if err != nil {
		log.Errorf("[GetTeachCourse] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	teacher, err := db.LoadTeacher(tc.TID)
	if err != nil {
		log.Errorf("[GetTeachCourse] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	WriteResultSuccess(ctx, goutil.Map{
		"tc":      tc,
		"course":  course,
		"teacher": teacher,
	})
}

func ListStudentOfCourse(ctx context.Context) {
	id := ctx.Params().Get("id")
	tc, err := db.LoadTeachCourse(id)
	if err != nil {
		log.Errorf("[ListStudentOfCourse] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	var sids []string
	for _, item := range tc.StuInfo {
		sids = append(sids, item.GetString("sid"))
	}

	students, err := db.ListStudentByIds(sids)
	if err != nil {
		log.Errorf("[ListStudentOfCourse] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	var studentList []goutil.Map
	for _, s := range students {
		stu := goutil.Struct2Map(s)
		for _, item := range tc.StuInfo {
			stu.Set("selectTime", item.Get("create"))
			if item.Get("grade") != nil {
				stu.Set("grade", item.GetFloat64("grade"))
			}
			if item.Get("ordinaryGrade") != nil {
				stu.Set("ordinaryGrade", item.GetFloat64("ordinaryGrade"))
			}
			if item.Get("examGrade") != nil {
				stu.Set("examGrade", item.GetFloat64("examGrade"))
			}
		}
		studentList = append(studentList, stu)
	}

	WriteResultSuccess(ctx, goutil.Map{
		"tc":      tc,
		"studentList":  studentList,
	})
}



func UpdateStudentForTc(ctx context.Context)  {
	var data goutil.Map
	err := ctx.ReadJSON(&data)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	switch data.GetString("method") {
	case "select":
		err = db.StuSelectCourse([]string{data.GetString("id")}, data.GetString("sid"))
	case "cancel":
		err = db.StuCancelCourse([]string{data.GetString("id")}, data.GetString("sid"))
	default:
		WriteResultWithArgErr(ctx, errors.New("unknown method"))
		return
	}
	if err != nil {
		log.Errorf("[UpdateStudentForTc] error(%v)", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	go db.NotifyTcFull2Adm([]string{data.GetString("id")})

	WriteResultSuccess(ctx, "OK")
}

// 老师导入课程的学生成绩
func SetGradeForTc(ctx context.Context) {
	var info struct {
		Tc  *db.TeachCourse `json:"tc"`
	}
	err := ctx.ReadJSON(&info)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	err = db.TeaSettingGrade(info.Tc.ID, info.Tc.StuInfo)
	if err != nil {
		log.Errorf("[SetGradeForTc] add tc(%v) error(%v)", goutil.Struct2Json(info), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, "OK")

}