package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	_, err = C(CollectionTeachCourse).UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": args})
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

func ListTeachCourse(status string, sort []string, skip, limit int) ([]goutil.Map, int, error) {
	return nil, 0, nil
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

//管理员：增加选课、修改选课、删除选课
//老师：打分 查看选择的学生信息

func SettingGrade(tcid string, method string, info []goutil.Map) error {
	return nil
}

func ListStuOfCourse() {

}

//学生：选课，查看、取消选课

func StuSelectCourse(tcid, sid string) error {
	return nil
}

func StuCancelCourse(tcid, sid string) error {
	return nil
}

func ListStuLearnCourse(status, sid string, sort []string, skip, limit int) ([]goutil.Map, int, error) {
	return nil, 0, nil
}
