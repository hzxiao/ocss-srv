package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2/bson"
)

func AddStudent(stu *Student) error {
	if stu == nil {
		return errors.New("stu is nil")
	}
	if stu.ID == "" {
		return errors.New("id is empty")
	}
	stu.Status = UserStatsNormal
	stu.Create = tools.NowMillisecond()
	stu.Update = stu.Create

	user := &User{
		Username: stu.ID,
		Status:   UserStatsNormal,
		Role:     RoleStudent,
	}
	err := AddUser(user)
	if err != nil {
		return err
	}

	return C(CollectionStudent).Insert(stu)
}

func UpdateStudent(stu *Student) error {
	if stu == nil {
		return errors.New("stu is nil")
	}
	if stu.ID == "" {
		return errors.New("id is empty")
	}
	args, err := tools.Struct2BsonMap(stu)
	if err != nil {
		return err
	}
	delete(args, "_id")
	delete(args, "create")
	delete(args, "status")

	args["update"] = tools.NowMillisecond()
	return C(CollectionStudent).UpdateId(stu.ID, bson.M{"$set": args})
}

func UpdateStudentByIDs(ids []string, update goutil.Map) error {
	if len(ids) == 0 {
		return nil
	}
	if len(update) == 0 {
		return errors.New("update is empty")
	}

	_, err := C(CollectionStudent).UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": tools.ToBsonMap(update)})
	return err
}

func LoadStudent(id string) (*Student, error) {
	var stu Student
	err := one(CollectionStudent, bson.M{"_id": id}, nil, &stu)
	return &stu, err
}

func ListStudent(exactCond, fuzzyCond goutil.Map, sort []string, skip, limit int) ([]*Student, int, error) {
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
	var stuList []*Student
	total, err := list(CollectionStudent, finder, nil, sort, skip, limit, &stuList)
	if err != nil {
		return nil, 0, err
	}
	return stuList, total, nil
}

func CountStudent(cond goutil.Map) (int, error) {
	return count(CollectionStudent, tools.ToBsonMap(cond))
}
