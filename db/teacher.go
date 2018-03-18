package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2/bson"
)

func AddTeacher(t *Teacher) error {
	if t == nil {
		return errors.New("t is nil")
	}
	if t.ID == "" {
		return errors.New("id is empty")
	}
	t.Status = UserStatsNormal
	t.Create = tools.NowMillisecond()
	t.Update = t.Create

	user := &User{
		Username: t.ID,
		Status:   UserStatsNormal,
		Role:     RoleTeacher,
	}
	err := AddUser(user)
	if err != nil {
		return err
	}

	return C(CollectionTeacher).Insert(t)
}

func UpdateTeacher(t *Teacher) error {
	if t == nil {
		return errors.New("t is nil")
	}
	if t.ID == "" {
		return errors.New("id is empty")
	}
	args, err := tools.Struct2BsonMap(t)
	if err != nil {
		return err
	}
	delete(args, "_id")
	delete(args, "create")
	delete(args, "status")

	args["update"] = tools.NowMillisecond()
	return C(CollectionTeacher).UpdateId(t.ID, bson.M{"$set": args})
}

func UpdateTeacherByIDs(ids []string, update goutil.Map) error {
	if len(ids) == 0 {
		return nil
	}
	if len(update) == 0 {
		return errors.New("update is empty")
	}

	_, err := C(CollectionTeacher).UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": tools.ToBsonMap(update)})
	return err
}

func ListTeacher(exactCond, fuzzyCond goutil.Map, sort []string, skip, limit int) ([]*Teacher, int, error) {
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
	var teacherList []*Teacher
	total, err := list(CollectionTeacher, finder, nil, sort, skip, limit, &teacherList)
	if err != nil {
		return nil, 0, err
	}
	return teacherList, total, nil
}

func CountTeacher(cond goutil.Map) (int, error) {
	return count(CollectionTeacher, tools.ToBsonMap(cond))
}

func LoadTeacher(id string) (*Teacher, error) {
	var t Teacher
	err := one(CollectionTeacher, bson.M{"_id": id}, nil, &t)
	return &t, err
}

func ListTeacherByIds(ids []string) ([]*Teacher, error) {
	finder := bson.M{"_id": bson.M{"$in": ids}}
	var teList []*Teacher
	_, err := list(CollectionTeacher, finder, nil, nil, 0, 0, &teList)
	if err != nil {
		return nil, err
	}
	return teList, nil
}
