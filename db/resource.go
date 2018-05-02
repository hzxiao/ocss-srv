package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func AddFile(file *File) error {
	if file == nil {
		return errors.New("file is nil")
	}
	if file.ID == "" {
		return errors.New("id is empty")
	}

	return C(CollectionFile).Insert(file)
}

func LoadFile(id string) (*File, error) {
	var file File
	err := one(CollectionFile, bson.M{"_id": id}, nil, &file)
	return &file, err
}

func AddResource(r *CourseResource) error {
	if r == nil {
		return errors.New("r is nil")
	}

	if r.TCID == "" {
		return errors.New("tcid is empty")
	}
	if r.TID == "" {
		return errors.New("tid is empty")
	}
	if r.File == nil {
		return errors.New("file is nil")
	}

	r.ID = tools.GenerateUniqueId()
	r.Create = tools.NowMillisecond()
	r.Update = r.Create
	r.Status = StatusNormal
	return C(CollectionResource).Insert(r)
}

func NotifyResourceToStudent(r CourseResource) error {
	log.Info("NotifyResourceToStudent: id=", r)
	teacher, err := LoadTeacher(r.TID)
	if err != nil {
		log.Error(err)
		return err
	}
	tc, err := LoadTeachCourse(r.TCID)
	if err != nil {
		log.Error(err)
		return err
	}
	crs, err := LoadCourse(tc.CID)
	if err != nil {
		log.Error(err)
		return err
	}
	title := fmt.Sprintf("%v在%v这门课程上传了课件！！！", teacher.Name, crs.Name)
	content := fmt.Sprintf("任课教师：%v，课程名称：%v，课件名称：%v",
		teacher.Name,  crs.Name, r.File.Name)
	stu := []string{}
	if tc.StuInfo != nil &&  len(tc.StuInfo) > 0 {
		for j := range tc.StuInfo {
			stu = append(stu, tc.StuInfo[j].GetString("sid"))
		}
	}
	log.Info("NotifyTcOverSelect2Student: stu(%v)", stu)
	if len(stu) > 0 {
		SendNotice(RoleStudent, goutil.Map{
			"title":   title,
			"content": content,
			"sid": stu,
		})
	}

	return nil
}

func DelCourseResource(tid, tcid string, rids []string) error {
	cond := bson.M{"_id": bson.M{"$in": rids}}
	if tid != "" {
		cond["tid"] = tid
	}
	if tcid != "" {
		cond["tcid"] = tcid
	}

	_, err := C(CollectionResource).UpdateAll(cond, bson.M{"$set": bson.M{"status": StatusDeleted}})
	return err
}

func ListCourseResource(cond goutil.Map, sort []string, skip, limit int) ([]*CourseResource, int, error) {
	var resourceList []*CourseResource
	total, err := list(CollectionResource, tools.ToBsonMap(cond), nil, sort, skip, limit, &resourceList)
	if err != nil {
		return nil, 0, err
	}
	return resourceList, total, nil
}
