package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

func AddTeachCourse(tc *TeachCourse) error {
	if tc == nil {
		return errors.New("tc is nil")
	}
	if tc.CID == "" {
		return errors.New("cid is empty")
	}
	if tc.TID == "" {
		return errors.New("tid is empty")
	}
	if tc.Capacity == 0 {
		tc.Capacity = 50
	}
	tc.Margin = tc.Capacity
	tc.Status = TeachCourseStatusSelectable
	tc.Create = tools.NowMillisecond()
	tc.Update = tc.Create
	tc.ID = tools.GenerateUniqueId()

	info, err := C(CollectionTeachCourse).Find(bson.M{"cid": tc.CID, "tid": tc.TID}).Apply(mgo.Change{
		Update:    bson.M{"$setOnInsert": tc},
		Upsert:    true,
		ReturnNew: true,
	}, tc)
	if err != nil {
		return err
	}
	if info.UpsertedId == nil {
		return ErrAlreadyExist
	}
	return nil
}

func UpdateTeachCourseByIDs(ids []string, tc *TeachCourse) (err error) {
	if len(ids) == 0 {
		return nil
	}
	if tc == nil {
		return errors.New("tc is nil")
	}
	args := bson.M{}
	finder:=bson.M{"_id": bson.M{"$in": ids}}
	if tc.TID != "" {
		args["tid"] = tc.TID
	}
	if tc.CID != "" {
		args["cid"] = tc.CID
	}
	if tc.Capacity > 0 {
		args["capacity"] = tc.Capacity
		defer func() {
			err = setCourseMargin(ids)
		}()
	}
	if tc.Status > 0 {
		args["status"] = tc.Status
	}
	if len(tc.Addr) > 0 {
		args["addr"] = tc.Addr
	}
	if len(tc.TakeTimes) > 0 {
		args["takeTimes"] = tc.TakeTimes
	}
	if len(tc.TakeWeeks) > 0 {
		args["takeWeeks"] = tc.TakeWeeks
	}
	if tc.StartSelectTime > 0 {
		args["startSelectTime"] = tc.StartSelectTime
	}
	if tc.EndSelectTime > 0 {
		args["endSelectTime"] = tc.EndSelectTime
	}

	args["update"] = tools.NowMillisecond()

	_, err = C(CollectionTeachCourse).UpdateAll(finder, bson.M{"$set": args})
	if err != nil {
		return err
	}
	return nil
}

func setCourseMargin(ids []string) error {
	var tcs []*TeachCourse
	err := C(CollectionTeachCourse).Find(bson.M{"_id": bson.M{"$in": ids}}).All(&tcs)
	if err != nil {
		return err
	}

	for _, tc := range tcs {
		err = C(CollectionTeachCourse).UpdateId(tc.ID, bson.M{"$set": bson.M{"margin": tc.Capacity - len(tc.StuInfo)}})
		if err != nil {
			return err
		}
	}

	return nil
}

func DelTeachCourse(tcid string) error {
	var tc TeachCourse
	err := one(CollectionTeachCourse, bson.M{"_id": tcid}, nil, &tc)
	if err != nil {
		return err
	}
	if tc.Margin < tc.Capacity {
		return errors.New("the course has been selected by student")
	}

	return C(CollectionTeachCourse).RemoveId(tcid)
}

func ListTeachCourse(status int, sort []string, skip, limit int) ([]*TeachCourse, int, error) {
	finder := bson.M{}
	if status > 0 {
		finder["status"] = status

	}
	var teachCourseList []*TeachCourse
	total, err := list(CollectionTeachCourse, finder, nil, sort, skip, limit, &teachCourseList)
	if err != nil {
		return nil, 0, err
	}

	return teachCourseList, total, nil
}

