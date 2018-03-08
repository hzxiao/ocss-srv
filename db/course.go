package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2/bson"
)

func AddCourse(course *Course) (*Course,error) {
	if course == nil {
		return nil,errors.New("course is nil")
	}
	if course.Name == "" {
		return nil,errors.New("name is empty")
	}
	course.ID = tools.GenerateUniqueId()
	course.Status = CourseStatusChecking
	course.Create = tools.NowMillisecond()
	course.Update = course.Create

	//err := AddCourse(course)
	//if err != nil {
	//	return err
	//}

	return course,C(CollectionCourse).Insert(course)
}

func UpdateCourse(course *Course) error {
	if course == nil {
		return errors.New("course is nil")
	}
	if course.ID == "" {
		return errors.New("id is empty")
	}
	args, err := tools.Struct2BsonMap(course)
	if err != nil {
		return err
	}
	delete(args, "_id")
	delete(args, "create")
	delete(args, "creator")

	args["update"] = tools.NowMillisecond()
	return C(CollectionCourse).UpdateId(course.ID, bson.M{"$set": args})
}

func UpdateCourseByIDs(ids []string, update goutil.Map) error {
	if len(ids) == 0 {
		return nil
	}
	if len(update) == 0 {
		return errors.New("update is empty")
	}

	_, err := C(CollectionCourse).UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": tools.ToBsonMap(update)})
	return err
}

func LoadCourse(id string) (*Course, error) {
	var course Course
	err := one(CollectionCourse, bson.M{"_id": id}, nil, &course)
	return &course, err
}

func ListCourse(exactCond, fuzzyCond goutil.Map, sort []string, skip, limit int) ([]*Course, int, error) {
	finder := tools.ToBsonMap(exactCond)
	var fuzzyConds []bson.M
	for k := range fuzzyCond {
		fuzzyConds = append(fuzzyConds, bson.M{k: bson.M{"$regex": tools.ParseRegex(fuzzyCond.GetString(k)), "$options": "i"}})
	}
	if len(fuzzyConds) > 0 {
		if finder == nil {
			finder = bson.M{}
		}
		if len(fuzzyConds) == 1 {
			for k, v := range fuzzyConds[0] {
				finder[k] = v
			}
		} else {
			finder["$or"] = fuzzyConds
		}
	}
	var courseList []*Course
	total, err := list(CollectionCourse, finder, nil, sort, skip, limit, &courseList)
	if err != nil {
		return nil, 0, err
	}
	return courseList, total, nil
}

func CountCourse(cond goutil.Map) (int, error) {
	return count(CollectionCourse, tools.ToBsonMap(cond))
}
