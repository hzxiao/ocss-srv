package db

import (
	"encoding/json"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
)

//InitDept int dept data
func InitDept(filename string) error {
	var deptList []*Dept
	err := unmarshalJsonFile(filename, &deptList)
	if err != nil {
		return err
	}
	for i := range deptList {
		err = AddDept(deptList[i])
		if err != nil {
			return err
		}
	}
	return nil
}

//AddDept upsert dept data by id
func AddDept(dept *Dept) error {
	if dept == nil {
		return errors.New("dept is nil")
	}
	if dept.ID == "" {
		return errors.New("id is empty")
	}
	_, err := C(CollectionDept).FindId(dept.ID).Apply(mgo.Change{
		Update: bson.M{
			"$set": dept,
		},
		Upsert: true,
	}, dept)
	return err
}

func FindAllDept() ([]*Dept, error) {
	var deptList []*Dept
	err := C(CollectionDept).Find(nil).All(&deptList)
	return deptList, err
}

func InitMajor(filename string) error {
	var majorList []*Major
	err := unmarshalJsonFile(filename, &majorList)
	if err != nil {
		return err
	}
	for i := range majorList {
		err = AddMajor(majorList[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func AddMajor(major *Major) error {
	if major == nil {
		return errors.New("major is nil")
	}
	if major.ID == "" {
		return errors.New("id is empty")
	}
	if major.DeptID == "" {
		return errors.New("deptId is empty")
	}

	_, err := C(CollectionMajor).FindId(major.ID).Apply(mgo.Change{
		Update: bson.M{
			"$set": major,
		},
		Upsert: true,
	}, major)
	return err
}

func FindAllMajor() ([]*Major, error) {
	var majorList []*Major
	err := C(CollectionMajor).Find(nil).All(&majorList)
	return majorList, err
}

func unmarshalJsonFile(filename string, dest interface{}) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	return json.Unmarshal(buf, &dest)
}