func IsCourseSelectable(tcid string) (bool, error) {
	now := tools.NowMillisecond()
	find := bson.M{
		"_id":             tcid,
		"margin":          bson.M{"$gt": 0},
		"startSelectTime": bson.M{"$lte": now},
		"endSelectTime":   bson.M{"$lte": now},
		"status":          TeachCourseStatusSelectable,
	}
	var tc TeachCourse
	err := one(CollectionTeachCourse, find, bson.M{"_id": 1}, &tc)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CountTeachCourse(cond goutil.Map) (int, error) {
	return count(CollectionTeachCourse, tools.ToBsonMap(cond))
}

//管理员：增加选课、修改选课、删除选课
//通过选课

//老师：打分 查看选择的学生信息

func TeaSettingGrade(tcid string, info []goutil.Map) error {
	if tcid == ""{
		return errors.New("tcid is nil")
	}
	if info == nil || len(info) == 0{
		return errors.New("grade info is nil")
	}
	var sids []string
	var stuInfo goutil.Map = goutil.Map{}
	finder:=bson.M{"_id":tcid}
	update_time := tools.NowMillisecond()
	for _,v:= range info {
		sids = append(sids,v.GetStringP("sid"))
		sin:=goutil.Map{}
		if v.GetFloat64P("grade") > 0 {
			finder["stuInfo.cstatus"] = SCourseStatusSelected
			finder["status"] = TeachCourseStatusLearning
			sin.Set("stuInfo.$.grade",v.GetFloat64P("grade"))
		}
		if v.GetFloat64P("ordinaryGrade") > 0 {
			finder["stuInfo.cstatus"] = SCourseStatusSelected
			finder["status"] = TeachCourseStatusLearning
			sin.Set("stuInfo.$.ordinaryGrade",v.GetFloat64P("ordinaryGrade"))
		}
		if v.GetFloat64P("examGrade") > 0 {
			finder["stuInfo.cstatus"] = SCourseStatusSelected
			finder["status"] = TeachCourseStatusLearning
			sin.Set("stuInfo.$.examGrade",v.GetFloat64P("examGrade"))
		}
		if v.GetFloat64P("cstatus") > 0 {
			sin.Set("stuInfo.$.cstatus",v.GetInt64P("cstatus"))
		}
		sin.Set("stuInfo.$.update",update_time)

		stuInfo.Set(v.GetStringP("sid"),sin)
	}
	//finder["stuInfo.sid"] = bson.M{"$in":sids}
	//fmt.Println(tools.Struct2BsonMap(stuInfo))
	//先查
	var teachCourseList []*TeachCourse
	total, err := list(CollectionTeachCourse, finder, bson.M{"stuInfo":1}, nil, 0, 0, &teachCourseList)
	if err != nil {
		return err
	}
	if total == 0 || teachCourseList == nil || teachCourseList[0].StuInfo == nil || len(teachCourseList[0].StuInfo) < len(sids){
		return errors.New(fmt.Sprintf("One or more student is not in this learn course(%v)", goutil.Struct2Json(teachCourseList)))
	}
	scout:=len(sids)
	flag := ""
	for _,v:= range sids {
		for j,s:= range teachCourseList[0].StuInfo  {
			if j == len(teachCourseList[0].StuInfo) && v!=s.GetStringP("sid"){
				flag = 	v
				break
			}
			if v == s.GetStringP("sid") {
				scout --
				break
			}
		}
		if flag != "" {
			break
		}
	}
	if flag != ""{
		return errors.New(fmt.Sprintf("One or more student(%v) is not in this learn course(%v)", flag,goutil.Struct2Json(flag)))
	}


	changeI := bson.M{}
	for _,v:= range teachCourseList[0].StuInfo {
		changeI = tools.ToBsonMap(stuInfo.GetMapP(v.GetStringP("sid")))
		changeI["update"] = update_time
		finder["stuInfo.sid"] = v.GetStringP("sid")
		fmt.Println("finder,",goutil.Struct2Json(finder),";change",goutil.Struct2Json(changeI))
		_, err = C(CollectionTeachCourse).UpdateAll(finder,
			bson.M{ "$set":changeI})
	}

	if err != nil {
		return err
	}

	return nil
}

func TeaListStuOfCourse(tid string, sids []string, status int, sort []string, skip, limit int) ([]*TeachCourse, int, error) {
	if tid == "" {
		return nil, 0, errors.New("tid is nil")
	}
	finder := bson.M{"tid":tid}
	if status != 0 {
		finder["stuInfo.cstatus"] = status
	}
	if sids != nil && len(sids) > 0 {
		finder["stuInfo.sid"] = bson.M{"$in":sids}
	}
	var teachCourseList []*TeachCourse
	selector := bson.M{"stuInfo.$":1, "_id":1, "cid":1, "tid":1, "takeWeeks":1, "takeTimes":1, "startSelectTime":1, "endSelectTime":1, "addr":1, "capacity":1, "margin":1, "status":1, "create":1, "update":1}
	total, err := list(CollectionTeachCourse, finder, selector, sort, skip, limit, &teachCourseList)
	if err != nil {
		return nil, 0, err
	}
	return teachCourseList, total, nil
}

//学生：选课，查看、取消选课
func StuSelectCourse(tcids []string, sid string) error {
	if tcids == nil || len(tcids) == 0 {
		return errors.New("tcid is nil")
	}
	if sid == "" {
		return errors.New("sid is nil")
	}
	exactCond := []bson.M{}
	exactCond = append(exactCond, bson.M{"_id":bson.M{"$in":tcids}})
	exactCond = append(exactCond, bson.M{"status":TeachCourseStatusSelectable})
	//判断该些课程是否满员可选
	var teachCourseCon []*TeachCourse
	finder := bson.M{"_id":bson.M{"$in":tcids}, "status":TeachCourseStatusSelectable, "margin":bson.M{"$gte":1}}
	total, err := list(CollectionTeachCourse, finder, bson.M{"_id":1}, nil, 0, 0, &teachCourseCon)
	if err != nil {
		return err
	}
	if total != len(tcids)  {
		return errors.New(fmt.Sprintf("This Student has select the course(%v)", goutil.Struct2Json(teachCourseCon)))
	}
	//判断该学生是否已经选了该谢门课程
	var teachCourseList []*TeachCourse
	finder = bson.M{"_id":bson.M{"$in":tcids}, "status":TeachCourseStatusSelectable, "stuInfo.sid":sid}
	total, err = list(CollectionTeachCourse, finder, bson.M{"_id":1}, nil, 0, 0, &teachCourseList)
	if err != nil {
		return err
	}
	if total > 0 || teachCourseList != nil {
		return errors.New(fmt.Sprintf("This Student has select the course(%v)", goutil.Struct2Json(teachCourseList)))
	}

	// 若没有存在已选课程，则添加课程
	create_time := tools.NowMillisecond()
	update_time := create_time
	_, err = C(CollectionTeachCourse).UpdateAll(bson.M{"_id":bson.M{"$in":tcids}, "status":TeachCourseStatusSelectable},
		bson.M{"$inc":bson.M{"margin":-1}, "$set":bson.M{"update":update_time}, "$addToSet":bson.M{"stuInfo":bson.M{"sid":sid, "cstatus":SCourseStatusSelecting, "create":create_time, "update":update_time}}})
	//
	if err != nil {
		return err
	}

	return nil
}

func StuCancelCourse(tcids []string, sid string) error {
	if tcids == nil || len(tcids) == 0 {
		return errors.New("tcid is nil")
	}
	if sid == "" {
		return errors.New("sid is nil")
	}
	update_time := tools.NowMillisecond()
	// 取消正在选的课，即删除课程
	_, err := C(CollectionTeachCourse).UpdateAll(bson.M{"_id":bson.M{"$in":tcids}, "status":TeachCourseStatusSelectable,
		"stuInfo.sid":sid, "stuInfo.cstatus":SCourseStatusSelecting},
		bson.M{"$inc":bson.M{"margin":1}, "$set":bson.M{"update":update_time}, "$pull":bson.M{"stuInfo":bson.M{"sid":sid}}})
	if err != nil {
		return err
	}
	return nil
}

func ListStuLearnCourse(cstatus int, sid string, sort []string, skip, limit int) ([]*TeachCourse, int, error) {
	if sid == "" {
		return nil, 0, errors.New("sid is nil")
	}
	finder := bson.M{"stuInfo.sid":sid}

	if cstatus != 0 {
		finder["stuInfo.cstatus"] = cstatus
	}
	var teachCourseList []*TeachCourse
	total, err := list(CollectionTeachCourse, finder, bson.M{"_id":1}, sort, skip, limit, &teachCourseList)
	if err != nil {
		return nil, 0, err
	}
	return teachCourseList, total, nil
}

